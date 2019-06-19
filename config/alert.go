package config

import (
	. "github.com/totoval/framework/config"
)

func init() {
	alert := make(map[string]interface{})

	alert["proxy"] = Env("ALERT_PROXY", "")
	alert["schedule_pair"] = Env("ALERT_SCHEDULE_PAIR", "btcusdt")
	alert["schedule_duration_minutes"] = Env("ALERT_SCHEDULE_DURATION", 5)
	alert["schedule_difference"] = Env("ALERT_SCHEDULE_DIFFERENCE", "0.01")

	alert["schedule_interval_seconds"] = Env("ALERT_SCHEDULE_INTERVAL", 1)

	Add("alert", alert)
}
