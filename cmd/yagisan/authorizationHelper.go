package main

import (
	"github.com/jinzhu/gorm"
	"github.com/shimiwaka/yagisan/schema"
)

func validateAccessToken(db *gorm.DB, token string) schema.Box {
	box := schema.Box{}
	accessToken := schema.AccessToken{}

	err := db.First(&accessToken, "token = ?", token).Error

	if err != nil || accessToken.Box == 0 {
		return schema.Box{}
	}

	err = db.First(&box, "ID = ?", accessToken.Box).Error

	if err != nil {
		return schema.Box{}
	}
	return box
}
