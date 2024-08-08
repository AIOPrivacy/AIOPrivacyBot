package utils

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// SendPhotoWithCaption 发送带有文字的图片
func SendPhotoWithCaption(chatID int64, messageID int, photoPath, caption string, bot *tgbotapi.BotAPI) error {
	log.Printf("Sending photo with caption to chat ID %d: %s", chatID, caption)

	photoFile, err := os.Open(photoPath)
	if err != nil {
		log.Printf("Error opening photo: %v", err)
		return err
	}
	defer photoFile.Close()

	photo := tgbotapi.NewPhoto(chatID, tgbotapi.FileReader{
		Name:   photoPath,
		Reader: photoFile,
	})
	photo.Caption = caption
	photo.ReplyToMessageID = messageID // 设置回复消息ID
	_, err = bot.Send(photo)
	if err != nil {
		log.Printf("Error sending photo: %v", err)
		return err
	}

	log.Printf("Photo sent successfully to chat ID %d", chatID)
	return nil
}

// SendMessage 发送文本消息
func SendMessage(chatID int64, text string, messageID int, bot *tgbotapi.BotAPI) error {
	log.Printf("Sending message to chat ID %d: %s", chatID, text)
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	msg.ReplyToMessageID = messageID // 设置回复消息ID
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Error sending message: %v", err)
		return err
	}
	log.Printf("Message sent successfully to chat ID %d", chatID)
	return nil
}

// SendMarkdownMessage 发送 Markdown 格式的文本消息，并回复到用户
func SendMarkdownMessage(chatID int64, messageID int, text string, bot *tgbotapi.BotAPI) error {
	log.Printf("Sending Markdown message to chat ID %d: %s", chatID, text)
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"       // 使用 Markdown 解析模式
	msg.ReplyToMessageID = messageID // 回复到原始消息

	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Error sending Markdown message: %v", err)
		return err
	}
	log.Printf("Markdown message sent successfully to chat ID %d", chatID)
	return nil
}

// SendMarkdownMessageWithInlineKeyboard 发送带有内联键盘的 Markdown 格式的消息
func SendMarkdownMessageWithInlineKeyboard(chatID int64, messageID int, text string, buttons []tgbotapi.InlineKeyboardButton, bot *tgbotapi.BotAPI) error {
	log.Printf("Sending Markdown message with inline keyboard to chat ID %d: %s", chatID, text)
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyToMessageID = messageID // 回复到原始消息

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(buttons...),
	)
	msg.ReplyMarkup = keyboard

	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Error sending Markdown message with inline keyboard: %v", err)
		return err
	}
	log.Printf("Markdown message with inline keyboard sent successfully to chat ID %d", chatID)
	return nil
}

// SendPlainTextMessage 发送纯文本消息
func SendPlainTextMessage(chatID int64, text string, messageID int, bot *tgbotapi.BotAPI) error {
	log.Printf("Sending plain text message to chat ID %d: %s", chatID, text)
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = ""               // 不使用任何解析模式
	msg.ReplyToMessageID = messageID // 回复到原始消息
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Error sending plain text message: %v", err)
		return err
	}
	log.Printf("Plain text message sent successfully to chat ID %d", chatID)
	return nil
}
