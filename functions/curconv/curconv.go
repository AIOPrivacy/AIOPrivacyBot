package curconv

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// API URL 模板
const apiURL = "https://www.mastercard.com/settlement/currencyrate/conversion-rate?fxDate=0000-00-00&transCurr=%s&crdhldBillCurr=%s&bankFee=0&transAmt=%s"

// HandleCurconvCommand 处理汇率转换命令
func HandleCurconvCommand(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	args := strings.Fields(message.CommandArguments())

	if len(args) < 2 {
		reply := "用法: <b>/curconv 货币源 货币目标 [金额]</b>"
		SendMessage(message.Chat.ID, reply, message.MessageID, bot)
		return
	}

	// 校验输入，只允许英文字母、数字、空格和小数点
	validInput := regexp.MustCompile(`^[a-zA-Z0-9\s.]+$`).MatchString
	for _, arg := range args {
		if !validInput(arg) {
			reply := "无效的输入，请使用: <b>/curconv 货币源 货币目标 [金额]</b>"
			SendMessage(message.Chat.ID, reply, message.MessageID, bot)
			return
		}
	}

	sourceCurrency := args[0]
	targetCurrency := args[1]
	amount := "100" // 默认金额

	if len(args) == 3 {
		// 尝试解析 amount 参数为浮点数
		if _, err := strconv.ParseFloat(args[2], 64); err != nil {
			reply := "无效的金额，请输入一个有效的数字。用法: <b>/curconv 货币源 货币目标 [金额]</b>"
			SendMessage(message.Chat.ID, reply, message.MessageID, bot)
			return
		}
		amount = args[2]
	}

	log.Printf("Fetching conversion rate for %s to %s with amount %s", sourceCurrency, targetCurrency, amount)
	url := fmt.Sprintf(apiURL, sourceCurrency, targetCurrency, amount)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error fetching conversion rate: %v", err)
		SendMessage(message.Chat.ID, "获取汇率失败，请稍后再试。", message.MessageID, bot)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("API request failed with status code: %d", resp.StatusCode)
		SendMessage(message.Chat.ID, "API请求失败，请稍后再试。", message.MessageID, bot)
		return
	}

	var result struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Date        string `json:"date"`
		Type        string `json:"type"`
		Data        struct {
			ConversionRate float64 `json:"conversionRate"`
			CrdhldBillAmt  float64 `json:"crdhldBillAmt"`
			FxDate         string  `json:"fxDate"`
			TransCurr      string  `json:"transCurr"`
			CrdhldBillCurr string  `json:"crdhldBillCurr"`
			TransAmt       float64 `json:"transAmt"`
			ErrorCode      string  `json:"errorCode"`
			ErrorMessage   string  `json:"errorMessage"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("Error decoding API response: %v", err)
		SendMessage(message.Chat.ID, "解析汇率数据失败，请稍后再试。", message.MessageID, bot)
		return
	}

	// 处理API错误响应
	if result.Type == "error" {
		log.Printf("API error: %s", result.Data.ErrorMessage)
		reply := "货币不可用，请在 <a href=\"https://www.mastercard.com/settlement/currencyrate/settlement-currencies\">这里</a> 选择有效的货币"
		SendMessage(message.Chat.ID, reply, message.MessageID, bot)
		return
	}

	var reply string
	if len(args) == 3 {
		reply = fmt.Sprintf(
			"<b>%s</b> （UTC）<b>\n%s</b> - <b>%s</b> 的汇率为 <b>%.6f</b>\n若您的金额为 <b>%.2f %s</b>，那么你可以兑换 <b>%.2f %s</b>",
			result.Date, result.Data.TransCurr, result.Data.CrdhldBillCurr,
			result.Data.ConversionRate, result.Data.TransAmt, result.Data.TransCurr,
			result.Data.CrdhldBillAmt, result.Data.CrdhldBillCurr,
		)
	} else {
		reply = fmt.Sprintf(
			"<b>%s</b> （UTC）<b>\n%s</b> - <b>%s</b> 的汇率为 <b>%.6f</b>",
			result.Date, result.Data.TransCurr, result.Data.CrdhldBillCurr,
			result.Data.ConversionRate,
		)
	}

	log.Printf("Sending response: %s", reply)
	SendMessage(message.Chat.ID, reply, message.MessageID, bot)
}

// SendMessage 发送文本消息
func SendMessage(chatID int64, text string, messageID int, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "HTML"
	msg.ReplyToMessageID = messageID // 设置回复消息ID
	_, err := bot.Send(msg)
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}
}
