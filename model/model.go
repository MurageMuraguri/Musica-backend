package model

import (
	"time"

	_ "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type User struct {
	Common
	UserName      string `json:"user_name"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	ApiKey        string `json:"api_key"`
	EmailVerified bool   `json:"email_verified"` // audio duration
	Enabled       bool   `json:"enabled"`        // whether as user is activated....defaults
	PhoneNumber   string `json:"phone_number"`
	Provider      string `json:"provider"` // email, thirdparty like google
	ProviderData  string `json:"provider_data"`
	PhotoUrl      string `json:"photo_url"`
}
type Login struct {
	ID        string `gorm:"primaryKey"`
	CreatedAt time.Time
	UserId    string //uuid of user who logged in
	IpAddress string
	UserAgent string
}
type Common struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Audio struct {
	IsMatch        bool   `json:"is_match"`
	MatchedAudioId string `json:"matched_audio_id"`

}
type File struct {
	Common

	ImgURL       string `json:"img_url"` //link to file
	ImgURL300    string `json:"img_url_300"`
	WavURL       string `json:"wav_url"`
	Mp3URL       string `json:"mp3_url"`
	ZipURL       string `json:"zip_url"`
	Size         int64  `json:"size"`
	OriginalName string `json:"original_name"`
	MimeType     string `json:"mime_type"`
	UserId       string `json:"user_id"`
}
