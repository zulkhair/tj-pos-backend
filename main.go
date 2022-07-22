package main

import (
	"dromatech/pos-backend/app"
	"dromatech/pos-backend/config"
	"dromatech/pos-backend/global"
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/sirupsen/logrus"
)

func main() {
	// load configuration
	fmt.Println("Load Configuration")
	cfg, err := config.New("config.yaml")
	if err != nil {
		fmt.Println(err.Error())
	}
	global.CONFIG = cfg

	// set up log
	f, err := os.OpenFile(cfg.LogConfig.LogFilename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0755)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.SetOutput(f)
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// set up database
	dsn := fmt.Sprintf("host=%s user=%s dbname=%s port=%s sslmode=disable", cfg.Database.Host, cfg.Database.User, cfg.Database.DBName, cfg.Database.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Fatal(err)
	}
	global.DBCON = db

	err = app.StartApp()
	if err != nil {
		logrus.Error(err.Error())
	}
}
