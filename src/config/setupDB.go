package config

import (
	"fmt"
	"os"
	"posts-ms/src/entity"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupDB() (*gorm.DB, error) {
	host := os.Getenv("DATABASE_DOMAIN")
	user := os.Getenv("DATABASE_USERNAME")
	password := os.Getenv("DATABASE_PASSWORD")
	name := os.Getenv("DATABASE_SCHEMA")
	port := os.Getenv("DATABASE_PORT")

	connectionString := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host,
		user,
		password,
		name,
		port,
	)

	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})

	db.AutoMigrate(&entity.Post{Tbl: "posts"})
	db.AutoMigrate(&entity.Like{Tbl: "likes"})
	db.AutoMigrate(&entity.Comment{Tbl: "comments"})

	return db, err
}
