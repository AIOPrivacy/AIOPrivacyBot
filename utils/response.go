package utils

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// SendPhotoWithCaption 发送带有文字的图片
func SendPhotoWithCaption(chatID int64, photoPath, caption string, bot *tgbotapi.BotAPI) error {
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
	_, err = bot.Send(photo)
	if err != nil {
		log.Printf("Error sending photo: %v", err)
		return err
	}

	log.Printf("Photo sent successfully to chat ID %d", chatID)
	return nil
}

// SendMessage 发送文本消息
func SendMessage(chatID int64, text string, bot *tgbotapi.BotAPI) error {
	log.Printf("Sending message to chat ID %d: %s", chatID, text)
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Error sending message: %v", err)
		return err
	}
	log.Printf("Message sent successfully to chat ID %d", chatID)
	return nil
}
