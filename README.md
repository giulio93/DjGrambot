# DjGrambot - Go Telegram Bot
A simple Telegram bot written in Go, that download and (will) mix your favourite music

## How to run this bot

- Set your env variable: export TELEGRAM_APITOKEN='<your_telegram_api_token>'
- go run main.go

## ChangeLog

### Features
- Reply to specific commands
- Upload and send gif
- Show keyboard and send the pressed key
- Send a link
- Send a position
- Download a video/audio using @vid bot
- Go routine for video download
- Reply to the @vid with the downloaded audio
### Bug fix

- Handle errors during video download
- Too many request error handled via timer
- Avoid streaming channels and playlist
- Avoid download too long file



## TODO
- Mixing track
- Get the local playlist
- Get the list of command
- Show download percentage


## Dependencies
- [Go Telegram Bot API](https://go-telegram-bot-api.dev/getting-started/index.html)
- [Youtube downloader](github.com/kkdai/youtube)
- [Mixing tracks](https://github.com/go-mix/mix)