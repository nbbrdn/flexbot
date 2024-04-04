package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	baseURL = "https://api.telegram.org/bot"
)

var (
	token  string
	logger *log.Logger
)

type Config struct {
	Token string
	Port  string
}

type Update struct {
	UpdateID int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	MessageID int    `json:"message_id"`
	Chat      Chat   `json:"chat"`
	Text      string `json:"text"`
}

type Chat struct {
	ID int `json:"id"`
}

func main() {
	var config Config
	flag.StringVar(&config.Token, "token", "", "Telegram bot token")
	flag.StringVar(&config.Port, "port", "8080", "Port to listen")
	flag.Parse()

	if config.Token == "" {
		log.Println("Bot token must be specified.")
		return
	}

	logger = log.New(os.Stdout, "[flexbot] ", log.Ldate|log.Ltime|log.Lshortfile)

	setWebhook(config.Token)

	http.HandleFunc("/webhook", webhookHandler)
	logger.Println("Starting server at :" + config.Port)
	err := http.ListenAndServe(":"+config.Port, nil)
	if err != nil {
		logger.Fatalf("Error starting server: %v\n", err)
	}
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Println("Error reading request:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	logger.Println(string(bytes))

	var update Update
	if err := json.Unmarshal(bytes, &update); err != nil {
		logger.Println("Error parsing update:", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	processUpdate(update)
}

func processUpdate(update Update) {
	message := update.Message.Text
	chatID := update.Message.Chat.ID

	if strings.HasPrefix(message, "/start") {
		sendMessage(chatID, "Hi! I'm a bot. Can I help you?")
	} else if strings.HasPrefix(message, "/stop") {
		sendMessage(chatID, "Chao! By the way you can write me anytyme...")
	} else if strings.HasPrefix(message, "/") {
		sendMessage(chatID, "Unknown command.")
	} else {
		sendMessage(chatID, message)
	}
}

func sendMessage(chatID int, text string) {
	apiURL := fmt.Sprintf("%s%s/sendMessage", baseURL, token)

	values := url.Values{}
	values.Set("chat_id", fmt.Sprintf("%d", chatID))
	values.Set("text", text)

	_, err := http.PostForm(apiURL, values)
	if err != nil {
		logger.Println("Error sending message:", err)
		return
	}
}

func setWebhook(token string) {
	webhookUrl := "https://dev.bayborodin.ru/webhook"

	url := fmt.Sprintf("%s%s/setWebhook?url=%s", baseURL, token, webhookUrl)

	_, err := http.Get(url)
	if err != nil {
		logger.Println("Error setting webhook:", err)
		return
	}

	logger.Println("Webhook registered successfully.")
}
