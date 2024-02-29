package db

import (
	"fmt"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"notificator/config"
	"time"
)

func InsertNotification(message string) (int, error) {
	// Пример запроса на вставку элемента в массив в таблице

	var id int
	err := DB.QueryRow("INSERT INTO notification(message) VALUES ($1) RETURNING id", message).Scan(&id)

	if err != nil {
		fmt.Printf("Problem with SQL INSERT notification . Error: %s", err)
	}

	fmt.Println("Элемент успешно вставлен в массив.")
	return id, err
}

func SelectNotifications() error {

	var size int

	err := DB.QueryRow("SELECT COUNT(*) FROM notification").Scan(&size)
	if err != nil {
		fmt.Printf("Problem with SQL INSERT readBy . Error: %s", err)
	}
	config.NotificationListPtr = make([]*config.Notification, size)

	rows, err := DB.Query("SELECT * from notification")
	if err != nil {
		fmt.Printf("Problem with SQL INSERT readBy . Error: %s", err)
	}

	for i := 0; rows.Next(); i++ {
		config.NotificationListPtr[i] = new(config.Notification)
		if err = rows.Scan(&config.NotificationListPtr[i].ID, &config.NotificationListPtr[i].Message); err != nil {
			return err
		}
	}

	fmt.Println("SELECT notification is access")
	return err
}

func InsertReadByList(notificationID int) error {
	// Пример запроса на вставку элемента в массив в таблице
	err := DB.QueryRow("INSERT INTO read_by(notification_id) VALUES ($1)", notificationID)
	if err.Err() != nil {
		fmt.Printf("Problem with SQL INSERT readBy . Error: %s", err.Err())
	}

	fmt.Println("Элемент успешно вставлен в массив.")
	return err.Err()
}

func UpdateReadBy(notificationID, doctorSurname, ipAddress string, time time.Time) error {

	err := DB.QueryRow("UPDATE read_by SET read_by_doctors_list = array_append(read_by_doctors_list, $1),read_on_computer_list = array_append(read_on_computer_list, $2),  time = array_append(time,$3) WHERE notification_ID = $4", doctorSurname, ipAddress, time, notificationID)

	if err.Err() != nil {
		fmt.Printf("Problem with SQL UPDATE readByList . Error: %s", err.Err())
	}

	fmt.Println("Элемент успешно вставлен в массив.")
	return err.Err()
}

func SelectAllReadByList(notificationID string) error {
	var size int

	err := DB.QueryRow("SELECT COUNT(*) FROM notification WHERE notification.id = $1", notificationID).Scan(&size)
	if err != nil {
		fmt.Printf("Problem with SQL INSERT readBy . Error: %s", err)
	}
	config.ReadByListPtr = make([]*config.ReadBy, size)

	rows, err := DB.Query("SELECT notification.message,read_by.read_by_doctors_list,read_by.time FROM notification, read_by WHERE notification.id = $1 AND read_by.notification_id=notification.id", notificationID)
	if err != nil {
		fmt.Printf("Problem with SQL INSERT readBy . Error: %s", err)
	}

	for i := 0; rows.Next(); i++ {

		config.ReadByListPtr[i] = new(config.ReadBy)
		if err = rows.Scan(&config.ReadByListPtr[i].NotificationName, pq.Array(&config.ReadByListPtr[i].ReadByDoctorList), pq.Array(&config.ReadByListPtr[i].ReadByTime)); err != nil {
			return err
		}
	}

	fmt.Println("SELECT read_by list is access")
	return err
}

func SelectAllReadByName(notificationID string) error {
	var size int

	err := DB.QueryRow("SELECT COUNT(*) FROM notification WHERE notification.id = $1", notificationID).Scan(&size)
	if err != nil {
		fmt.Printf("Problem with SQL INSERT readBy . Error: %s", err)
	}
	config.ReadByListPtr = make([]*config.ReadBy, size)

	rows, err := DB.Query("SELECT read_by.read_by_doctors_list FROM notification, read_by WHERE notification.id = $1 AND read_by.notification_id=notification.id", notificationID)
	if err != nil {
		fmt.Printf("Problem with SQL INSERT readBy . Error: %s", err)
	}

	for i := 0; rows.Next(); i++ {

		config.ReadByListPtr[i] = new(config.ReadBy)
		if err = rows.Scan(pq.Array(&config.ReadByListPtr[i].ReadByDoctorList)); err != nil {
			return err
		}
	}

	fmt.Println("SELECT read_by list is access")
	return err
}
