package config

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"practise/models"
)

var (
	DB *gorm.DB
)

func ConnectDatabase() {
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")

	dsn := dbUser + ":" + dbPass + "@tcp(127.0.0.1:3306)/" + dbName

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database!")
	}
	err = database.AutoMigrate(&models.User{}, &models.Picture{})
	if err != nil {
		panic("Failed to perform auto migration!")
	}

	DB = database
}
