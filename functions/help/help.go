package help

import (
	"log"

	"AIOPrivacyBot/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SendHelpMessage(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	chatID := message.Chat.ID
	photoPath := "functions/help/help_image.png"
	photoCaption := "这是帮助图片\n该机器人正在测试中"

	// 使用utils.SendPhotoWithCaption发送带有文字的图片
	err := utils.SendPhotoWithCaption(chatID, photoPath, photoCaption, bot)
	if err != nil {
		log.Printf("Error sending help image: %v", err)
	}
}
