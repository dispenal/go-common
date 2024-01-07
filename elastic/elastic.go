package elastic

import (
	"os"

	common_utils "github.com/dispenal/go-common/utils"
	"github.com/elastic/elastic-transport-go/v8/elastictransport"
	"github.com/elastic/go-elasticsearch/v8"
)

func NewElasticSearchClient(cfg common_utils.BaseConfig) (*elasticsearch.Client, error) {
	config := elasticsearch.Config{
		Addresses: cfg.ElasticsearchHost,
		Username:  cfg.ElasticsearchUser,
		Password:  cfg.ElasticsearchPassword,
	}

	if cfg.ElasticsearchLogging {
		config.Logger = &elastictransport.ColorLogger{Output: os.Stdout, EnableRequestBody: true, EnableResponseBody: true}
	}

	client, err := elasticsearch.NewClient(config)
	if err != nil {
		return nil, err
	}

	return client, nil
}
