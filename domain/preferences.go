package domain

type Lang string

const (
	English Lang = "en"
	Amharic Lang = "am"
)

type Preferences struct {
	ID                string
	UserID            string
	PreferredLang     Lang
	PushNotification  bool
	EmailNotification bool
}
