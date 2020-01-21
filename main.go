// Package main provides ...
package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const token = ""

var registerChat = []int64{}

func listenCmd(bot *tgbotapi.BotAPI) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if !update.Message.IsCommand() {
			continue
		}

		var messageText string

		switch update.Message.Command() {
		case "chatid":
			messageText = fmt.Sprintf("%d", update.Message.Chat.ID)
		default:
			messageText = "unknown command"
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, messageText)
		bot.Send(msg)
	}
}

func main() {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = false
	log.Printf("Authorized on account %s", bot.Self.UserName)
	go listenCmd(bot)

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.Run()
}
