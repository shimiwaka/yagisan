package schema

import (
	"github.com/jinzhu/gorm"
)

type Box struct {
	gorm.Model `json:"-"`
	Username   string `json:"username"`
	Password   string `json:"password"`
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
	Mail       string `json:"mail"`
	IP         string `json:"ip"`
}
