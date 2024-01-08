package elastic

import (
	"bytes"
	"context"

	common_utils "github.com/dispenal/go-common/utils"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/pkg/errors"
)

func Index(ctx context.Context, transport esapi.Transport, index, documentID string, v any) (*esapi.Response, error) {
	reqBytes, err := common_utils.Marshal(v)
	if err != nil {
		return nil, errors.Wrap(err, "json.Marshal")
	}

	request := esapi.IndexRequest{
		Index:      index,
		DocumentID: documentID,
		Body:       bytes.NewBuffer(reqBytes),
	}

	return request.Do(ctx, transport)
}

func BulkIndex(ctx context.Context, transport esapi.Transport, index, documentID string, v []any) (*esapi.Response, error) {
	reqBytes, err := common_utils.Marshal(v)
	if err != nil {
		return nil, errors.Wrap(err, "json.Marshal")
	}

	request := esapi.BulkRequest{
		Index: index,
		Body:  bytes.NewBuffer(reqBytes),
	}

	return request.Do(ctx, transport)
}

func Update(ctx context.Context, transport esapi.Transport, index, documentID string, document any) (*esapi.Response, error) {
	doc := Doc{Doc: document}
	reqBytes, err := common_utils.Marshal(&doc)
	if err != nil {
		return nil, err
	}

	request := esapi.UpdateRequest{
		Index:      index,
		DocumentID: documentID,
		Body:       bytes.NewReader(reqBytes),
		Refresh:    "true",
	}

	return request.Do(ctx, transport)
}

func Delete(ctx context.Context, transport esapi.Transport, index, documentID string) (*esapi.Response, error) {
	request := esapi.DeleteRequest{
		Index:      index,
		DocumentID: documentID,
		Refresh:    "true",
	}

	return request.Do(ctx, transport)
}
