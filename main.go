package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	youtube "github.com/kkdai/youtube/v2"
)

var bot *tgbotapi.BotAPI
var err error

func main() {
	bot, err = tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 90
	updates := bot.GetUpdatesChan(updateConfig)

	t := time.Now()
	// Let's go through each update that we're getting from Telegram.
	for update := range updates {

		if update.Message != nil { // ignore any non-Message updates

			if update.Message.ViaBot != nil && update.Message.ViaBot.UserName == "vid" {

				link := update.Message.Text

				err = DownloadRoutine(update, link)
				if err != nil {
					sendTextMessage(update, err.Error(), true)
				}

			} else {

				CommandHandler(update, &t)

			}
		}

	}
}

func DownloadRoutine(update tgbotapi.Update, link string) error {

	filename, err := DownloadVideo(update, link)
	if err != nil {
		return err
	}
	msg := SendAudio(update, filename)
	msg.ReplyToMessageID = update.Message.MessageID
	if _, err := bot.Send(msg); err != nil {
		return err
	}
	return nil

}

func CommandHandler(update tgbotapi.Update, times *time.Time) {
	// Extract the command from the Message.
	switch update.Message.Command() {

	case "help":
		sendTextMessage(update, "Use this bot via @vid <you tube video>", false)

	case "start":
		fmt.Println(update.Message.Time().Unix() - times.Unix())
		if update.Message.Time().Unix()-times.Unix() > 10 {
			*times = update.Message.Time()
			SendGif("gif/1.gif", update)
			sendTextMessage(update, "Oh Mona! Devi usare @vid per scaricare un video!", false)
		}

	default:
		sendTextMessage(update, "Oh Mona! Devi usare @vid per scaricare un video!", false)
	}

}

func sendTextMessage(update tgbotapi.Update, text string, replyToUpdate bool) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)

	if replyToUpdate {
		msg.ReplyToMessageID = update.Message.MessageID
	}

	if _, err := bot.Send(msg); err != nil {
		log.Panic(err)
	}

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

func SendGif(path string, update tgbotapi.Update) {
	f_reader, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	file := tgbotapi.FileReader{
		Name:   "1.gif",
		Reader: f_reader,
	}

	msg := tgbotapi.NewAnimation(update.Message.Chat.ID, file)

	if _, err := bot.Send(msg); err != nil {
		log.Panic(err)
	}
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

func DownloadVideo(update tgbotapi.Update, link string) (string, error) {

	client := youtube.Client{Debug: true}

	videoID, err := youtube.ExtractVideoID(link)
	if err != nil {
		return "", err
	}

	video, err := client.GetVideo(videoID)
	if err != nil {
		return "", err
	}

	formats := video.Formats.WithAudioChannels().Type("audio/mp4") // only get videos with audio

	stream, _, err := client.GetStream(video, &formats[0])
	if err != nil {
		return "", err
	}
	regexTitle := regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(video.Title, "")
	fileTitle := regexTitle + "-" + strconv.Itoa(formats[0].AverageBitrate)
	if _, err := os.Stat("playlist/" + fileTitle + ".mp4"); err != nil {

		file, err := os.Create("playlist/" + fileTitle + ".mp4")
		if err != nil {
			return "", err
		}
		defer file.Close()

		_, err = io.Copy(file, stream)
		if err != nil {
			return "", err
		}
	}

	return fileTitle, nil
}
