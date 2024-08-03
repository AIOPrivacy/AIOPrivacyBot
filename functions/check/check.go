package check

import (
	"bytes"
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
	dataURL       = "https://raw.githubusercontent.com/AIOPrivacy/AIOPrivacyBot/main/functions/check/data.minify.json"
	providers     map[string]Provider
	providersLock sync.RWMutex
	apiKey        string
)

func Init(apiKeyFromConfig string) {
	apiKey = apiKeyFromConfig
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

type ThreatEntry struct {
	URL string `json:"url"`
}

type ThreatInfo struct {
	ThreatTypes      []string      `json:"threatTypes"`
	PlatformTypes    []string      `json:"platformTypes"`
	ThreatEntryTypes []string      `json:"threatEntryTypes"`
	ThreatEntries    []ThreatEntry `json:"threatEntries"`
}

type SafeBrowsingRequest struct {
	Client     ClientInfo `json:"client"`
	ThreatInfo ThreatInfo `json:"threatInfo"`
}

type ClientInfo struct {
	ClientID      string `json:"clientId"`
	ClientVersion string `json:"clientVersion"`
}

type SafeBrowsingResponse struct {
	Matches []struct {
		ThreatType      string `json:"threatType"`
		PlatformType    string `json:"platformType"`
		ThreatEntryType string `json:"threatEntryType"`
		Threat          struct {
			URL string `json:"url"`
		} `json:"threat"`
		CacheDuration string `json:"cacheDuration"`
	} `json:"matches"`
}

func checkURLSafety(input string) (*SafeBrowsingResponse, error) {
	safeBrowsingURL := fmt.Sprintf("https://safebrowsing.googleapis.com/v4/threatMatches:find?key=%s", apiKey)

	clientInfo := ClientInfo{
		ClientID:      "yourcompanyname",
		ClientVersion: "1.5.2",
	}

	threatInfo := ThreatInfo{
		ThreatTypes:      []string{"MALWARE", "SOCIAL_ENGINEERING", "UNWANTED_SOFTWARE", "POTENTIALLY_HARMFUL_APPLICATION"},
		PlatformTypes:    []string{"WINDOWS", "LINUX", "ANDROID", "IOS", "OSX", "CHROME"},
		ThreatEntryTypes: []string{"URL", "EXECUTABLE"},
		ThreatEntries: []ThreatEntry{
			{URL: input},
		},
	}

	requestBody := SafeBrowsingRequest{
		Client:     clientInfo,
		ThreatInfo: threatInfo,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(safeBrowsingURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, response: %s", resp.StatusCode, body)
	}

	var safeBrowsingResponse SafeBrowsingResponse
	err = json.Unmarshal(body, &safeBrowsingResponse)
	if err != nil {
		return nil, err
	}

	return &safeBrowsingResponse, nil
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

	// Check URL safety using Google Safe Browsing API
	safetyResponse, err := checkURLSafety(cleanedURL)
	if err != nil {
		log.Printf("Error checking URL safety: %v", err)
		return
	}

	safetyMessage := "The URL is safe."
	if len(safetyResponse.Matches) > 0 {
		threatDetails := ""
		for _, match := range safetyResponse.Matches {
			threatDetails += fmt.Sprintf("Threat Type: %s\nPlatform Type: %s\nThreat Entry Type: %s\nThreat URL: %s\nCache Duration: %s\n---\n",
				match.ThreatType, match.PlatformType, match.ThreatEntryType, match.Threat.URL, match.CacheDuration)
		}
		safetyMessage = fmt.Sprintf("The URL is not safe. Threat details:\n%s", threatDetails)
	}

	results := []interface{}{
		tgbotapi.NewInlineQueryResultArticleHTML(
			inlineQuery.ID,
			"URL Clean & Safe Check",
			fmt.Sprintf("网址（去跟踪）：<a href=\"%s\">%s</a>\n\n\n%s", cleanedURL, cleanedURL, safetyMessage),
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
