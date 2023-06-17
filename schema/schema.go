package schema

import (
	"github.com/jinzhu/gorm"
)

type Box struct {
	gorm.Model `json:"-"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Description   string `json:"description"`
}

type Question struct {
	gorm.Model `json:"-"`
	Box        uint   `json:"box"`
	Mail       string `json:"mail"`
	IP         string `json:"ip"`
	UserAgent  string `json:"user_agent"`
	Body       string `json:"body"`
	Token      string `json:"token"`
}

type Answer struct {
	gorm.Model `json:"-"`
	Question   uint   `json:"question"`
	Body       string `json:"body"`
}

type Block struct {
	gorm.Model `json:"-"`
	Box        uint   `json:"box"`
	Type       int `json:"type"`
	Value      int `json:"value"`
}

type BlockMail struct {
	gorm.Model `json:"-"`
	Box        uint   `json:"box"`
	Value      int `json:"value"`
}

type AccessToken struct {
	gorm.Model `json:"-"`
	Box        uint   `json:"box"`
	Token	string `json:"token"`
}

type PasswordResetToken struct {
	gorm.Model `json:"-"`
	Box        uint   `json:"box"`
	Token	string `json:"token"`
}