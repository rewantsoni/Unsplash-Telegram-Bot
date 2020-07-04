package main

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"os"
)

func main() {
	createbot()
}

//Creating the bot
func createbot() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("API_TOKEN"))
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("Hello, I am " + bot.Self.FirstName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		fmt.Println(err.Error())
	}
	getUpdates(bot, updates)
}

//Waiting for updates
func getUpdates(bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		go handleUpdate(bot, update)
	}
}
func sendImage(bot *tgbotapi.BotAPI, d data, ChatID int64) {
	req := newPhotoRequest(d)
	res, err := getResponse(req)
	if err != nil {
		var msg tgbotapi.Chattable
		if err.Error() == "couldn't find a image" {
			msg = tgbotapi.NewMessage(ChatID, "Couldn't find a image for "+d.Query)
		} else {
			msg = tgbotapi.NewMessage(ChatID, err.Error())
		}
		bot.Send(msg)
		return
	}

	if len(res) == 0 {
		msg := tgbotapi.NewMessage(ChatID, "Couldn't find a image for "+d.Query)
		bot.Send(msg)
		return
	}
	document := tgbotapi.NewDocumentShare(ChatID, res[0].Urls.Full)
	document.Caption = "By " + res[0].User.FirstName + " " + res[0].User.LastName + " On Unsplash\n" + res[0].User.Links.Html

	document.ReplyMarkup = newInlineKeyboard(data{Query: d.Query, Random: d.Random})
	_, err = bot.Send(document)
	if err != nil {
		panic(err.Error())
	}
}

//Handling Updates
func handleUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	var d data
	if update.CallbackQuery != nil {
		_ = json.Unmarshal([]byte(update.CallbackQuery.Data), &d)
		sendImage(bot, d, update.CallbackQuery.Message.Chat.ID)
		newCallbackConfig := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
		_, _ = bot.AnswerCallbackQuery(newCallbackConfig)
		return
	}
	if update.Message == nil || update.Message.Text == "" {
		return
	}
	if update.Message.Chat.IsPrivate() {
		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "random":
				d.Random = true
			case "start":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Welcome to Unsplash Bot")
				bot.Send(msg)
				return
			}
		} else {
			d.Query = update.Message.Text
		}
		sendImage(bot, d, update.Message.Chat.ID)
	}
}
