package config

import (
	. "github.com/totoval/framework/config"
)

func init() {
	pushover := make(map[string]interface{})

	pushover["token"] = Env("PUSHOVER_TOKEN", "")
	pushover["user"] = Env("PUSHOVER_USER", "")
	pushover["device"] = Env("PUSHOVER_DEVICE", "")

	Add("pushover", pushover)
}
