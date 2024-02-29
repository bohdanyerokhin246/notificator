package config

type DBConfig struct {
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type Notification struct {
	ID                  int    `json:"ID"`
	Message             string `json:"message"`
	ActivationArguments string `json:"activationArguments"`
}

type NotificationList struct {
	ID                  int    `json:"ID"`
	Message             string `json:"message"`
	ActivationArguments string `json:"activationArguments"`
}

type ReadBy struct {
	ID               int      `json:"ID,omitempty"`
	NotificationName string   `json:"notificationName,omitempty"`
	ReadByDoctorList []string `json:"readByDoctorList,omitempty"`
	ReadByTime       []string `json:"readByTime,omitempty"`
}

var NotificationPtr *Notification = nil
var NotificationListPtr []*Notification = nil
var ReadByListPtr []*ReadBy = nil
