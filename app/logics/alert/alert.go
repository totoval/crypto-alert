package alert

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/totoval/framework/helpers/cache"
	"github.com/totoval/framework/helpers/hash"
	"github.com/totoval/framework/helpers/log"
	"github.com/totoval/framework/helpers/zone"
	"github.com/totoval/framework/logs"
	"github.com/totoval/framework/model/types/bigfloat"
	"strconv"
	"strings"
)

type alert struct {
	pair       string
	duration   zone.Duration
	difference string

	upDifference   bigfloat.BigFloat
	downDifference bigfloat.BigFloat
}

func New(pair string, duration zone.Duration, difference string) (*alert, error) {
	a := &alert{pair: pair, duration: duration, difference: difference}
	err := a.initDifference()
	return a, err
}

type Response struct {
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

func (alert *alert) Fetch(fetcher Fetcher, notifier Notifier) error {
	resp, err := fetcher.Fetch(alert.pair)
	if err != nil {
		return err
	}

	log.Info("data", logs.Field{
		"response": resp,
	})

	alert.addToCache(resp)

	direction, err := alert.checkDirection(&nowData(resp).Price)
	if err != nil {
		return err
	}

	if direction != Draw {
		return alert.notify(notifier, direction, nowData(resp))
	}

	return nil
}

func (alert *alert) notify(notifier Notifier, direction Direction, nowData *TickData) error {

	cache.Put(alert.pair+"_notified", nowData.Price.String(), zone.Now().Add(alert.duration))

	diff, _ := strconv.ParseFloat(alert.difference, 64)

	dataStr, _ := json.Marshal(nowData)

	return notifier.Notify(strings.ToUpper(alert.pair), direction, fmt.Sprintf("%.2f%%", diff*100), nowData.Price.String(), string(dataStr))
}

func nowData(resp *Response) *TickData {
	if len(resp.Tick.Data) <= 0 {
		return nil
	}
	return resp.Tick.Data[0]
}

func (alert *alert) addToCache(resp *Response) {
	raw := zone.Duration(nowData(resp).Ts) * zone.Millisecond
	nanoSecond := int64(raw % zone.Second)
	second := int64(int64(raw)-nanoSecond) / int64(zone.Second)
	createdAt := zone.Unix(second, nanoSecond)
	cache.Put(alert.cacheKey(createdAt), nowData(resp).Price.String(), zone.Now().Add(alert.duration+5*zone.Minute))
}

func (alert *alert) cacheKey(createdAt zone.Time) string {
	str := fmt.Sprintf("%s|%d|%d|%d|%d|%d", alert.pair, createdAt.Year(), createdAt.Month(), createdAt.Day(), createdAt.Hour(), createdAt.Minute())
	return hash.Md5(str)
}

type Direction = string

const (
	Up   Direction = "Up"
	Down           = "Down"
	Draw           = "Draw"
)

func (alert *alert) initDifference() error {
	var difference bigfloat.BigFloat // 0.1
	if err := difference.CreateFromString(alert.difference, bigfloat.ToNearestEven); err != nil {
		return err
	}
	var one bigfloat.BigFloat // 1
	if err := one.CreateFromString("1", bigfloat.ToNearestEven); err != nil {
		return err
	}
	alert.upDifference.Add(one, difference)   // 1.1
	alert.downDifference.Sub(one, difference) // 0.9

	return nil
}

func (alert *alert) checkDirection(nowPrice *bigfloat.BigFloat) (Direction, error) {
	beforePrice := cache.Get(alert.pair + "_notified") // beforeNotifiedPrice  first use before notified price for checking, to avoid notify every interval
	if beforePrice == nil {
		beforePrice = cache.Get(alert.cacheKey(zone.Now().Add(-alert.duration)))
		if beforePrice == nil {
			return Draw, errors.New("not enough history data")
		}
	}

	//debug.Dump(beforePrice)

	var beforePriceBigFloat bigfloat.BigFloat
	if err := beforePriceBigFloat.UnmarshalBinary([]byte(beforePrice.(string))); err != nil {
		return Draw, err
	}

	var cmpUpPrice bigfloat.BigFloat
	cmpUpPrice.Mul(beforePriceBigFloat, alert.upDifference) // before * 1.1
	var cmpDownPrice bigfloat.BigFloat
	cmpDownPrice.Mul(beforePriceBigFloat, alert.downDifference) // before * 0.9

	//debug.Dump(cmpUpPrice, cmpDownPrice, nowPrice)

	if nowPrice.Cmp(cmpDownPrice) <= 0 { // check before * 0.9 > now
		return Down, nil
	}
	if nowPrice.Cmp(cmpUpPrice) >= 0 { // check before * 1.1 < now
		return Up, nil
	}
	return Draw, nil
}
