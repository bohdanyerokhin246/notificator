package main

import (
	"encoding/json"
	"fmt"
	t "github.com/go-toast/toast"
	"io"
	"net/http"
)

type Notifications struct {
	Message             string `json:"message"`
	ActivationArguments string `json:"activationArguments"`
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	// Ваша логика обработки запроса здесь

	// Чтение тела запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Ошибка чтения тела запроса", http.StatusInternalServerError)
		return
	}

	var notification Notifications
	err = json.Unmarshal(body, &notification)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	fmt.Println(notification.ActivationArguments)
	fmt.Println(notification.Message)

	toast := t.Notification{
		AppID:    "СПОВІЩЕННЯ КНП МП №6",
		Message:  notification.Message,
		Icon:     "C:/notificator/logo.png",
		Duration: t.Long,
		Actions: []t.Action{
			{"protocol", "Прочитати", notification.ActivationArguments},
		},
	}

	if err := toast.Push(); err != nil {
		fmt.Println("Ошибка при выводе уведомления:", err)
	}

	_, err = fmt.Fprintf(w, "Запрос успешно обработан!")
	if err != nil {
		return
	}
}

func main() {
	http.HandleFunc("/notification", handleRequest)
	fmt.Println("Client start on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}
