package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/barkhayot/auth-fiber/models"  
)

// Global variable
var DB *gorm.DB

func Connect() {
	dsn := "host=localhost user=postgres password=010203 dbname=auth-fiber port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	connection, err  := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("couldn't connect to database ")
	}

	DB = connection

	connection.AutoMigrate(&models.User{})
}