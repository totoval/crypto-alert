package notifiers

import (
	"fmt"
	"github.com/totoval/framework/biu"
	"github.com/totoval/framework/config"
	"totoval/app/logics/alert"
)
const (
	PushOverMessageUrl = "https://api.pushover.net/1/messages.json"
)

type Pushover struct {

}
func (po *Pushover)Notify(pair string, direction alert.Direction, differencePercentage string, price string, rawDataStr string) error {
	directionStr := "-"
	if direction == alert.Up {
		directionStr = "📈"
	}else{
		directionStr = "️📉"
	}

	_, err := biu.Ready(biu.MethodPost, PushOverMessageUrl, &biu.Options{
		Body: &biu.Body{
			"token":   config.GetString("pushover.token"),
			"user":    config.GetString("pushover.user"),
			"device":  config.GetString("pushover.device"),

			"title":   fmt.Sprintf("[%s] %s %s !!! {%s}", pair,  directionStr, differencePercentage, price),
			"message": string(rawDataStr),
		},
	}).Biu().Status()
	return err
}
