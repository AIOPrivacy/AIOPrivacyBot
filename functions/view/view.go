package view

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/playwright-community/playwright-go"
)

const telegraphAPI = "https://api.telegra.ph/createPage"

var accessToken string

type CreatePageResponse struct {
	OK     bool   `json:"ok"`
	Result Result `json:"result"`
}

type Result struct {
	URL string `json:"url"`
}

type PageData struct {
	Title string `json:"title"`
	Node  string `json:"node"`
}

func Init(tgToken string) {
	accessToken = tgToken
}

func HandleViewCommand(inlineQuery *tgbotapi.InlineQuery, bot *tgbotapi.BotAPI) {
	query := strings.TrimSpace(inlineQuery.Query)
	if !strings.HasPrefix(query, "-view") {
		return
	}

	url := strings.TrimPrefix(query, "-view")
	url = strings.TrimSpace(url)
	if url == "" {
		return
	}

	pageData, err := fetchPageData(url)
	if err != nil {
		log.Printf("Error fetching page data: %v", err)
		return
	}

	if pageData.Title == "" || pageData.Node == "" {
		results := []interface{}{
			tgbotapi.NewInlineQueryResultArticleHTML(
				inlineQuery.ID,
				"链接不支持",
				"不支持此链接，请使用 -check功能",
			),
		}
		inlineConf := tgbotapi.InlineConfig{
			InlineQueryID: inlineQuery.ID,
			Results:       results,
			IsPersonal:    true,
		}
		if _, err := bot.Request(inlineConf); err != nil {
			log.Printf("Error sending inline query response: %v", err)
		}
		return
	}

	telegraphURL, err := postToTelegraph(pageData)
	if err != nil {
		log.Printf("Error posting to Telegraph: %v", err)
		return
	}

	results := []interface{}{
		tgbotapi.NewInlineQueryResultArticleHTML(
			inlineQuery.ID,
			"URL Clean & View",
			fmt.Sprintf("我通过 @AIOPrivacyBot 分享了文章：《%s》\n\n查看完整内容请点击：<a href=\"%s\">这里</a>", pageData.Title, telegraphURL),
		),
	}

	inlineConf := tgbotapi.InlineConfig{
		InlineQueryID: inlineQuery.ID,
		Results:       results,
		IsPersonal:    true,
	}

	if _, err := bot.Request(inlineConf); err != nil {
		log.Printf("Error sending inline query response: %v", err)
	}
}

func fetchPageData(targetURL string) (PageData, error) {
	// 自动安装浏览器依赖
	if err := playwright.Install(); err != nil {
		return PageData{}, fmt.Errorf("could not install browsers: %w", err)
	}

	pw, err := playwright.Run()
	if err != nil {
		return PageData{}, fmt.Errorf("could not start playwright: %w", err)
	}
	defer pw.Stop()

	browser, err := pw.Chromium.Launch()
	if err != nil {
		return PageData{}, fmt.Errorf("could not launch browser: %w", err)
	}
	defer browser.Close()

	page, err := browser.NewPage()
	if err != nil {
		return PageData{}, fmt.Errorf("could not create page: %w", err)
	}

	if _, err = page.Goto(targetURL); err != nil {
		return PageData{}, fmt.Errorf("could not goto: %w", err)
	}

	jsFilePath := "./functions/view/export.js"
	jsCode, err := os.ReadFile(jsFilePath)
	if err != nil {
		return PageData{}, fmt.Errorf("could not read JavaScript file: %w", err)
	}

	result, err := page.Evaluate(string(jsCode))
	if err != nil {
		return PageData{}, fmt.Errorf("could not execute JavaScript: %w", err)
	}

	if result == nil {
		return PageData{}, nil // 返回空的 PageData 以便上层处理
	}

	data := result.(map[string]interface{})
	title, titleOk := data["title"].(string)
	node, nodeOk := data["node"].(string)
	if !titleOk || !nodeOk {
		return PageData{}, nil // 返回空的 PageData 以便上层处理
	}

	return PageData{
		Title: title,
		Node:  node,
	}, nil
}

func postToTelegraph(pageData PageData) (string, error) {
	data := map[string]interface{}{
		"title":        pageData.Title,
		"author_name":  "AIOPrivacyBot",
		"content":      pageData.Node,
		"access_token": accessToken,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("error marshalling JSON: %w", err)
	}

	resp, err := http.Post(telegraphAPI, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	var response CreatePageResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("error decoding response: %w", err)
	}

	if response.OK {
		return response.Result.URL, nil
	} else {
		return "", fmt.Errorf("failed to create page")
	}
}
