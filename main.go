package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"AIOPrivacyBot/functions/help"
	"AIOPrivacyBot/functions/play"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Config struct {
	Token       string   `json:"token"`
	SuperAdmins []string `json:"super_admins"`
}

var (
	botUsername string
	config      Config
)

func main() {
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatalf("Error opening config file: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatalf("Error decoding config file: %v", err)
	}

	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		log.Fatalf("Error creating new bot: %v", err)
	}

	botUsername = bot.Self.UserName
	log.Printf("Authorized on account %s", botUsername)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			log.Printf("Received message from %s: %s", update.Message.From.UserName, update.Message.Text)
			processMessage(update.Message, bot)
		}
	}
}

func processMessage(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	log.Printf("Processing message from %s: %s", message.From.UserName, message.Text)

	if message.IsCommand() {
		command := message.Command()
		if command == "help" && (message.Chat.IsPrivate() || strings.Contains(message.Text, fmt.Sprintf("@%s", botUsername))) {
			help.SendHelpMessage(message, bot)
		} else if command == "play" && strings.Contains(message.Text, fmt.Sprintf("@%s", botUsername)) {
			play.HandlePlayCommand(message, bot)
		}
	}
}
