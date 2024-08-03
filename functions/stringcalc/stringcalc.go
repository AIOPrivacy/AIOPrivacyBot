package stringcalc

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"net/url"
	"strings"

	"AIOPrivacyBot/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/net/idna"
)

// HandleStringCommand 处理 /string 命令
func HandleStringCommand(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	input := message.CommandArguments()
	parts := strings.Fields(input)
	if len(parts) < 2 {
		utils.SendMessage(message.Chat.ID, "请输入有效的字符串和命令，例如 /string@AIOPrivacyBot -url example.com", message.MessageID, bot)
		return
	}

	command := strings.ToLower(parts[0])
	str := strings.Join(parts[1:], " ")

	var response string

	switch command {
	case "-url":
		response = handleURL(str)
	case "-crc":
		response = handleCRC(str)
	case "-base":
		response = handleBase(str)
	case "-unicode":
		response = handleUnicode(str)
	case "-ascii":
		response = handleASCII(str)
	case "-md5":
		response = handleMD5(str)
	case "-sha":
		response = handleSHA(str)
	case "-all":
		response = fmt.Sprintf("输入字符串：<b>%s</b>\n\n%s\n\n%s\n\n%s\n\n%s\n\n%s\n\n%s\n\n%s", str, handleURL(str), handleBase(str), handleUnicode(str), handleASCII(str), handleMD5(str), handleCRC(str), handleSHA(str))
	default:
		utils.SendMessage(message.Chat.ID, "未知命令，请输入正确的命令参数，例如 /string@bot -url example.com", message.MessageID, bot)
		return
	}

	utils.SendMessage(message.Chat.ID, response, message.MessageID, bot)
}

func handleURL(str string) string {
	encoded := url.QueryEscape(str)
	punycode, _ := idna.ToASCII(str)

	return fmt.Sprintf("<b>URL：</b>\nURLCode：<code>%s</code>\nPunycode：<code>%s</code>", encoded, punycode)
}

func handleCRC(str string) string {
	crc32q := crc32.MakeTable(0xD5828281)
	crc := crc32.Checksum([]byte(str), crc32q)
	return fmt.Sprintf("<b>CRC：</b><code>%d</code>", crc)
}

func handleBase(str string) string {
	base64Encoded := base64.StdEncoding.EncodeToString([]byte(str))
	base32Encoded := base32.StdEncoding.EncodeToString([]byte(str))
	base16Encoded := hex.EncodeToString([]byte(str))

	return fmt.Sprintf("<b>Base：</b>\nBase64：<code>%s</code>\nBase32：<code>%s</code>\nBase16：<code>%s</code>", base64Encoded, base32Encoded, base16Encoded)
}

func handleUnicode(str string) string {
	utf8Encoded := toUTF8Hex(str)
	utf16Encoded := toUTF16Hex(str)
	utf32Encoded := toUTF32Hex(str)

	return fmt.Sprintf("<b>Unicode：</b>\nUTF-8：<code>%s</code>\nUTF-16：<code>%s</code>\nUTF-32：<code>%s</code>", utf8Encoded, utf16Encoded, utf32Encoded)
}

func handleASCII(str string) string {
	var asciiCodes []string
	for _, r := range str {
		asciiCodes = append(asciiCodes, fmt.Sprintf("%d", r))
	}
	return fmt.Sprintf("<b>ASCII十进制：</b><code>%s</code>", strings.Join(asciiCodes, " "))
}

func handleMD5(str string) string {
	hash := md5.Sum([]byte(str))
	return fmt.Sprintf("<b>MD5：</b><code>%s</code>", hex.EncodeToString(hash[:]))
}

func handleSHA(str string) string {
	sha1Hash := sha1.Sum([]byte(str))
	sha256Hash := sha256.Sum256([]byte(str))
	sha512Hash := sha512.Sum512([]byte(str))

	return fmt.Sprintf("<b>SHA：</b>\nSHA1：<code>%s</code>\nSHA256：<code>%s</code>\nSHA512：<code>%s</code>", hex.EncodeToString(sha1Hash[:]), hex.EncodeToString(sha256Hash[:]), hex.EncodeToString(sha512Hash[:]))
}

func toUTF8Hex(s string) string {
	var result strings.Builder
	for _, r := range s {
		result.WriteString(fmt.Sprintf("%02X", r))
	}
	return result.String()
}

func toUTF16Hex(s string) string {
	var result strings.Builder
	for _, r := range s {
		result.WriteString(fmt.Sprintf("%04X", r))
	}
	return result.String()
}

func toUTF32Hex(s string) string {
	var result strings.Builder
	for _, r := range s {
		result.WriteString(fmt.Sprintf("%08X", r))
	}
	return result.String()
}
