package play

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"AIOPrivacyBot/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandlePlayCommand(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	if message.Chat.IsPrivate() {
		err := utils.SendMessage(message.Chat.ID, "此命令只可用于群组", message.MessageID, bot)
		if err != nil {
			log.Printf("Error sending message: %v", err)
		}
		return
	}

	if message.CommandArguments() == "" {
		err := utils.SendMessage(message.Chat.ID, "需要使用参数", message.MessageID, bot)
		if err != nil {
			log.Printf("Error sending message: %v", err)
		}
		return
	}

	args := parseArguments(message.CommandArguments())

	if len(args) < 2 {
		err := utils.SendMessage(message.Chat.ID, "需要更多参数", message.MessageID, bot)
		if err != nil {
			log.Printf("Error sending message: %v", err)
		}
		return
	}

	senderMention := fmt.Sprintf("<a href=\"tg://user?id=%d\">%s %s</a>", message.From.ID, message.From.FirstName, message.From.LastName)
	targetMention := senderMention
	if message.ReplyToMessage != nil {
		targetMention = fmt.Sprintf("<a href=\"tg://user?id=%d\">%s %s</a>", message.ReplyToMessage.From.ID, message.ReplyToMessage.From.FirstName, message.ReplyToMessage.From.LastName)
	}

	commandType := args[0]
	action := args[1]
	extra := ""
	if len(args) > 2 {
		extra = args[2]
	}

	var response string
	if commandType == "-t" {
		if extra == "" {
			response = fmt.Sprintf("%s %s了 %s！", senderMention, action, targetMention)
		} else {
			response = fmt.Sprintf("%s %s %s %s", senderMention, action, targetMention, extra)
		}
	} else if commandType == "-p" {
		if extra == "" {
			response = fmt.Sprintf("%s 被 %s %s了！", senderMention, targetMention, action)
		} else {
			response = fmt.Sprintf("%s %s %s %s", targetMention, action, senderMention, extra)
		}
	}

	err := utils.SendMessage(message.Chat.ID, response, message.MessageID, bot)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

func parseArguments(command string) []string {
	re := regexp.MustCompile(`'[^']*'|"[^"]*"|\S+`)
	matches := re.FindAllString(command, -1)
	for i := range matches {
		matches[i] = strings.Trim(matches[i], `"'`)
	}
	return matches
}
