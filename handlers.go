package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"notificator/config"
	"notificator/db"
	"os"
	"time"
)

func handleMainPage(w http.ResponseWriter, _ *http.Request) {

	var tmpl = template.Must(template.ParseFiles("templates/index.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		fmt.Printf("Problem with execution templates/index.html. Error: %s", err)
	}

}

func handleAdminPage(w http.ResponseWriter, _ *http.Request) {

	var tmpl = template.Must(template.ParseFiles("templates/adminPage.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		fmt.Printf("Problem with execution templates/adminPage.html. Error: %s", err)
	}
}

func handleAdminMessage(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		ipGroup := r.FormValue("ipGroup")
		message := r.FormValue("message")

		id, err := db.InsertNotification(message)
		if err != nil {
			fmt.Printf("Problem with InsertNotification . Error: %s", err)
		}

		err = db.InsertReadByList(id)
		if err != nil {
			fmt.Printf("Problem with InsertReadByList . Error: %s", err)
		}

		config.NotificationPtr = new(config.Notification)
		config.NotificationPtr.ID = id
		config.NotificationPtr.Message = message

		ipAddresses, _ := readStringArrayFromFile("config/" + ipGroup + ".txt")
		var request [255]*http.Request

		for i := 0; i < len(ipAddresses); i++ {

			config.NotificationPtr.ActivationArguments = "http://10.20.77.37:8081/answer"

			// Кодируем структуру User в JSON
			requestBody, err := json.Marshal(config.NotificationPtr)
			if err != nil {
				fmt.Printf("Problem with marshaling adminMessage structur to JSON . Error: %s", err)

			}

			request[i], err = http.NewRequest("POST", fmt.Sprintf("http://%s:8080/notification", ipAddresses[i]), bytes.NewBuffer(requestBody))
			if err != nil {
				fmt.Printf("Problem with creating new request. Error: %s", err)

			}
		}

		for i := 0; i < len(request)-1; i++ {
			go sendNotification(request[i])
			fmt.Println(request[i].URL, "Done")
		}
	}
	handleNotificationSent(w, r)
}

func sendNotification(request *http.Request) {
	client := &http.Client{}
	_, err := client.Do(request)
	if err != nil {
		//fmt.Printf("Problem with sending message to client. Error: %s", err)
	}
}

func handleNotificationSent(w http.ResponseWriter, _ *http.Request) {

	var tmpl = template.Must(template.ParseFiles("templates/notificationSent.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		fmt.Printf("Problem with execution templates/notificationSent.html. Error: %s", err)
	}
}

func handleShowNotification(w http.ResponseWriter, _ *http.Request) {

	err := db.SelectNotifications()
	var tmpl = template.Must(template.ParseFiles("templates/showNotifications.html"))
	err = tmpl.Execute(w, &config.NotificationListPtr)
	if err != nil {
		fmt.Printf("Problem with execution templates/answerToClient.html. Error: %s", err)
	}

}

func handleShowReadByList(w http.ResponseWriter, r *http.Request) {

	err := db.SelectAllReadByList(r.URL.RawQuery)
	var tmpl = template.Must(template.ParseFiles("templates/showReadByList.html"))
	err = tmpl.Execute(w, &config.ReadByListPtr)
	if err != nil {
		fmt.Printf("Problem with execution templates/answerToClient.html. Error: %s", err)
	}

}

func handleShowReadByName(w http.ResponseWriter, r *http.Request) {

	err := db.SelectAllReadByName(r.URL.RawQuery)
	var tmpl = template.Must(template.ParseFiles("templates/showReadByName.html"))
	err = tmpl.Execute(w, &config.ReadByListPtr)
	if err != nil {
		fmt.Printf("Problem with execution templates/answerToClient.html. Error: %s", err)
	}

}

func handleAnswerFromClient(w http.ResponseWriter, _ *http.Request) {

	var tmpl = template.Must(template.ParseFiles("templates/answerToClient.html"))
	err := tmpl.Execute(w, &config.NotificationPtr)
	if err != nil {
		fmt.Printf("Problem with execution templates/answerToClient.html. Error: %s", err)
	}

}

func handleSaveInfo(w http.ResponseWriter, r *http.Request) {

	var doctorSurname string
	var answer = [2]string{"Дякуємо за те, що ознайомились зі сповіщенням, можете закрити цю вкладку", "Сталася помилка звернітся до адміністратора"}

	var tmpl = template.Must(template.ParseFiles("templates/answerAccess.html"))

	if r.Method == http.MethodPost {
		doctorSurname = r.FormValue("surname")
		ipAddress := r.RemoteAddr
		notificationID := r.URL.RawQuery

		err := db.UpdateReadBy(notificationID, doctorSurname, ipAddress, time.Now())
		if err != nil {
			fmt.Printf("UpdateReadBy error. Error: %s", err)
		}

		err = tmpl.Execute(w, answer[0])
		if err != nil {
			fmt.Printf("Problem with execution templates/answerAccess.html with answer[0]. Error: %s", err)
		}
	} else {

		err := tmpl.Execute(w, answer[1])
		if err != nil {
			fmt.Printf("Problem with execution templates/answerAccess.html with answer[1]. Error: %s", err)
		}
	}
}

func readStringArrayFromFile(filename string) ([]string, error) {

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {

		}
	}(file)

	scanner := bufio.NewScanner(file)
	var strings []string
	for scanner.Scan() {
		strings = append(strings, scanner.Text())
	}

	if err = scanner.Err(); err != nil {
		return nil, err
	}

	fmt.Printf("IP addresses was read %s\n", filename)
	return strings, nil
}
