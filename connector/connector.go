package connector

import (
	_ "embed"
	"fmt"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Settings struct {
	Username string `yaml:"db_username"`
	Pass     string `yaml:"db_pass"`
	Host     string `yaml:"db_host"`
	Port     int    `yaml:"db_port"`
	Name     string `yaml:"db_name"`
}

func ConnectDB() *gorm.DB {
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
		return nil
	}

	return db
}

func ConnectTestDB() *gorm.DB {
	settings := Settings{}
	b, _ := os.ReadFile("config_test.yaml")
	yaml.Unmarshal(b, &settings)

	connectQuery := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		settings.Username, settings.Pass, settings.Host, settings.Port, settings.Name)

	db, err := gorm.Open("mysql", connectQuery)

	if err != nil {
		fmt.Println(connectQuery)
		fmt.Println("error: failed to connect DB.")
		fmt.Println(err)
		return nil
	}

	return db
}