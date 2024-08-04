package ask

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"AIOPrivacyBot/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// AIRequest represents the structure of the request to the AI service.
type AIRequest struct {
	Model    string `json:"model"`
	Messages []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
	Stream bool `json:"stream"`
}

// AIResponseChunk represents a single chunk of the AI response.
type AIResponseChunk struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
		Index        int         `json:"index"`
		FinishReason interface{} `json:"finish_reason"`
	} `json:"choices"`
}

// HandleAsk processes messages and interacts with the AI service for academic questions.
func HandleAskCommand(message *tgbotapi.Message, bot *tgbotapi.BotAPI) {
	if message.Chat.IsPrivate() || message.CommandArguments() != "" {
		userMessage := generateUserMessage(message.CommandArguments())
		aiRequest := createAIRequest(userMessage)

		requestBody, err := json.Marshal(aiRequest)
		if err != nil {
			log.Fatalf("Error marshaling AI request: %v", err)
		}

		sendTypingAction(message.Chat.ID, bot)

		resp, err := http.Post("http://127.0.0.1:8080/v1/chat/completions", "application/json", bytes.NewBuffer(requestBody))
		if err != nil {
			log.Fatalf("Error making AI request: %v", err)
		}
		defer resp.Body.Close()

		combinedMessage, err := parseAIResponse(resp.Body)
		if err != nil {
			log.Fatalf("Error parsing AI response: %v", err)
		}

		err = utils.SendMarkdownMessage(message.Chat.ID, message.MessageID, combinedMessage, bot)
		if err != nil {
			log.Printf("Error sending Markdown message: %v", err)
			// Use HTML format as a fallback
			htmlMessage := convertMarkdownToHTML(combinedMessage)
			err = utils.SendMessage(message.Chat.ID, htmlMessage, message.MessageID, bot)
			if err != nil {
				log.Printf("Error sending AI response: %v", err)
			}
		}
	} else {
		utils.SendMarkdownMessage(message.Chat.ID, message.MessageID, "请发送 /ask@AIOPrivacyBot 你要说的内容", bot)
	}
}

// generateUserMessage creates the full message to be sent to the AI service.
func generateUserMessage(userInput string) string {
	systemPrompt := `我是一个能力超强的AI，我可以回答用户的一切问题，我非常注重用户隐私，我叫做AIOPrivacyBot，这是用户发来的消息，请你回答：`
	return systemPrompt + "\n" + userInput
}

// createAIRequest initializes the AI request payload.
func createAIRequest(content string) AIRequest {
	return AIRequest{
		Model: "gpt-4o-mini",
		Messages: []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}{
			{
				Role:    "user",
				Content: content,
			},
		},
		Stream: true,
	}
}

// sendTypingAction sends the "typing" action to indicate the bot is processing.
func sendTypingAction(chatID int64, bot *tgbotapi.BotAPI) {
	chatAction := tgbotapi.NewChatAction(chatID, tgbotapi.ChatTyping)
	if _, err := bot.Request(chatAction); err != nil {
		log.Printf("Error sending chat action: %v", err)
	}
}

// parseAIResponse reads and combines chunks from the AI response.
func parseAIResponse(responseBody io.Reader) (string, error) {
	body, err := ioutil.ReadAll(responseBody)
	if err != nil {
		return "", err
	}

	chunks := strings.Split(string(body), "\n")
	var combinedMessage string
	for _, chunk := range chunks {
		if !strings.HasPrefix(chunk, "data: ") {
			continue
		}
		var aiResponseChunk AIResponseChunk
		if err := json.Unmarshal([]byte(chunk[6:]), &aiResponseChunk); err != nil {
			log.Printf("Error unmarshaling AI response chunk: %v", err)
			continue
		}
		combinedMessage += aiResponseChunk.Choices[0].Delta.Content
	}
	return combinedMessage, nil
}

// convertMarkdownToHTML converts Markdown text to HTML text.
func convertMarkdownToHTML(markdownText string) string {
	// This is a basic conversion. You may need to handle more cases or use a library for a comprehensive conversion.
	htmlText := strings.ReplaceAll(markdownText, "*", "<b>")
	htmlText = strings.ReplaceAll(htmlText, "_", "<i>")
	htmlText = strings.ReplaceAll(htmlText, "`", "<code>")
	return htmlText
}
