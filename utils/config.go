package common_utils

import (
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type BaseConfig struct {
	ServiceName            string        `mapstructure:"SERVICE_NAME"`
	ServiceEnv             string        `mapstructure:"SERVICE_ENV"`
	ServiceHost            string        `mapstructure:"SERVICE_HOST"`
	ServiceHttpPort        string        `mapstructure:"SERVICE_HTTP_PORT"`
	ServiceGrpcPort        string        `mapstructure:"SERVICE_GRPC_PORT"`
	GcpProjectId           string        `mapstructure:"GCP_PROJECT_ID"`
	PubsubDlq              string        `mapstructure:"PUBSUB_DLQ_TOPIC"`
	RedisHost              string        `mapstructure:"REDIS_HOST"`
	RedisPort              string        `mapstructure:"REDIS_PORT"`
	RedisUser              string        `mapstructure:"REDIS_USER"`
	RedisPassword          string        `mapstructure:"REDIS_PASSWORD"`
	RedisCacheExpire       int           `mapstructure:"REDIS_DEFAULT_CACHE_EXPIRE"`
	MongoHost              string        `mapstructure:"MONGO_HOST"`
	MongoPort              string        `mapstructure:"MONGO_PORT"`
	MongoUser              string        `mapstructure:"MONGO_USER"`
	MongoPassword          string        `mapstructure:"MONGO_PASSWORD"`
	MongoDb                string        `mapstructure:"MONGO_DATABASE"`
	MongoReplicaSet        string        `mapstructure:"MONGO_REPLICA_SET,default=rs0"`
	PostgresHost           string        `mapstructure:"POSTGRES_HOST"`
	PostgresPort           string        `mapstructure:"POSTGRES_PORT"`
	PostgresUser           string        `mapstructure:"POSTGRES_USER"`
	PostgresPassword       string        `mapstructure:"POSTGRES_PASSWORD"`
	PostgresDb             string        `mapstructure:"POSTGRES_DATABASE"`
	ElasticsearchHost      []string      `mapstructure:"ELASTICSEARCH_HOST"`
	ElasticsearchUser      string        `mapstructure:"ELASTICSEARCH_USER"`
	ElasticsearchPassword  string        `mapstructure:"ELASTICSEARCH_PASSWORD"`
	ElasticsearchLogging   bool          `mapstructure:"ELASTICSEARCH_LOGGING,default=false"`
	JaegerEnable           bool          `mapstructure:"JAEGER_ENABLE"`
	JaegerHost             string        `mapstructure:"JAEGER_HOST"`
	JaegerPort             string        `mapstructure:"JAEGER_PORT"`
	JaegerLogSpans         bool          `mapstructure:"JAEGER_LOG_SPANS"`
	S3Endpoint             string        `mapstructure:"S3_ENDPOINT"`
	S3AccessKey            string        `mapstructure:"S3_ACCESS_KEY"`
	S3SecretKey            string        `mapstructure:"S3_SECRET_KEY"`
	S3Region               string        `mapstructure:"S3_REGION"`
	S3PublicBucket         string        `mapstructure:"S3_PUBLIC_BUCKET"`
	S3PrivateBucket        string        `mapstructure:"S3_PRIVATE_BUCKET"`
	S3PublicUrl            string        `mapstructure:"S3_PUBLIC_URL"`
	S3PreSignedExpire      time.Duration `mapstructure:"S3_PRESIGNED_EXPIRE"`
	JwtSecretKey           string        `mapstructure:"JWT_SECRET_KEY"`
	JwtAccessTokenExpire   time.Duration `mapstructure:"JWT_ACCESS_TOKEN_EXPIRE"`
	JwtRefreshTokenExpire  time.Duration `mapstructure:"JWT_REFRESH_TOKEN_EXPIRE"`
	KafkaBrokers           []string      `mapstructure:"KAFKA_BROKERS"`
	KafkaTopics            []string      `mapstructure:"KAFKA_TOPICS"`
	KafkaGroupID           string        `mapstructure:"KAFKA_GROUP_ID"`
	KafkaDlqTopic          string        `mapstructure:"KAFKA_DLQ_TOPIC"`
	KafkaDlqRetry          int           `mapstructure:"KAFKA_DLQ_RETRY,default=3"`
	KafkaAutoCommit        bool          `mapstructure:"KAFKA_AUTO_COMMIT,default=false"`
	KafkaAutoTopicCreation bool          `mapstructure:"KAFKA_AUTO_TOPIC_CREATION,default=true"`
	KafkaReplicationFactor int           `mapstructure:"KAFKA_REPLICATION_FACTOR,default=1"`
}

func LoadBaseConfig(path string, configName string) (*BaseConfig, error) {
	if configName != "" {
		viper.SetConfigName(configName)
	} else {
		viper.SetConfigName(".env")
	}

	if path == "" {
		path = "."
	}

	conf := &BaseConfig{}

	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "viper.ReadInConfig")
	}

	if err := viper.Unmarshal(conf); err != nil {
		return nil, errors.Wrap(err, "viper.Unmarshal")
	}

	return conf, nil
}

func CheckAndSetConfig(path string, configName string) *BaseConfig {
	config, err := LoadBaseConfig(path, configName)
	if err != nil {
		panic(err)
	}

	if config.ServiceEnv == TEST {
		os.Setenv("SERVICE_ENV", fmt.Sprintf("%s-test", config.ServiceEnv))
		config, err = LoadBaseConfig(path, "test")
		if err != nil {
			panic(err)
		}
	}

	if config.ServiceEnv == DEVELOPMENT {
		os.Setenv("SERVICE_ENV", fmt.Sprintf("%s-local", config.ServiceEnv))
	}

	return config
}
