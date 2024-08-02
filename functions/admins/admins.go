package admins

import (
	"fmt"
	"log"
	"strings"

	"AIOPrivacyBot/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleAdminsCommand(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	if message.Chat.IsPrivate() {
		return // 不支持私聊使用
	}

	chatID := message.Chat.ID
	msgText := `⚠️ *警告：您即将召唤本群所有管理员，请确认您要这样做！（如果故意随意@管理员可能导致封禁）*`

	buttons := []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("确认", fmt.Sprintf("confirm_admins:%d", message.From.ID)), // 添加发送者用户ID作为确认的标识
		tgbotapi.NewInlineKeyboardButtonData("取消", fmt.Sprintf("cancel_admins:%d", message.From.ID)),  // 添加发送者用户ID作为取消的标识
	}

	err := utils.SendMarkdownMessageWithInlineKeyboard(chatID, message.MessageID, msgText, buttons, bot)
	if err != nil {
		log.Printf("Error sending admins command message: %v", err)
	}
}

func HandleCallbackQuery(callbackQuery *tgbotapi.CallbackQuery, bot *tgbotapi.BotAPI) {
	chatID := callbackQuery.Message.Chat.ID
	callbackID := callbackQuery.ID
	messageID := callbackQuery.Message.MessageID

	// 解析确认的用户ID
	dataParts := strings.Split(callbackQuery.Data, ":")
	action := dataParts[0]

	if len(dataParts) < 2 {
		log.Printf("Invalid callback data: %s", callbackQuery.Data)
		return
	}

	senderID := dataParts[1] // 获取发送者用户ID

	if fmt.Sprintf("%d", callbackQuery.From.ID) != senderID {
		// 如果不是消息的发送者，则忽略
		callback := tgbotapi.NewCallback(callbackID, "您没有权限执行此操作。")
		_, err := bot.Request(callback)
		if err != nil {
			log.Printf("Error sending permission error callback: %v", err)
		}
		return
	}

	if action == "confirm_admins" {
		admins, err := bot.GetChatAdministrators(tgbotapi.ChatAdministratorsConfig{
			ChatConfig: tgbotapi.ChatConfig{ChatID: chatID},
		})
		if err != nil {
			log.Printf("Error getting chat administrators: %v", err)
			return
		}

		adminMentions := ""
		for _, admin := range admins {
			adminMentions += fmt.Sprintf("<a href=\"tg://user?id=%d\">%s %s</a> . ", admin.User.ID, admin.User.FirstName, admin.User.LastName)
		}

		msgText := "⚠️ 召唤本群所有管理员：" + adminMentions

		// 发送 @ 管理员的消息
		err = utils.SendMessage(chatID, msgText, 0, bot)
		if err != nil {
			log.Printf("Error sending confirm command message: %v", err)
		}

		// 修改确认消息为
		editMsg := tgbotapi.NewEditMessageText(chatID, messageID, "消息已确认，已经召唤所有管理员，无法取消或修改，请等待管理员回复")
		editMsg.ParseMode = "Markdown"
		_, err = bot.Send(editMsg)
		if err != nil {
			log.Printf("Error editing confirmed command message: %v", err)
		}
	} else if action == "cancel_admins" {
		msgText := "已取消召唤所有管理员"

		editMsg := tgbotapi.NewEditMessageText(chatID, messageID, msgText)
		_, err := bot.Send(editMsg)
		if err != nil {
			log.Printf("Error sending cancel command message: %v", err)
		}
	}

	callback := tgbotapi.NewCallback(callbackID, "")
	_, err := bot.Request(callback)
	if err != nil {
		log.Printf("Error sending callback: %v", err)
	}
}
