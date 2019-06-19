package console

import (
	"fmt"
	"github.com/totoval/framework/cmd"
	"github.com/totoval/framework/config"
)

func Schedule(schedule *cmd.Schedule) {
	schedule.Command(
		fmt.Sprintf("crypto:alert %s %d %s",
			config.GetString("alert.schedule_pair"),
			config.GetUint("alert.schedule_duration_minutes"),
			config.GetString("alert.schedule_difference"),
		),
	).EverySeconds(config.GetUint("alert.schedule_interval_seconds"))
}
