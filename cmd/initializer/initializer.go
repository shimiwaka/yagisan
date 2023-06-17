package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
	"github.com/jinzhu/gorm"
	"github.com/shimiwaka/yagisan/schema"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Settings struct {
	Username string `yaml:"db_username"`
	Pass string `yaml:"db_pass"`
	Host string `yaml:"db_host"`
	Port int `yaml:"db_port"`
	Name string `yaml:"db_name"`
}

func main() {
	settings := Settings{}
	b, _ := os.ReadFile("config.yaml")
	yaml.Unmarshal(b, &settings)

	connectQuery := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		settings.Username, settings.Pass, settings.Host, settings.Port, settings.Name)

	db, err := gorm.Open("mysql", connectQuery)

	if err != nil {
		fmt.Println(connectQuery)
		fmt.Println("error: failed to connect DB.")
		fmt.Println(err)
		return
	}

	db.Exec("DROP TABLE boxes")
	db.AutoMigrate(&schema.Box{})

	// db.Exec("DROP TABLE questions")
	db.AutoMigrate(&schema.Question{})

	// db.Exec("DROP TABLE answers")
	db.AutoMigrate(&schema.Answer{})

	// db.Exec("DROP TABLE blocks")
	db.AutoMigrate(&schema.Block{})
	
	// db.Exec("DROP TABLE blockmails")
	db.AutoMigrate(&schema.BlockMail{})

	fmt.Println("Successfully initialized.")
}