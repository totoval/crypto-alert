# Crypto Alert
Help you build your own crypto currency price alert system using [@Totoval](https://github.com/totoval/totoval) and [Pushover](https://pushover.net)

<img src="https://raw.githubusercontent.com/totoval/crypto-alert/master/readme_assets/WechatIMG539.jpeg" alt="pushover" width="300" /> <img src="https://raw.githubusercontent.com/totoval/crypto-alert/master/readme_assets/WechatIMG538.jpeg" alt="huobi.pro" width="300" />

# What is it?
Alert me by Pushover **if `btcusdt`'s price is higher/lower than `5min ago`'s price of `1%`**

# How to Use
## 0x01 Configure

```bash
tee ./.env.json <<- "EOF"
{
  "APP_NAME": "CryptoAlert",
  "APP_ENV": "develop",
  "APP_DEBUG": true,
  "APP_PORT": 8080,
  "APP_LOCALE": "en",
  "APP_KEY": "xxx",

  "CACHE_DRIVER": "memory",

  "ALERT_PROXY": "",
  "ALERT_SCHEDULE_DURATION": 30,
  "ALERT_SCHEDULE_PAIR": "btcusdt",
  "ALERT_SCHEDULE_DIFFERENCE": "0.01",
  "ALERT_SCHEDULE_INTERVAL": 1,

  "PUSHOVER_TOKEN": "{YOUR-PUSHOVER-TOKEN}",
  "PUSHOVER_USER": "{YOUR-PUSHOVER-USER}",
  "PUSHOVER_DEVICE": ""
}
EOF
```

1. Set a `pair` like `btcusdt` referred at [here](https://huobiapi.github.io/docs/spot/v1/cn/#0e505d18dc)
2. Set an interval about which the price of time will be used to be compared with current ticker, for example: 
    set `5` means will let the program compare the price of `now` and `5min ago`
3. Set a difference rate, for example: if price is higher/lower than `5min ago`'s price of `1%`, then set the difference to `0.01` 

---

* You could configure it in **command line** with `go run artisan.go crypto:alert btcusdt 5 0.01`  
* For using **`.env.json`**, just copy `.env.example.json` to `.env.json`, set your env, then rock!

## 0x02 Run

For `loop`:
```bash
go run artisan.go schedule:run
```
Or for `once`:
```bash
go run artisan.go crypto:alert btcusdt 5 0.01
```

**Or you could build a binary~, as you wish!**   
**For linux:**  
`./artisan_linux schedule:run`  
**For Mac:**  
`./artisan_mac schedule:run`

## 0x03 Or just use our example
**Just follow the [link](https://pushover.net/subscribe/CryptoAlert-7q6qr9j2qr4mm39) to subscribe a `huobi.pro`-`30-min`-`btcusdt`-`1%UP/DOWN` alert!**

# Implement with Yours
* By implement the `Fetcher` Or `Notifier`, you could build your own **Crypto-Alert**
* `Fetcher`: `app/logics/alert/fetchers`
* `Notifier`: `app/logics/alert/notifiers`

**Thanks for [@Totoval](https://github.com/totoval/totoval)**
