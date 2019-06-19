package config

import (
	. "github.com/totoval/framework/config"
)

func init() {
	cache := make(map[string]interface{})

	cache["default"] = Env("CACHE_DRIVER", "memory")
	cache["stores"] = map[string]interface{}{
		"memory": map[string]interface{}{
			"driver":                    "memory",
			"default_expiration_minute": 5,
			"cleanup_interval_minute":   5,
			"prefix":                    Env("APP_NAME", "totoval").(string) + "_cache_",
		},
		"redis": map[string]interface{}{
			"driver":     "redis",
			"connection": "cache",
		},
	}

	Add("cache", cache)
}
