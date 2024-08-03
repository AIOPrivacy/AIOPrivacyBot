package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"AIOPrivacyBot/functions/admins"
	"AIOPrivacyBot/functions/ai_chat"
	"AIOPrivacyBot/functions/ask"
	"AIOPrivacyBot/functions/check"
	"AIOPrivacyBot/functions/getid"
	"AIOPrivacyBot/functions/help"
	"AIOPrivacyBot/functions/num"
	"AIOPrivacyBot/functions/play"
	"AIOPrivacyBot/functions/status"
	"AIOPrivacyBot/functions/stringcalc"
	"AIOPrivacyBot/functions/view"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Config struct {
	Token                string   `json:"token"`
	SuperAdmins          []string `json:"super_admins"`
	SafeBrowsingAPIKey   string   `json:"safe_browsing_api_key"`
	TelegraphAccessToken string   `json:"telegraph_access_token"`
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

	// Initialize check package with SafeBrowsingAPIKey
	check.Init(config.SafeBrowsingAPIKey)

	// Initialize view package with Telegraph access token
	view.Init(config.TelegraphAccessToken)

	// 设置命令
	setBotCommands(bot)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			log.Printf("Received message from %s: %s", update.Message.From.UserName, update.Message.Text)
			processMessage(update.Message, bot)
		} else if update.InlineQuery != nil {
			processInlineQuery(update.InlineQuery, bot)
		} else if update.CallbackQuery != nil {
			admins.HandleCallbackQuery(update.CallbackQuery, bot)
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
		} else if command == "ask" && (message.Chat.IsPrivate() || strings.Contains(message.Text, fmt.Sprintf("@%s", botUsername))) {
			ask.HandleAskCommand(message, bot)
		} else if command == "getid" && (message.Chat.IsPrivate() || strings.Contains(message.Text, fmt.Sprintf("@%s", botUsername))) {
			getid.HandleGetIDCommand(message, bot, config.SuperAdmins)
		} else if command == "status" && (message.Chat.IsPrivate() || strings.Contains(message.Text, fmt.Sprintf("@%s", botUsername))) {
			status.HandleStatusCommand(message, bot)
		} else if command == "admins" && (message.Chat.IsGroup() || message.Chat.IsSuperGroup()) {
			admins.HandleAdminsCommand(message, bot)
		} else if command == "num" && (message.Chat.IsPrivate() || strings.Contains(message.Text, fmt.Sprintf("@%s", botUsername))) {
			num.HandleNumCommand(message, bot)
		} else if command == "string" && (message.Chat.IsPrivate() || strings.Contains(message.Text, fmt.Sprintf("@%s", botUsername))) {
			stringcalc.HandleStringCommand(message, bot)
		}
	} else if (message.Chat.IsGroup() || message.Chat.IsSuperGroup()) && isReplyToBot(message) && shouldTriggerResponse() {
		ai_chat.HandleAIChat(message, bot)
	}
}

func processInlineQuery(inlineQuery *tgbotapi.InlineQuery, bot *tgbotapi.BotAPI) {
	if strings.HasPrefix(inlineQuery.Query, "-view") {
		view.HandleViewCommand(inlineQuery, bot)
	} else if strings.HasPrefix(inlineQuery.Query, "-check") {
		check.HandleInlineQuery(inlineQuery, bot)
	}
}

func isReplyToBot(message *tgbotapi.Message) bool {
	if message.ReplyToMessage != nil && message.ReplyToMessage.From.UserName == botUsername {
		return true
	}
	return false
}

func shouldTriggerResponse() bool {
	rand.Seed(time.Now().UnixNano())
	randomValue := rand.Intn(100) + 1
	return randomValue > 0
}

func setBotCommands(bot *tgbotapi.BotAPI) {
	commands := []tgbotapi.BotCommand{
		{Command: "help", Description: "获取帮助信息"},
		{Command: "play", Description: "互动游玩"},
		{Command: "ask", Description: "提问AI"},
		{Command: "getid", Description: "获取ID"},
		{Command: "status", Description: "获取机器人状态"},
		{Command: "admins", Description: "召唤管理员"},
		{Command: "num", Description: "数字进制转换"},
		{Command: "string", Description: "字符串编码"},
	}

	config := tgbotapi.NewSetMyCommands(commands...)

	_, err := bot.Request(config)
	if err != nil {
		log.Fatalf("Error setting bot commands: %v", err)
	}

	log.Println("Bot commands set successfully")
}
