package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

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
	tokenPtr := flag.String("token", "", "Telegram bot token")
	portPtr := flag.String("port", "8080", "Port to listen")
	flag.Parse()

	if *tokenPtr == "" {
		log.Println("Bot token must be specified.")
		return
	}

	setWebhook(*tokenPtr)

	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		webhookHandler(w, r, *tokenPtr)
	})
	log.Println("Starting server at :" + *portPtr)
	err := http.ListenAndServe(":"+*portPtr, nil)
	if err != nil {
		log.Printf("Error starting server: %v\n", err)
	}
}

func webhookHandler(w http.ResponseWriter, r *http.Request, token string) {
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading request:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	log.Println(string(bytes))

	var update Update
	if err := json.Unmarshal(bytes, &update); err != nil {
		log.Println("Error parsing update:", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	response := update.Message.Text
	sendMessage(update.Message.Chat.ID, response, token)
}

func sendMessage(chatID int, text string, token string) {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)

	values := url.Values{}
	values.Set("chat_id", fmt.Sprintf("%d", chatID))
	values.Set("text", text)

	_, err := http.PostForm(apiURL, values)
	if err != nil {
		log.Println("Error sending message:", err)
		return
	}
}

func setWebhook(token string) {
	webhookUrl := "https://dev.bayborodin.ru/webhook"

	url := fmt.Sprintf("https://api.telegram.org/bot%s/setWebhook?url=%s", token, webhookUrl)

	_, err := http.Get(url)
	if err != nil {
		log.Println("Error setting webhook:", err)
		return
	}

	log.Println("Webhook registered successfully.")
}
