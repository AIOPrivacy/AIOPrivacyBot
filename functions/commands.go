package functions

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleSlashCommands(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	text := strings.TrimPrefix(message.Text, "/") // Remove the initial "/" or "//"
	text = strings.TrimPrefix(text, "/")          // Remove the second "/" if present

	var response string
	var replyToName string

	if message.ReplyToMessage != nil {
		replyToName = fmt.Sprintf("<a href=\"tg://user?id=%d\">%s %s</a>", message.ReplyToMessage.From.ID, message.ReplyToMessage.From.FirstName, message.ReplyToMessage.From.LastName)
	} else {
		replyToName = "自己"
	}

	parts := strings.SplitN(text, " ", 2)
	if len(parts) == 2 {
		response = fmt.Sprintf("<a href=\"tg://user?id=%d\">%s %s</a> %s %s %s", message.From.ID, message.From.FirstName, message.From.LastName, parts[0], replyToName, parts[1])
	} else {
		response = fmt.Sprintf("<a href=\"tg://user?id=%d\">%s %s</a> %s了 %s", message.From.ID, message.From.FirstName, message.From.LastName, parts[0], replyToName)
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, response)
	msg.ParseMode = "HTML"
	bot.Send(msg)
}

func HandleDollarCommand(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	text := strings.TrimPrefix(message.Text, "/$") // Remove the initial "/$"
	var response string

	parts := strings.SplitN(text, " ", 2)
	if len(parts) == 2 {
		if message.ReplyToMessage != nil {
			response = fmt.Sprintf("<a href=\"tg://user?id=%d\">%s %s</a> %s <a href=\"tg://user?id=%d\">%s %s</a> %s", message.ReplyToMessage.From.ID, message.ReplyToMessage.From.FirstName, message.ReplyToMessage.From.LastName, parts[0], message.From.ID, message.From.FirstName, message.From.LastName, parts[1])
		} else {
			response = fmt.Sprintf("<a href=\"tg://user?id=%d\">%s %s</a> %s 自己 %s", message.From.ID, message.From.FirstName, message.From.LastName, parts[0], parts[1])
		}
	} else {
		if message.ReplyToMessage != nil {
			response = fmt.Sprintf("<a href=\"tg://user?id=%d\">%s %s</a> 被 <a href=\"tg://user?id=%d\">%s %s</a> %s了", message.From.ID, message.From.FirstName, message.From.LastName, message.ReplyToMessage.From.ID, message.ReplyToMessage.From.FirstName, message.ReplyToMessage.From.LastName, text)
		} else {
			response = fmt.Sprintf("<a href=\"tg://user?id=%d\">%s %s</a> 被 自己 %s了", message.From.ID, message.From.FirstName, message.From.LastName, text)
		}
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, response)
	msg.ParseMode = "HTML"
	bot.Send(msg)
}

func IsAlphanumeric(s string) bool {
	for _, r := range s {
		if !('a' <= r && r <= 'z' || 'A' <= r && r <= 'Z' || '0' <= r && r <= '9' || r == '@' || r == '_') {
			return false
		}
	}
	return true
}
