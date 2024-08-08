package color

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"AIOPrivacyBot/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/image/colornames"
)

// HandleColorCommand 处理 /color 命令
func HandleColorCommand(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	args := strings.Fields(message.CommandArguments())
	if len(args) == 0 {
		utils.SendMessage(message.Chat.ID, "请提供一个颜色名称、RGB 值或十六进制颜色值。", message.MessageID, bot)
		return
	}

	var r, g, b int
	var colorName string

	if len(args) == 1 {
		colorName = args[0]
		if strings.HasPrefix(colorName, "#") {
			if matched, _ := regexp.MatchString(`^#[0-9a-fA-F]{6}$`, colorName); !matched {
				utils.SendMessage(message.Chat.ID, "无效的十六进制颜色值。", message.MessageID, bot)
				return
			}
			r, g, b = hexToRGB(colorName)
			colorName = "unnamed"
		} else {
			r, g, b = getColorFromName(colorName)
			if r == -1 && g == -1 && b == -1 {
				utils.SendMessage(message.Chat.ID, "无效的颜色名称。", message.MessageID, bot)
				return
			}
		}
	} else if len(args) == 3 {
		// 解析 RGB 值
		var err error
		r, err = strconv.Atoi(args[0])
		if err != nil {
			utils.SendMessage(message.Chat.ID, "无效的 RGB 值。", message.MessageID, bot)
			return
		}
		g, err = strconv.Atoi(args[1])
		if err != nil {
			utils.SendMessage(message.Chat.ID, "无效的 RGB 值。", message.MessageID, bot)
			return
		}
		b, err = strconv.Atoi(args[2])
		if err != nil {
			utils.SendMessage(message.Chat.ID, "无效的 RGB 值。", message.MessageID, bot)
			return
		}
		colorName = "unnamed"
	} else {
		utils.SendMessage(message.Chat.ID, "请提供一个有效的颜色名称、RGB 值或十六进制颜色值。", message.MessageID, bot)
		return
	}

	sendColorResponse(message.Chat.ID, message.MessageID, colorName, r, g, b, bot)
}

func getColorFromName(name string) (int, int, int) {
	col, ok := colornames.Map[strings.ToLower(name)]
	if !ok {
		return -1, -1, -1 // 返回无效的 RGB 值
	}
	r, g, b, _ := col.RGBA()
	return int(r / 257), int(g / 257), int(b / 257) // 将 16 位的值转换为 8 位
}

func hexToRGB(hex string) (int, int, int) {
	hex = strings.TrimPrefix(hex, "#")
	r, _ := strconv.ParseInt(hex[0:2], 16, 0)
	g, _ := strconv.ParseInt(hex[2:4], 16, 0)
	b, _ := strconv.ParseInt(hex[4:6], 16, 0)
	return int(r), int(g), int(b)
}

func sendColorResponse(chatID int64, messageID int, name string, r, g, b int, bot *tgbotapi.BotAPI) {
	hex := fmt.Sprintf("#%02x%02x%02x", r, g, b)
	rgbPercent := fmt.Sprintf("(%.2f%%, %.2f%%, %.2f%%)", float64(r)/255*100, float64(g)/255*100, float64(b)/255*100)

	// 获取推荐颜色
	recommendedColors := getRecommendedColors(r, g, b)

	text := fmt.Sprintf("名称 (Name): <code>%s</code>\nRGB: <code>(%d, %d, %d)</code>\n十六进制 (Hex): <code>%s</code>\nRGB 百分比 (RGB Percent): <code>%s</code>\n\n推荐颜色:\n%s", name, r, g, b, hex, rgbPercent, recommendedColors)

	imagePath := fmt.Sprintf("/tmp/color_%d_%d_%d.png", r, g, b)
	err := createColorImage(r, g, b, imagePath)
	if err != nil {
		log.Printf("Error creating color image: %v", err)
		utils.SendMessage(chatID, "生成颜色图片时出错。", messageID, bot)
		return
	}

	err = utils.SendPhotoWithCaption(chatID, messageID, imagePath, text, bot)
	if err != nil {
		log.Printf("Error sending photo: %v", err)
		utils.SendMessage(chatID, "发送颜色图片时出错。", messageID, bot)
		return
	}

	// 删除临时图片文件
	err = os.Remove(imagePath)
	if err != nil {
		log.Printf("Error removing temp image file: %v", err)
	}
}

func getRecommendedColors(r, g, b int) string {
	var recommendations string
	// 生成不同亮度的推荐颜色
	for i := 1; i <= 3; i++ {
		factor := float64(i) / 4.0
		newR := int(float64(r) * factor)
		newG := int(float64(g) * factor)
		newB := int(float64(b) * factor)
		hex := fmt.Sprintf("#%02x%02x%02x", newR, newG, newB)
		recommendations += fmt.Sprintf("<code>%s</code>\n", hex)
	}
	return recommendations
}

func createColorImage(r, g, b int, path string) error {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	fillColor := color.RGBA{uint8(r), uint8(g), uint8(b), 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{fillColor}, image.Point{}, draw.Src)

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, img)
}
