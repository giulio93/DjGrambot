package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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

		msg := SendGif("1.gif", update)
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

	msg := tgbotapi.NewAnimation(update.Message.Chat.ID, file)

	return msg
}
