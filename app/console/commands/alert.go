package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/totoval/framework/config"
	"net/http"
	"strconv"
	"strings"

	"github.com/totoval/framework/biu"
	"github.com/totoval/framework/helpers/cache"
	"github.com/totoval/framework/helpers/hash"
	"github.com/totoval/framework/helpers/log"
	"github.com/totoval/framework/helpers/zone"
	"github.com/totoval/framework/logs"
	"github.com/totoval/framework/model/types/bigfloat"

	"github.com/totoval/framework/cmd"
)

func init() {
	cmd.Add(&Alert{})
}

const (
	HuobiMarketTradeBaseUrl = "https://api.huobi.pro/market/trade"
	PushOverMessageUrl = "https://api.pushover.net/1/messages.json"
)

type Alert struct {
	pair           string
	duration       zone.Duration
	difference     string
	upDifference   bigfloat.BigFloat
	downDifference bigfloat.BigFloat
}

func (hw *Alert) Command() string {
	return "crypto:alert {pair} {duration-minute} {difference}"
}

func (hw *Alert) Description() string {
	return "Alert crypto ticker"
}

type response struct {
	Status string `json:"status"`
	Ch     string `json:"ch"`
	Ts     int64  `json:"ts"`
	Tick   *Tick  `json:"tick"`
}
type Tick struct {
	Id   int64       `json:"id"`
	Ts   int64       `json:"ts"`
	Data []*TickData `json:"data"`
}
type TickData struct {
	Amount    bigfloat.BigFloat `json:"amount"`
	Ts        int64             `json:"ts"`
	Id        json.Number       `json:"id"`
	Price     bigfloat.BigFloat `json:"price"`
	Direction string            `json:"direction"`
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

	hw.pair = *pair
	durationInt64, err := strconv.ParseInt(*duration, 10, 64)
	if err != nil {
		return err
	}
	hw.duration = zone.Duration(durationInt64) * zone.Minute
	hw.difference = *difference

	if err := hw.initDifference(); err != nil {
		return err
	}

	var resp response
	statusCode, err := biu.Ready(biu.MethodGet, HuobiMarketTradeBaseUrl, &biu.Options{
		ProxyUrl: config.GetString("alert.proxy"),
		UrlParam: &biu.UrlParam{
			"symbol": hw.pair,
		},
	}).Biu().Object(&resp)

	if err != nil {
		return err
	}
	if statusCode != http.StatusOK {
		return log.Error(errors.New("status is not 200"), logs.Field{
			"status_code": statusCode,
		})
	}

	log.Info("data", logs.Field{
		"response": resp,
	})

	hw.addToCache(&resp)

	direction, err := hw.checkDirection(&nowData(&resp).Price)
	if err != nil {
		return err
	}

	if direction != Draw {
		hw.notify(direction, nowData(&resp))
	}

	return nil
}

func (hw *Alert) notify(direction Direction, nowData *TickData) {
	// every minute should notify once
	if cache.Has(hw.cacheKey(zone.Now())+"_notified"){
		return
	}

	diff, _ := strconv.ParseFloat(hw.difference, 64)

	dataStr, _ := json.Marshal(nowData)

	biu.Ready(biu.MethodPost, PushOverMessageUrl, &biu.Options{
		Body: &biu.Body{
			"token":   config.GetString("pushover.token"),
			"user":    config.GetString("pushover.user"),
			"device":  config.GetString("pushover.device"),

			"title":   fmt.Sprintf("%s %s %.2f%% !!! now price: %s", strings.ToUpper(hw.pair), direction, diff*100, nowData.Price.String()),
			"message": string(dataStr),
		},
	}).Biu()

	cache.Put(hw.cacheKey(zone.Now())+"_notified", true, zone.Now().Add(hw.duration+5*zone.Minute))
}

func nowData(resp *response) *TickData {
	if len(resp.Tick.Data) <= 0 {
		return nil
	}
	return resp.Tick.Data[0]
}

func (hw *Alert) addToCache(resp *response) {
	raw := zone.Duration(nowData(resp).Ts) * zone.Millisecond
	nanoSecond := int64(raw % zone.Second)
	second := int64(int64(raw)-nanoSecond) / int64(zone.Second)
	createdAt := zone.Unix(second, nanoSecond)
	cache.Put(hw.cacheKey(createdAt), nowData(resp).Price.String(), zone.Now().Add(hw.duration+5*zone.Minute))
}

func (hw *Alert) cacheKey(createdAt zone.Time) string {
	str := fmt.Sprintf("%s|%d|%d|%d|%d|%d", hw.pair, createdAt.Year(), createdAt.Month(), createdAt.Day(), createdAt.Hour(), createdAt.Minute())
	return hash.Md5(str)
}

type Direction = string

const (
	Up   Direction = "Up"
	Down           = "Down"
	Draw           = "Draw"
)

func (hw *Alert) initDifference() error {
	var difference bigfloat.BigFloat // 0.1
	if err := difference.CreateFromString(hw.difference, bigfloat.ToNearestEven); err != nil {
		return err
	}
	var one bigfloat.BigFloat // 1
	if err := one.CreateFromString("1", bigfloat.ToNearestEven); err != nil {
		return err
	}
	hw.upDifference.Add(one, difference)   // 1.1
	hw.downDifference.Sub(one, difference) // 0.9

	return nil
}

func (hw *Alert) checkDirection(nowPrice *bigfloat.BigFloat) (Direction, error) {
	beforePrice := cache.Get(hw.cacheKey(zone.Now().Add(-hw.duration)))
	if beforePrice == nil {
		return Draw, errors.New("not enough history data")
	}

	//debug.Dump(beforePrice)

	var beforePriceBigFloat bigfloat.BigFloat
	if err := beforePriceBigFloat.UnmarshalBinary([]byte(beforePrice.(string))); err != nil {
		return Draw, err
	}

	var cmpUpPrice bigfloat.BigFloat
	cmpUpPrice.Mul(beforePriceBigFloat, hw.upDifference) // before * 1.1
	var cmpDownPrice bigfloat.BigFloat
	cmpDownPrice.Mul(beforePriceBigFloat, hw.downDifference) // before * 0.9

	//debug.Dump(cmpUpPrice, cmpDownPrice, nowPrice)

	if nowPrice.Cmp(cmpDownPrice) <= 0 { // check before * 0.9 > now
		return Down, nil
	}
	if nowPrice.Cmp(cmpUpPrice) >= 0 { // check before * 1.1 < now
		return Up, nil
	}
	return Draw, nil
}
