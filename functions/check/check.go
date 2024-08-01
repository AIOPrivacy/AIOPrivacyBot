package check

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
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
	dataURL       = "https://raw.githubusercontent.com/iuu6/AIOPrivacyBot/main/functions/check/data.minify.json"
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
	defer providersLock.Unlock()
	providers = newProviders.Providers

	log.Println("Successfully updated providers")
	return nil
}

func cleanURL(input string) (string, error) {
	providersLock.RLock()
	defer providersLock.RUnlock()

	for _, provider := range providers {
		if matched, _ := regexp.MatchString(provider.URLPattern, input); matched {
			//log.Printf("Matching provider found: %v", provider.URLPattern)
			for _, exception := range provider.Exceptions {
				re := regexp.MustCompile(fmt.Sprintf(`(?i)%s`, exception))
				if re.MatchString(input) {
					log.Printf("Exception matched, returning URL: %s", input)
					return input, nil
				}
			}
			for _, rawRule := range provider.RawRules {
				re := regexp.MustCompile(fmt.Sprintf(`(?i)%s`, rawRule))
				input = re.ReplaceAllString(input, "")
				//log.Printf("Applied raw rule: %s", rawRule)
			}
			parsed, err := url.Parse(input)
			if err != nil {
				return input, err
			}
			values := parsed.Query()
			for key := range parsed.Query() {
				for _, rule := range provider.Rules {
					re := regexp.MustCompile(fmt.Sprintf(`(?i)%s`, rule))
					if re.MatchString(key) {
						values.Del(key)
					}
					//log.Printf("Applied rule: %s", rule)
				}
				for _, refParam := range provider.ReferralMarketing {
					re := regexp.MustCompile(fmt.Sprintf(`(?i)%s`, refParam))
					if re.MatchString(key) {
						values.Del(key)
					}
					//log.Printf("Applied referral marketing rule: %s", refParam)
				}
			}
			parsed.RawQuery = values.Encode()
			input = parsed.String()
		}
	}
	return input, nil
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

	cleanedURL, err := cleanURL(url)

	if err != nil {
		return
	}

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
