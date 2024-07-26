package database

import (
	"database1/model"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

type DbInstance struct {
	DB *gorm.DB
}

var Database DbInstance

func Connectdb() {

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", os.Getenv("HOST"), os.Getenv("USER"), os.Getenv("PASSWORD"), os.Getenv("DATABASE"), os.Getenv("DATABASEPORT"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(2)
	}
	db.AutoMigrate(&model.User{}, &model.Category{}, &model.Product{}, &model.Cart{})
	Database = DbInstance{DB: db}
}
