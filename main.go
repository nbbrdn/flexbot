package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	MessageId int    `json:"message_id"`
	Chat      Chat   `json:"chat"`
	Text      string `json:"text"`
}

type Chat struct {
	Id int `json:"id"`
}

func main() {
	tokenPtr := flag.String("token", "", "Telegram bot token")
	portPtr := flag.String("port", "8080", "Port to listen")
	flag.Parse()

	if *tokenPtr == "" {
		fmt.Println("Bot token must be specified.")
		return
	}

	setWebhook(*tokenPtr)

	http.HandleFunc("/webhook", webhookHandler)
	fmt.Println("Starting server at :" + *portPtr)
	http.ListenAndServe(":"+*portPtr, nil)
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error reading request:", err)
		return
	}
	defer r.Body.Close()

	fmt.Println(string(bytes))

	w.Write([]byte("ok"))

	var update Update
	err = json.Unmarshal(bytes, &update)
	if err != nil {
		fmt.Println("Error parsing update:", err)
		return
	}

	response := "Hello! I'm a simple Telegram bot written in Go."
	sendMessage(update.Message.Chat.Id, response)
}

func sendMessage(chatId int, text string) {
	token := flag.Lookup("token").Value.(flag.Getter).Get().(string)
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)

	values := url.Values{}
	values.Set("chat_id", fmt.Sprintf("%d", chatId))
	values.Set("text", text)

	_, err := http.PostForm(apiURL, values)
	if err != nil {
		fmt.Println("Error sending message:", err)
		return
	}
}

func setWebhook(token string) {
	webhookUrl := "https://dev.bayborodin.ru/webhook"

	url := fmt.Sprintf("https://api.telegram.org/bot%s/setWebhook?url=%s", token, webhookUrl)

	_, err := http.Get(url)
	if err != nil {
		fmt.Println("Error setting webhook:", err)
		return
	}

	fmt.Println("Webhook registered successfully.")
}
