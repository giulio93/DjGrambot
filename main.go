package main

import (
	"io"
	"log"
	"os"
	"regexp"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	youtube "github.com/kkdai/youtube/v2"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := bot.GetUpdatesChan(updateConfig)

	// Let's go through each update that we're getting from Telegram.
	for update := range updates {

		if update.Message == nil { // ignore any non-Message updates
			continue
		}

		if update.Message.ViaBot != nil && update.Message.ViaBot.UserName == "vid" {

			link := update.Message.Text

			filename := DownloadVideo(update, link)

			msg := SendAudio(update, filename)
			if _, err := bot.Send(msg); err != nil {
				log.Panic(err)
			}

		} else {
			msg := CommandHandler(update)

			if _, err := bot.Send(msg); err != nil {
				log.Panic(err)
			}

		}

	}
}

func CommandHandler(update tgbotapi.Update) tgbotapi.MessageConfig {

	// Create a new MessageConfig. We don't have text yet,
	// so we leave it empty.
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	// Extract the command from the Message.
	switch update.Message.Command() {
	case "help":
		msg.Text = "Use this bot via @vid"
	default:
		msg.Text = "Oh Mona! Devi usare @vid per scaricare un video!"
	}

	return msg

}

func NumericKeyboard(update tgbotapi.Update) tgbotapi.MessageConfig {
	var numericKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("1"),
			tgbotapi.NewKeyboardButton("2"),
			tgbotapi.NewKeyboardButton("3"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("4"),
			tgbotapi.NewKeyboardButton("5"),
			tgbotapi.NewKeyboardButton("6"),
		),
	)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

	switch update.Message.Text {
	case "open":
		msg.ReplyMarkup = numericKeyboard
	case "close":
		msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	}

	return msg
}

func SendGif(path string, update tgbotapi.Update) tgbotapi.AnimationConfig {
	f_reader, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	file := tgbotapi.FileReader{
		Name:   "1.gif",
		Reader: f_reader,
	}

	return tgbotapi.NewAnimation(update.Message.Chat.ID, file)

}

func SendAudio(update tgbotapi.Update, filename string) tgbotapi.AudioConfig {

	f_reader, err := os.Open("playlist/" + filename + ".mp4")
	if err != nil {
		panic(err)
	}

	file := tgbotapi.FileReader{
		Name:   filename,
		Reader: f_reader,
	}

	return tgbotapi.NewAudio(update.Message.Chat.ID, file)

}

func SendLocation(update tgbotapi.Update, lat float64, long float64) tgbotapi.LocationConfig {

	return tgbotapi.NewLocation(update.Message.Chat.ID, lat, long)
}

func SendLink(update tgbotapi.Update, link string) tgbotapi.MessageConfig {

	return tgbotapi.NewMessage(update.Message.Chat.ID, link)
}

func DownloadVideo(update tgbotapi.Update, link string) string {

	client := youtube.Client{Debug: true}

	videoID, err := youtube.ExtractVideoID(link)
	if err != nil {
		panic(err)
	}

	video, err := client.GetVideo(videoID)
	if err != nil {
		panic(err)
	}

	formats := video.Formats.WithAudioChannels().Type("audio/mp4") // only get videos with audio

	stream, _, err := client.GetStream(video, &formats[0])
	if err != nil {
		panic(err)
	}
	regexTitle := regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(video.Title, "")
	fileTitle := regexTitle + "-" + strconv.Itoa(formats[0].AverageBitrate)
	file, err := os.Create("playlist/" + fileTitle + ".mp4")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = io.Copy(file, stream)
	if err != nil {
		panic(err)
	}

	return fileTitle
}
