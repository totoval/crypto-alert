package fetchers

import (
	"errors"
	"github.com/totoval/framework/biu"
	"github.com/totoval/framework/config"
	"github.com/totoval/framework/helpers/log"
	"github.com/totoval/framework/logs"
	"net/http"
	"totoval/app/logics/alert"
)

const (
	HuobiMarketTradeBaseUrl = "https://api.huobi.pro/market/trade"
)

type HuoBi struct {
}

func (hb *HuoBi) Fetch(pair string) (resp *alert.Response, err error) {
	statusCode, err := biu.Ready(biu.MethodGet, HuobiMarketTradeBaseUrl, &biu.Options{
		ProxyUrl: config.GetString("alert.proxy"),
		UrlParam: &biu.UrlParam{
			"symbol": pair,
		},
	}).Biu().Object(&resp)

	if err != nil {
		return nil, err
	}
	if statusCode != http.StatusOK {
		return nil, log.Error(errors.New("status is not 200"), logs.Field{
			"status_code": statusCode,
		})
	}

	return resp, nil
}
