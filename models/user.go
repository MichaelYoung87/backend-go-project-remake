package models

type User struct {
	UserID   int    `json:"user_id" gorm:"PRIMARY_KEY;AUTO_INCREMENT;NOT NULL"`
	Username string `json:"username" gorm:"unique"`
	Password []byte `json:"-"`
	ImageURL string `json:"image_url" gorm:"default:''"`
}