package main

import (
	"errors"
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

				go DownloadRoutine(update, link)

			} else {

				CommandHandler(update, &t)

			}
		}

	}
}

func DownloadRoutine(update tgbotapi.Update, link string) {

	policeEmoji := "\xF0\x9F\x9A\xA8"
	clapboardEmoji := "\xF0\x9F\x8E\xAC"
	speakerEmoji := "\xF0\x9F\x94\x8A"
	envelopeEmoji := "\xE2\x9C\x89"
	sadface := "\xF0\x9F\x98\xA2"
	filename, err := DownloadVideo(update, link)
	if err != nil {
		sorryMessage := "\n" + sadface + "Sorry i cannot download this video, try with another version of the same song" + speakerEmoji
		sendTextMessage(update, policeEmoji+err.Error()+clapboardEmoji+sorryMessage, true)
		return
	}
	msg, err := SendAudio(update, filename)
	if err != nil {
		sendTextMessage(update, policeEmoji+err.Error()+speakerEmoji, true)
		return
	}
	msg.ReplyToMessageID = update.Message.MessageID
	if _, err := bot.Send(msg); err != nil {
		sendTextMessage(update, policeEmoji+err.Error()+envelopeEmoji, true)
		return

	}

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

	if video.Duration.Minutes() > 10 || video.Duration.Minutes() == 0 {
		fmt.Println(video.Duration.Minutes())
		return "", errors.New("this video is too long, i don't support streaming or playlist")
	}

	formats := video.Formats.WithAudioChannels().Type("audio/mp4") // only get videos with audio

	stream, _, err := client.GetStream(video, &formats[0])
	if err != nil {
		return "", err
	}
	regexTitle := regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(video.Title, "")
	fileTitle := regexTitle + "-" + strconv.Itoa(formats[0].AverageBitrate)
	if _, err := os.Stat("playlist/" + fileTitle + ".mp4"); err != nil {
		sendTextMessage(update, "\xF0\x9F\xA4\x98 In Download ==> "+fileTitle+" \xF0\x9F\x8E\xB5", true)

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

func CommandHandler(update tgbotapi.Update, times *time.Time) {
	// Extract the command from the Message.
	switch update.Message.Command() {

	case "help":
		sendTextMessage(update, "Use this bot via @vid <youtube video>, like: @vid free bird", false)

	case "start":
		fmt.Println(update.Message.Time().Unix() - times.Unix())
		if update.Message.Time().Unix()-times.Unix() > 10 {
			*times = update.Message.Time()
			err := SendGif("gif/demo.gif", update)
			if err != nil {
				sendTextMessage(update, err.Error(), true)
			}
			sendTextMessage(update, "Use this bot via @vid <youtube video>, like: @vid free bird", false)
		}

	default:
		sendTextMessage(update, "Use this bot via @vid <youtube video>, like: @vid free bird", false)
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

func SendGif(path string, update tgbotapi.Update) error {
	f_reader, err := os.Open(path)
	if err != nil {
		return err
	}

	file := tgbotapi.FileReader{
		Name:   "1.gif",
		Reader: f_reader,
	}

	msg := tgbotapi.NewAnimation(update.Message.Chat.ID, file)

	if _, err := bot.Send(msg); err != nil {
		return err
	}
	return nil

}

func SendAudio(update tgbotapi.Update, filename string) (tgbotapi.AudioConfig, error) {

	f_reader, err := os.Open("playlist/" + filename + ".mp4")
	if err != nil {
		return tgbotapi.AudioConfig{}, err
	}

	file := tgbotapi.FileReader{
		Name:   filename,
		Reader: f_reader,
	}

	return tgbotapi.NewAudio(update.Message.Chat.ID, file), nil

}

func SendLocation(update tgbotapi.Update, lat float64, long float64) tgbotapi.LocationConfig {

	return tgbotapi.NewLocation(update.Message.Chat.ID, lat, long)
}

func SendLink(update tgbotapi.Update, link string) tgbotapi.MessageConfig {

	return tgbotapi.NewMessage(update.Message.Chat.ID, link)
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
