package redis_client

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	common_utils "github.com/dispenal/go-common/utils"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RedisClient interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	Scan(ctx context.Context, cursor uint64, match string, count int64) *redis.ScanCmd
	Close() error
}

type CacheSvc interface {
	Set(ctx context.Context, key string, data any, duration ...time.Duration) error
	Get(ctx context.Context, key string, output any) error
	Del(ctx context.Context, key string) error
	DelByPrefix(ctx context.Context, prefixName string)
	GetOrSet(ctx context.Context, key string, function func() any, duration ...time.Duration) (any, error)
	CloseClient() error
}

type CacheSvcImpl struct {
	config  *common_utils.BaseConfig
	cacheDb RedisClient
}

func NewCacheSvc(config *common_utils.BaseConfig, cacheDb RedisClient) CacheSvc {
	return &CacheSvcImpl{
		config:  config,
		cacheDb: cacheDb,
	}
}

func NewRedisClient(config *common_utils.BaseConfig) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.RedisHost, config.RedisPort),
		Password: config.RedisPassword, // no password set
		DB:       0,                    // use default DB
		Username: config.RedisUser,
	})
	return rdb
}

func NewRedisClientForTesting(config *common_utils.BaseConfig) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.RedisHost, config.RedisPort),
		Password: config.RedisPassword, // no password set
		DB:       1,                    // use testing DB
		Username: config.RedisUser,
	})
	return rdb
}

func (s *CacheSvcImpl) Set(ctx context.Context, key string, data any, duration ...time.Duration) error {
	dataErr, isDataErr := data.(error)
	if isDataErr {
		common_utils.LogInfo("not save data to cache (data error)")
		return dataErr
	}

	appErr, isAppErr := data.(common_utils.AppError)
	if isAppErr {
		common_utils.LogInfo("not save data to cache (app error)")
		return &appErr
	}

	validationErrs, isValidationErrs := data.(common_utils.ValidationErrors)
	if isValidationErrs {
		common_utils.LogInfo("not save data to cache (validation errors)")
		return &validationErrs
	}

	if data != nil {
		if reflect.TypeOf(data).Kind() == reflect.Slice {
			if reflect.ValueOf(data).Len() == 0 {
				common_utils.LogInfo("no data to save, array is empty")
				return nil
			}
		}

		cacheData, err := common_utils.Marshal(data)
		if err != nil {
			return err
		}

		common_utils.LogInfo(fmt.Sprintf("set data to cache with key --> %s", key))

		expiration := time.Duration(s.config.RedisCacheExpire) * time.Second
		if len(duration) > 0 {
			expiration = duration[0]
		}

		return s.cacheDb.Set(ctx, key, cacheData, expiration).Err()
	}

	common_utils.LogInfo(fmt.Sprintf("not save data to cache, key --> %s", key))

	return nil
}

func (s *CacheSvcImpl) Get(ctx context.Context, key string, output any) error {
	val, err := s.cacheDb.Get(ctx, key).Result()
	if err != nil {
		common_utils.LogInfo(fmt.Sprintf("failed when getting key -> %s | error: %v", key, err))
		return err
	}

	err = json.Unmarshal([]byte(val), &output)

	if err != nil {
		common_utils.LogInfo("failed when unmarshal data")
		return err
	}

	common_utils.LogInfo(fmt.Sprintf("get data from cache with key --> %s", key))

	return nil
}

func (s *CacheSvcImpl) DelByPrefix(ctx context.Context, prefixName string) {
	var foundedRecordCount int = 0
	iter := s.cacheDb.Scan(ctx, 0, fmt.Sprintf("%s*", prefixName), 0).Iterator()
	common_utils.LogInfo(fmt.Sprintf("your search pattern: %s", prefixName))

	for iter.Next(ctx) {
		common_utils.LogInfo(fmt.Sprintf("deleted= %s", iter.Val()))
		s.cacheDb.Del(ctx, iter.Val())
		foundedRecordCount++
	}

	if err := iter.Err(); err != nil {
		common_utils.LogError("failed when deleting cache", zap.Error(err))
	}

	common_utils.LogInfo(fmt.Sprintf("deleted Count %d", foundedRecordCount))
}

func (s *CacheSvcImpl) GetOrSet(ctx context.Context, key string, function func() any, duration ...time.Duration) (any, error) {
	var data any
	err := s.Get(ctx, key, &data)

	if err != nil && err == redis.Nil {
		data = function()
		err := s.Set(ctx, key, data, duration...)

		return data, err
	}

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *CacheSvcImpl) Del(ctx context.Context, key string) error {
	err := s.cacheDb.Del(ctx, key).Err()
	if err != nil {
		common_utils.LogError(fmt.Sprintf("failed when deleting cache key: %s", key), zap.Error(err))
		return err
	}

	return nil
}

func (s *CacheSvcImpl) CloseClient() error {
	return s.cacheDb.Close()
}
