package schema

import (
	"github.com/jinzhu/gorm"
)

type Box struct {
	gorm.Model  `json:"-"`
	Username    string `json:"username" gorm:"unique"`
	Password    string `json:"password"`
	Email       string `jsoin:"email" gorm:"unique"`
	Description string `json:"description"`
}

type Question struct {
	gorm.Model
	Box        uint   `json:"-"`
	Email      string `json:"-"`
	IP         string `json:"-"`
	UserAgent  string `json:"-"`
	Body       string `json:"body"`
	Token      string `json:"token"`
	Visible    bool   `json:"-"`
}

type Answer struct {
	gorm.Model `json:"-"`
	Question   uint   `json:"question"`
	Body       string `json:"body"`
}

type Block struct {
	gorm.Model `json:"-"`
	Box        uint `json:"box"`
	Type       int  `json:"type"`
	Value      int  `json:"value"`
}

type BlockMail struct {
	gorm.Model `json:"-"`
	Box        uint `json:"box"`
	Value      int  `json:"value"`
}

type AccessToken struct {
	gorm.Model `json:"-"`
	Box        uint   `json:"box"`
	Token      string `json:"token" gorm:"unique"`
}

type PasswordResetToken struct {
	gorm.Model `json:"-"`
	Box        uint   `json:"box"`
	Token      string `json:"token"`
}
