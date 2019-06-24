package commands

import (
	"errors"
	"strconv"
	"totoval/app/logics/alert"
	"totoval/app/logics/alert/fetchers"
	"totoval/app/logics/alert/notifiers"

	"github.com/totoval/framework/cmd"
	"github.com/totoval/framework/helpers/log"
	"github.com/totoval/framework/helpers/zone"
)

func init() {
	cmd.Add(&Alert{})
}


type Alert struct {
}

func (hw *Alert) Command() string {
	return "crypto:alert {pair} {duration-minute} {difference}"
}

func (hw *Alert) Description() string {
	return "Alert crypto ticker"
}


func (hw *Alert) Handler(arg *cmd.Arg) error {
	pair, err := arg.Get("pair")
	if err != nil {
		return err
	}
	duration, err := arg.Get("duration-minute")
	if err != nil {
		return err
	}
	difference, err := arg.Get("difference")
	if err != nil {
		return err
	}

	if pair == nil {
		return log.Error(errors.New("pair not set"))
	}
	if duration == nil {
		return log.Error(errors.New("duration not set"))
	}
	if difference == nil {
		return log.Error(errors.New("difference not set"))
	}

	durationInt64, err := strconv.ParseInt(*duration, 10, 64)
	if err != nil {
		return err
	}
	a, err := alert.New(*pair, zone.Duration(durationInt64) * zone.Minute, *difference)
	if err != nil {
		return err
	}

	// you can diy your own Fetchers or Notifiers
	return a.Fetch(new(fetchers.HuoBi), new(notifiers.Pushover))
}

