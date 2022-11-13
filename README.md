# GO Telegram Bot

A simple Telegram Bot in GO

## How to run

- Set your env variable: export TELEGRAM_APITOKEN='<your_telegram_api_token>'
- go run main.go

## ChangeLog
- Reply to specific commands
- Upload and send gif
- Show keyboard and send the pressed key
- Send a link
- Send a position
- Download a video/audio using @vid bot
- Go routine for video download
- Reply to the @vid with the downloaded audio


## TODO
- Get the local playlist
- Error Handling during video download
- Go routine for the download method
- Mixing track

## Dependencies
- [Go Telegram Bot API](https://go-telegram-bot-api.dev/getting-started/index.html)
- [Youtube downloader](github.com/kkdai/youtube)
- [Mixing tracks](https://github.com/go-mix/mix)