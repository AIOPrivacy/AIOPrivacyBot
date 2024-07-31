package check

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Provider struct {
	URLPattern        string   `json:"urlPattern"`
	Rules             []string `json:"rules"`
	RawRules          []string `json:"rawRules"`
	ReferralMarketing []string `json:"referralMarketing"`
	Exceptions        []string `json:"exceptions"`
	CompleteProvider  bool     `json:"completeProvider,omitempty"`
	ForceRedirection  bool     `json:"forceRedirection,omitempty"`
}

type ClearURLsData struct {
	Providers map[string]Provider `json:"providers"`
}

var (
	dataURL       = "https://rules2.clearurls.xyz/data.minify.json"
	providers     map[string]Provider
	providersLock sync.RWMutex
)

func init() {
	go refreshData()
}

func refreshData() {
	for {
		err := fetchData()
		if err != nil {
			log.Printf("Error fetching data: %v", err)
		}
		time.Sleep(30 * time.Minute)
	}
}

func fetchData() error {
	log.Println("Fetching data from URL:", dataURL)
	resp, err := http.Get(dataURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch data: status code %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var newProviders ClearURLsData
	err = json.Unmarshal(data, &newProviders)
	if err != nil {
		return fmt.Errorf("failed to load providers: %v", err)
	}

	providersLock.Lock()
	providers = newProviders.Providers
	providersLock.Unlock()

	log.Println("Successfully updated providers")
	return nil
}

func cleanURL(url string) string {
	providersLock.RLock()
	defer providersLock.RUnlock()

	for _, provider := range providers {
		if matched, _ := regexp.MatchString(provider.URLPattern, url); matched {
			log.Printf("Matching provider found: %v", provider.URLPattern)
			for _, rule := range provider.Rules {
				re := regexp.MustCompile(fmt.Sprintf(`[?&]%s=[^&]*`, rule))
				url = re.ReplaceAllString(url, "")
				log.Printf("Applied rule: %s", rule)
			}
			for _, rawRule := range provider.RawRules {
				re := regexp.MustCompile(rawRule)
				url = re.ReplaceAllString(url, "")
				log.Printf("Applied raw rule: %s", rawRule)
			}
			for _, refParam := range provider.ReferralMarketing {
				re := regexp.MustCompile(fmt.Sprintf(`[?&]%s=[^&]*`, refParam))
				url = re.ReplaceAllString(url, "")
				log.Printf("Applied referral marketing rule: %s", refParam)
			}
			for _, exception := range provider.Exceptions {
				re := regexp.MustCompile(exception)
				if re.MatchString(url) {
					log.Printf("Exception matched, returning URL: %s", url)
					return url
				}
			}
			url = cleanQueryString(url)
		}
	}
	return url
}

func cleanQueryString(url string) string {
	parts := strings.SplitN(url, "?", 2)
	if len(parts) < 2 {
		return url
	}
	baseURL := parts[0]
	query := parts[1]

	params := strings.Split(query, "&")
	paramMap := make(map[string]string)
	for _, param := range params {
		keyValue := strings.SplitN(param, "=", 2)
		if len(keyValue) == 2 {
			key := keyValue[0]
			value := keyValue[1]
			if _, exists := paramMap[key]; !exists {
				paramMap[key] = value
			}
		}
	}

	var cleanedQuery []string
	for key, value := range paramMap {
		cleanedQuery = append(cleanedQuery, fmt.Sprintf("%s=%s", key, value))
	}
	return baseURL + "?" + strings.Join(cleanedQuery, "&")
}

func HandleInlineQuery(inlineQuery *tgbotapi.InlineQuery, bot *tgbotapi.BotAPI) {
	log.Printf("Received inline query from %s: %s", inlineQuery.From.UserName, inlineQuery.Query)

	query := strings.TrimSpace(inlineQuery.Query)
	if !strings.HasPrefix(query, "-check") {
		return
	}

	url := strings.TrimPrefix(query, "-check")
	url = strings.TrimSpace(url)
	if url == "" {
		return
	}

	cleanedURL := cleanURL(url)

	results := []interface{}{
		tgbotapi.NewInlineQueryResultArticleHTML(
			inlineQuery.ID,
			"Cleaned URL",
			fmt.Sprintf("网址（去跟踪）：<a href=\"%s\">%s</a>", cleanedURL, cleanedURL),
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
