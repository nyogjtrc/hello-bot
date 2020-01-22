// Package main provides ...
package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/spf13/viper"
)

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
		case "register":
			registerChat = append(registerChat, update.Message.Chat.ID)
			messageText = "done"
		default:
			messageText = "unknown command"
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, messageText)
		bot.Send(msg)
	}
}

func loadConfig() error {
	viper.SetEnvPrefix("bot")
	viper.BindEnv("token")
	fmt.Println("token", viper.GetString("token"))
	return nil
}

func main() {
	err := loadConfig()
	if err != nil {
		log.Panic(err)
	}

	bot, err := tgbotapi.NewBotAPI(viper.GetString("token"))
	if err != nil {
		log.Panic("new bot ", err)
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

	log.Fatal(r.Run())
}
