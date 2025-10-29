package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Gagal terkoneksi ke database ", err)
	}
	DB = database
	fmt.Println("Terkoneksi dengan database ")

	//Tes koneksi
	sqlDB, err := database.DB()
	if err != nil {
		log.Fatal("gagal untuk database instance : ", err)
	}
	if err := sqlDB.Ping(); err != nil {
		log.Fatal("gagal ping ke database : ", err)
	}

	fmt.Println("Database Ping success!")
}
