# Crypto Alert
Help you build your own crypto currency price alert system using [@Totoval](https://github.com/totoval/totoval) and [Pushover](https://pushover.net)

<img src="https://raw.githubusercontent.com/totoval/crypto-alert/master/readme_assets/WechatIMG539.jpeg" alt="pushover" width="300" /> <img src="https://raw.githubusercontent.com/totoval/crypto-alert/master/readme_assets/WechatIMG538.jpeg" alt="huobi.pro" width="300" />

# What is it
Alert me by Pushover **if `btcusdt`'s price is higher/lower than `5min ago`'s price of `1%`**

# How to Use
## 0x01 Configure
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
**Thanks for [@Totoval](https://github.com/totoval/totoval)**

## 0x03 Or just use our example
**Just follow the [link](https://pushover.net/subscribe/CryptoAlert-7q6qr9j2qr4mm39) to subscribe a `huobi.pro`-`30-min`-`btcusdt`-`1%UP/DOWN` alert!**

---
---
---

<p align="center"><img src="https://raw.githubusercontent.com/totoval/art/master/repo_use/logo-with-words-landscape.png?s=200&v=4"></p>

![GitHub last commit](https://img.shields.io/github/last-commit/totoval/totoval.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/totoval/totoval)](https://goreportcard.com/report/github.com/totoval/totoval)
![Travis (.org)](https://img.shields.io/travis/totoval/totoval.svg)
![GitHub top language](https://img.shields.io/github/languages/top/totoval/totoval.svg)
![Inspired by Laravel](https://img.shields.io/badge/Inspired%20by-Laravel-red.svg)
![GitHub](https://img.shields.io/github/license/totoval/totoval.svg)

## About Totoval
Totoval is an API web framework that helps Golang engineers build a performance-boiled project quickly, easily, and securely. It is more like a scaffolding, respecting Golang's programming philosophy, supported by a number of highly acclaimed, high-performance core components, as well as many easy-to-use components to quickly adapt to more business scenarios. We also believe that development must be an enjoyable and creative experience. Totoval frees developers from the painful coding process. Do less, think more.



## Roadmap
- [x] Env Configuration
- [x] Groupable Router
- [x] Request Middleware
- [x] Request Validator
- [x] Database Migration
- [x] Model Validator
- [x] Model Helper - such as `Pagination`
- [x] BigInt,BigFloat Support
- [x] Orm: Mysql
- [x] User Token JWT Support
- [x] Random Code Generate and Verification
- [x] Random String Helper
- [x] Locale Middleware
- [x] Gin Validator Upgrade to v9
- [x] Password Encryption
- [x] Validation Error Multi-Language Support
- [x] Request Logger Middleware
- [x] Infinity User Affiliation System
- [x] User Email Validation via Notification
- [x] Views Support
- [x] Language Package
- [x] Cache: Memory
- [x] Cache: Redis
- [x] Queue, Worker `nsq`
- [x] Event, Listener
- [x] Custom Artisan Command Line
- [x] Task Scheduling
- [x] Logo
- [x] Http Request Package `biu`
- [ ] Error Handler
- [ ] Model Getter/Setter
- [ ] File Storage
- [ ] User Authorization
- [ ] Database Seeder
- [ ] More Unit Test
- [ ] Websocket Support
- [ ] Website && Document
- [ ] CI

## Thanks
* gin
* gorm
* validator.v9
* viper
* big
* jwt
* i18n
* urfave/cli
* fatih/color
* golang/protobuf
* nsqio/go-nsq
* robfig/cron
* ztrue/tracerr
* go-redis/redis
