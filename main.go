package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
)

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

	response := "Hi! I am a simple Go bot."
	w.Write([]byte(response))
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
