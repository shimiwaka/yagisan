package main

import (
	"time"
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/shimiwaka/yagisan/schema"
)

func validateAccessToken(db *gorm.DB, token string) schema.Box {
	box := schema.Box{}
	accessToken := schema.AccessToken{}

	err := db.First(&accessToken, "token = ? and created_at >= ?", token, time.Now().AddDate(0, 0, -7)).Error

	if err != nil || accessToken.Box == 0 {
		return schema.Box{}
	}

	err = db.First(&box, "ID = ?", accessToken.Box).Error

	if err != nil {
		return schema.Box{}
	}

	return box
}
