package common_utils

import (
	"fmt"
	"reflect"
)

func BuildCacheKey(key string, identifier string, funcName string, args ...any) string {
	cacheKey := fmt.Sprintf("%s-%s-%s", key, identifier, funcName)
	cacheArgs := ""
	for _, arg := range args {

		if cacheArgs != "" {
			cacheArgs += "|"
		}

		v := reflect.ValueOf(arg)
		for i := 0; i < v.NumField(); i++ {
			isEmpty := v.Field(i).Interface() == ""
			if !isEmpty && cacheArgs == "" {
				cacheArgs += fmt.Sprintf("%s:%v", v.Type().Field(i).Name, v.Field(i).Interface())
			} else if !isEmpty && string(cacheArgs[len(cacheArgs)-1]) == "|" {
				cacheArgs += fmt.Sprintf("%s:%v", v.Type().Field(i).Name, v.Field(i).Interface())
			} else if !isEmpty {
				cacheArgs += fmt.Sprintf(",%s:%v", v.Type().Field(i).Name, v.Field(i).Interface())
			}
		}
	}

	return fmt.Sprintf("%s|%s", cacheKey, cacheArgs)
}

func BuildPrefixKey(keys ...string) string {
	prefixKey := ""

	for _, key := range keys {
		if prefixKey == "" {
			prefixKey += key
		} else {
			prefixKey += fmt.Sprintf("-%s", key)
		}
	}

	return prefixKey
}
