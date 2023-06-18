package main

import (
	"github.com/jinzhu/gorm"
	"github.com/shimiwaka/yagisan/schema"
)

func initializeDB(db *gorm.DB) {
	db.Exec("DROP TABLE boxes")
	db.AutoMigrate(&schema.Box{})

	db.Exec("DROP TABLE questions")
	db.AutoMigrate(&schema.Question{})

	db.Exec("DROP TABLE answers")
	db.AutoMigrate(&schema.Answer{})

	db.Exec("DROP TABLE blocks")
	db.AutoMigrate(&schema.Block{})

	// db.Exec("DROP TABLE blockmails")
	db.AutoMigrate(&schema.BlockMail{})

	// db.Exec("DROP TABLE accesstokens")
	db.AutoMigrate(&schema.AccessToken{})

	// db.Exec("DROP TABLE passwordresettokens")
	db.AutoMigrate(&schema.PasswordResetToken{})
}
