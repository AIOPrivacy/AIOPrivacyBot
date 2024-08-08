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
		err := utils.SendMessage(message.Chat.ID,
			`<b>需要使用参数</b>
	以下演示都以<b>A回复B</b>模拟！

	<u>主动模式</u>
	<code>/play@AIOPrivacyBot -t xxxxx</code> 可以成功触发 <b>A xxxxx了 B！</b>
	<code>/play@AIOPrivacyBot -t xxxxx yyyyy</code> 可以成功触发 <b>A xxxxx B yyyyy</b>

	<u>被动模式</u>
	<code>/play@AIOPrivacyBot -p xxxxx</code> 可以成功触发 <b>A 被 B xxxxx了！</b>
	<code>/play@AIOPrivacyBot -p xxxxx yyyyy</code> 可以成功触发 <b>B xxxxx A yyyyy</b>

	<i>注意：可以使用英文 ' 或 " 包括发送内容来高于空格优先级，例如 <code>/play@AIOPrivacyBot -p "xx xxx" "yy yy y"</code></i>`,
			message.MessageID, bot)
		if err != nil {
			log.Printf("Error sending message: %v", err)
		}
		return
	}

	args := parseArguments(message.CommandArguments())

	if len(args) < 2 {
		err := utils.SendMessage(message.Chat.ID,
			`<b>需要更多参数</b>
	以下演示都以<b>A回复B</b>模拟！

	<u>主动模式</u>
	<code>/play@AIOPrivacyBot -t xxxxx</code> 可以成功触发 <b>A xxxxx了 B！</b>
	<code>/play@AIOPrivacyBot -t xxxxx yyyyy</code> 可以成功触发 <b>A xxxxx B yyyyy</b>

	<u>被动模式</u>
	<code>/play@AIOPrivacyBot -p xxxxx</code> 可以成功触发 <b>A 被 B xxxxx了！</b>
	<code>/play@AIOPrivacyBot -p xxxxx yyyyy</code> 可以成功触发 <b>B xxxxx A yyyyy</b>

	<i>注意：可以使用英文 ' 或 " 包括发送内容来高于空格优先级，例如 <code>/play@AIOPrivacyBot -p "xx xxx" "yy yy y"</code></i>`,
			message.MessageID, bot)

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
