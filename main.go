package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"notificator/db"
)

func main() {

	db.InitDB()
	defer func(DB *sql.DB) {
		err := DB.Close()
		if err != nil {

		}
	}(db.DB)
	http.HandleFunc("/", handleMainPage)
	http.HandleFunc("/adminPage", handleAdminPage)
	http.HandleFunc("/notificationSent", handleAdminMessage)
	http.HandleFunc("/selectNotifications", handleShowNotification)
	http.HandleFunc("/showReadByList", handleShowReadByList)
	http.HandleFunc("/showReadByName", handleShowReadByName)
	http.HandleFunc("/answer", handleAnswerFromClient)
	http.HandleFunc("/saveInfo", handleSaveInfo)

	fmt.Println("Server start on :8081")
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		return
	}
}
