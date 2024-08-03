package num

import (
	"fmt"
	"log"
	"math/big"

	"AIOPrivacyBot/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// HandleNumCommand 处理 /num 命令
func HandleNumCommand(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	input := message.CommandArguments()
	var decimalValue = new(big.Int)

	// 自动识别进制并转换为十进制
	decimalValue, success := decimalValue.SetString(input, 0)

	if !success {
		log.Printf("Error parsing number: %s", input)
		utils.SendMessage(message.Chat.ID, "输入的数字格式不正确，请检查后重新输入。", message.MessageID, bot)
		return
	}

	// 转换为其他进制表示
	binaryValue := decimalValue.Text(2)
	octalValue := decimalValue.Text(8)
	hexValue := decimalValue.Text(16)

	// 构建回复消息
	replyText := fmt.Sprintf(
		"输入 (Input): <code>%s</code>\n十进制 (Decimal): <code>%s</code>\n二进制 (Binary): <code>%s</code>\n八进制 (Octal): <code>%s</code>\n十六进制 (Hex): <code>%s</code>",
		input, decimalValue.String(), binaryValue, octalValue, hexValue,
	)

	utils.SendMessage(message.Chat.ID, replyText, message.MessageID, bot)
}
