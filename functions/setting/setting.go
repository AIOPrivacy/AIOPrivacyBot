package setting

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Add the list of valid features
var availableFeatures = []string{
	"help", "play", "ask", "getid", "status", "admins",
	"num", "string", "curconv", "color", "setting",
}

// Features that cannot be disabled
var nonDisablableFeatures = []string{"setting"}

func IsFeatureEnabled(db *sql.DB, groupID int64, featureName string) bool {
	// Check if the feature is non-disablable
	for _, feature := range nonDisablableFeatures {
		if feature == featureName {
			return true
		}
	}

	var featureOffList string
	err := db.QueryRow("SELECT feature_off FROM group_setting WHERE groupid = ?", groupID).Scan(&featureOffList)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Error querying feature status: %v", err)
		return true // Default to enabled if error occurs
	}

	if featureOffList == "" {
		return true
	}

	disabledFeatures := strings.Split(featureOffList, ",")
	for _, feature := range disabledFeatures {
		if feature == featureName {
			return false
		}
	}

	return true
}

func HandleSettingCommand(db *sql.DB, message *tgbotapi.Message, bot *tgbotapi.BotAPI, superAdmins []string) {
	if !isAdmin(bot, message, superAdmins) {
		msg := tgbotapi.NewMessage(message.Chat.ID, "你没有权限修改设置")
		bot.Send(msg)
		return
	}

	args := strings.Fields(message.CommandArguments())
	if len(args) != 2 {
		msg := tgbotapi.NewMessage(message.Chat.ID, "用法: /setting <enable/disable> <feature_name>")
		bot.Send(msg)
		return
	}

	action := args[0]
	feature := args[1]

	// Validate the feature and action input
	if !validateInput(action, feature) {
		msg := tgbotapi.NewMessage(message.Chat.ID, "无效的输入，请检查命令格式和功能名称")
		bot.Send(msg)
		return
	}

	// Check if the feature is in the list of available features
	if !contains(availableFeatures, feature) {
		msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("功能 %s 不存在", feature))
		bot.Send(msg)
		return
	}

	if contains(nonDisablableFeatures, feature) && action == "disable" {
		msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("功能 %s 不能被禁用", feature))
		bot.Send(msg)
		return
	}

	var featureOffList string
	err := db.QueryRow("SELECT feature_off FROM group_setting WHERE groupid = ?", message.Chat.ID).Scan(&featureOffList)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Error querying feature status: %v", err)
		return
	}

	disabledFeatures := strings.Split(featureOffList, ",")

	if action == "enable" {
		disabledFeatures = remove(disabledFeatures, feature)
	} else if action == "disable" {
		if !contains(disabledFeatures, feature) {
			disabledFeatures = append(disabledFeatures, feature)
		}
	}

	featureOffList = strings.Join(disabledFeatures, ",")
	_, err = db.Exec("REPLACE INTO group_setting (groupid, feature_off) VALUES (?, ?)", message.Chat.ID, featureOffList)
	if err != nil {
		log.Printf("Error updating feature status: %v", err)
		return
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("功能 %s 已被 %s", feature, action))
	bot.Send(msg)
}

func isAdmin(bot *tgbotapi.BotAPI, message *tgbotapi.Message, superAdmins []string) bool {
	userID := fmt.Sprintf("%d", message.From.ID)
	for _, admin := range superAdmins {
		if admin == userID {
			return true
		}
	}

	admins, err := bot.GetChatAdministrators(tgbotapi.ChatAdministratorsConfig{
		ChatConfig: tgbotapi.ChatConfig{
			ChatID: message.Chat.ID,
		},
	})
	if err != nil {
		log.Printf("Error fetching chat administrators: %v", err)
		return false
	}

	for _, admin := range admins {
		if admin.User.ID == message.From.ID {
			return true
		}
	}

	return false
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func remove(slice []string, item string) []string {
	result := []string{}
	for _, v := range slice {
		if v != item {
			result = append(result, v)
		}
	}
	return result
}

// validateInput ensures that the action and feature are safe from SQL injection
func validateInput(action, feature string) bool {
	// Allowed actions are only "enable" and "disable"
	validActions := []string{"enable", "disable"}
	if !contains(validActions, action) {
		return false
	}

	// Ensure feature name only contains alphanumeric characters to prevent SQL injection
	re := regexp.MustCompile("^[a-zA-Z0-9]+$")
	return re.MatchString(feature)
}
