package main

import (
	"io"
	"log"
	"os"

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

		link := linkToDownload(update)

		msg := downloadVideo(update, link)

		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
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
		msg.Text = "I understand /sayhi and /status."
	case "sayhi":
		msg.Text = "Hi :)"
	case "status":
		msg.Text = "I'm ok."
	case "play":
		msg.Text = update.Message.CommandArguments()
	default:
		msg.Text = "I don't know that command"
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

func SendLocation(update tgbotapi.Update, lat float64, long float64) tgbotapi.LocationConfig {

	return tgbotapi.NewLocation(update.Message.Chat.ID, lat, long)
}

func SendLink(update tgbotapi.Update, link string) tgbotapi.MessageConfig {

	return tgbotapi.NewMessage(update.Message.Chat.ID, link)
}

func linkToDownload(update tgbotapi.Update) string {

	if update.Message.ViaBot.IsBot && update.Message.ViaBot.UserName == "vid" {

		return update.Message.Text
	}
	return "Please use @vid nameofthevideo"

}

func downloadVideo(update tgbotapi.Update, link string) tgbotapi.MessageConfig {

	client := youtube.Client{Debug: true}

	videoID, err := youtube.ExtractVideoID(link)
	if err != nil {
		panic(err)
	}

	video, err := client.GetVideo(videoID)
	if err != nil {
		panic(err)
	}

	formats := video.Formats.WithAudioChannels() // only get videos with audio
	stream, _, err := client.GetStream(video, &formats[0])
	if err != nil {
		panic(err)
	}

	file, err := os.Create("video.mp3")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = io.Copy(file, stream)
	if err != nil {
		panic(err)
	}

	return tgbotapi.NewMessage(update.Message.Chat.ID, "Downloaded")

}
