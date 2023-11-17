package s3

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"regexp"
	"sort"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	common_utils "github.com/dispenal/go-common/utils"
)

func Init(baseConfig *common_utils.BaseConfig) (*s3.Client, error) {
	creeds := credentials.NewStaticCredentialsProvider(baseConfig.S3AccessKey, baseConfig.S3SecretKey, "")

	customEndpointResolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:       "aws",
			URL:               baseConfig.S3Endpoint,
			SigningRegion:     baseConfig.S3Region,
			HostnameImmutable: true,
		}, nil
	})

	logMode := aws.ClientLogMode(0)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	no_ssl_verify := &http.Client{Transport: tr}

	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithCredentialsProvider(creeds),
		config.WithRegion(baseConfig.S3Region),
		config.WithEndpointResolver(customEndpointResolver),
		config.WithClientLogMode(logMode),
		config.WithHTTPClient(no_ssl_verify),
		config.WithRetryer(func() aws.Retryer {
			return aws.NopRetryer{}
		}),
	)

	return s3.NewFromConfig(cfg), err
}

func PreSignClient(client *s3.Client) *s3.PresignClient {
	return s3.NewPresignClient(client)
}

type S3ClientImpl struct {
	client  S3Client
	preSign S3PreSign
	config  *common_utils.BaseConfig
}

func NewS3Client(
	client S3Client,
	preSign S3PreSign,
	config *common_utils.BaseConfig,
) S3File {
	return &S3ClientImpl{
		client:  client,
		preSign: preSign,
		config:  config,
	}
}

func (s *S3ClientImpl) UploadPrivateFile(ctx context.Context, file multipart.File, path string) (string, error) {

	_, err := s.client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String(s.config.S3PrivateBucket),
		Key:    aws.String(path),
		Body:   file,
	})

	if err != nil {
		return "", err
	}

	preSignUrl, err := s.GetPreSignUrl(ctx, path)

	if err != nil {
		return "", err
	}

	return preSignUrl, nil
}

func (s *S3ClientImpl) UploadPublicFile(ctx context.Context, file multipart.File, path string) (string, error) {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.config.S3PublicBucket),
		Key:         aws.String(path),
		Body:        file,
		ContentType: aws.String(getExt(path)),
	})

	if err != nil {
		return "", err
	}

	return s.BuildPublicUrl(path), nil
}

func (s *S3ClientImpl) UploadPartPublicFile(ctx context.Context, file multipart.File, path string) (string, error) {

	resp, err := s.client.CreateMultipartUpload(ctx, &s3.CreateMultipartUploadInput{
		Bucket:      aws.String(s.config.S3PublicBucket),
		Key:         aws.String(path),
		ContentType: aws.String(getExt(path)),
	})
	if err != nil {
		return "", err
	}

	uploadId := resp.UploadId

	var parts []types.CompletedPart
	var totalSize int64 = 0

	var buff bytes.Buffer

	_, err = file.Seek(0, 0)
	if err != nil {
		return "", err
	}

	_, err = io.Copy(&buff, file)
	if err != nil {
		return "", err
	}

	totalSize = int64(buff.Len())

	buffer := buff.Bytes()

	chunkSize := int64(5 * 1024 * 1024) // 5 MB
	numChunks := (totalSize + chunkSize - 1) / chunkSize
	chunks := make([][]byte, numChunks)
	for i := int64(0); i < numChunks; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if end > totalSize {
			end = totalSize
		}
		chunks[i] = buffer[start:end]
	}

	var wg sync.WaitGroup
	errCh := make(chan error, len(chunks))
	for i, chunk := range chunks {
		wg.Add(1)
		go func(i int, chunk []byte) {
			defer wg.Done()
			resp2, err := s.client.UploadPart(ctx, &s3.UploadPartInput{
				Bucket:        aws.String(s.config.S3PublicBucket),
				Key:           resp.Key,
				UploadId:      resp.UploadId,
				PartNumber:    *aws.Int32(int32(i) + 1),
				Body:          bytes.NewReader(chunk),
				ContentLength: *aws.Int64(int64(len(chunk))),
			})

			parts = append(parts, types.CompletedPart{
				ETag:       resp2.ETag,
				PartNumber: *aws.Int32(int32(i) + 1),
			})
			if err != nil {
				errCh <- err
			}

		}(i, chunk)
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		return "", err
	}

	sort.Slice(parts, func(i, j int) bool {
		return parts[i].PartNumber < parts[j].PartNumber
	})

	if len(parts) == len(chunks) {
		_, err = s.client.CompleteMultipartUpload(ctx, &s3.CompleteMultipartUploadInput{
			Bucket:   aws.String(s.config.S3PublicBucket),
			Key:      resp.Key,
			UploadId: uploadId,
			MultipartUpload: &types.CompletedMultipartUpload{
				Parts: parts,
			},
		})

		if err != nil {
			return "", err
		}

		return s.BuildPublicUrl(path), nil
	}

	return "", errors.New("upload failed")

}

func (s *S3ClientImpl) DeleteFile(ctx context.Context, bucketName string, path string) error {
	_, err := s.client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(path),
	})

	return err
}

func (s *S3ClientImpl) GetPreSignUrl(ctx context.Context, path string) (string, error) {
	params := &s3.GetObjectInput{
		Bucket:              aws.String(s.config.S3PrivateBucket),
		Key:                 aws.String(path),
		ResponseContentType: aws.String(getExt(path)),
	}

	resp, err := s.preSign.PresignGetObject(ctx, params, func(po *s3.PresignOptions) {
		po.Expires = time.Duration(s.config.S3PreSignedExpire) * time.Second
	})

	if err != nil {
		return "", err
	}

	return resp.URL, nil
}

func (s *S3ClientImpl) BuildPublicUrl(path string) string {

	var url string
	if s.config.S3PublicUrl != "" {
		url = s.config.S3PublicUrl + string(filepath.Separator) + path
	} else {
		url = fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s",
			s.config.S3PublicBucket,
			s.config.S3Region,
			path,
		)
	}

	re := regexp.MustCompile(`([^:])//+`)

	newURLString := re.ReplaceAllString(url, "$1/")

	return newURLString
}

func getExt(pathOrFilename string) string {
	return mime.TypeByExtension(filepath.Ext(pathOrFilename))
}
