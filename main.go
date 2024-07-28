package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"AIOPrivacyBot/functions"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Config struct {
	Token       string   `json:"token"`
	SuperAdmins []string `json:"super_admins"`
}

func loadConfig() Config {
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatalf("Error opening config file: %v", err)
	}
	defer file.Close()

	byteValue, _ := ioutil.ReadAll(file)

	var config Config
	json.Unmarshal(byteValue, &config)

	return config
}

func main() {
	config := loadConfig()
	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	botID := bot.Self.ID

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			go handleUpdate(bot, update.Message, botID)
		}
	}
}

func handleUpdate(bot *tgbotapi.BotAPI, message *tgbotapi.Message, botID int64) {
	text := message.Text
	if strings.HasPrefix(text, "//") || (strings.HasPrefix(text, "/") && !strings.HasPrefix(text, "/$")) {
		if strings.HasPrefix(text, "/") && functions.IsAlphanumeric(text[1:]) {
			return
		}
		functions.HandleSlashCommands(bot, message)
	} else if strings.HasPrefix(text, "/$") {
		functions.HandleDollarCommand(bot, message)
	} else if message.ReplyToMessage != nil && message.ReplyToMessage.From.ID == botID {
		functions.HandleAIChat(bot, message)
	} else {
		// Here you can add more handling for other types of messages in the future
	}
}
