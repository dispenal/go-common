package elastic

import (
	"context"

	common_utils "github.com/dispenal/go-common/utils"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/pkg/errors"
)

func MappingIndex(ctx context.Context, esClient *elasticsearch.Client, indexName string, mapping []byte) error {
	exists, err := isIndexExists(ctx, esClient, indexName)
	if err != nil {
		common_utils.LogError("error while checking index exists")
		return errors.Wrap(err, "error while checking index exists")
	}

	if exists {
		return nil
	}

	response, err := CreateIndex(ctx, esClient, indexName, mapping)
	if err != nil {
		common_utils.LogError("error while creating index")
		return errors.Wrap(err, "error while creating index")
	}

	defer response.Body.Close()

	if response.IsError() {
		return errors.New(response.String())
	}

	return nil
}

func isIndexExists(ctx context.Context, esClient *elasticsearch.Client, indexName string) (bool, error) {
	response, err := Exists(ctx, esClient, []string{indexName})
	if err != nil {
		common_utils.LogError("error while checking index exists")
		return false, errors.Wrap(err, "esclient.Exists")
	}
	defer response.Body.Close()

	if response.IsError() && response.StatusCode == 404 {
		common_utils.LogInfo("index not exists")
		return false, nil
	}

	return true, nil
}

func IntPointer(v int) *int {
	return &v
}

func Int32Pointer(v int32) *int32 {
	return &v
}

func Int64Pointer(v int64) *int64 {
	return &v
}

func StringPointer(v string) *string {
	return &v
}
