package getid

import (
	"fmt"

	"AIOPrivacyBot/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// HandleGetIDCommand handles the /getid command.
func HandleGetIDCommand(message *tgbotapi.Message, bot *tgbotapi.BotAPI, configSuperAdmins []string) {
	var response string
	user := message.From

	isSuperAdmin := "否"
	for _, admin := range configSuperAdmins {
		if admin == fmt.Sprintf("%d", user.ID) {
			isSuperAdmin = "是"
			break
		}
	}

	if message.Chat.IsPrivate() {
		response = fmt.Sprintf(`<b>个人ID信息：</b>
个人用户名：%s
个人昵称：%s %s
个人ID：%d
超级管理员：%s`,
			user.UserName, user.FirstName, user.LastName, user.ID, isSuperAdmin)
	} else {
		chat := message.Chat
		response = fmt.Sprintf(`<b>群聊ID信息：</b>
群组名称：%s
群组类型：%s
群组ID：%d

<b>个人ID信息：</b>
个人用户名：%s
个人昵称：%s %s
个人ID：%d
超级管理员：%s`,
			chat.Title, chat.Type, chat.ID,
			user.UserName, user.FirstName, user.LastName, user.ID, isSuperAdmin)
	}

	utils.SendMessage(message.Chat.ID, response, bot)
}
