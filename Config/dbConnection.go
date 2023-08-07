package Config

import (
	"fmt"
	models "go-social-media-api/Models"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	godotenv.Load()
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s port=5432 user=%s password=%s dbname=%s sslmode=disable", dbHost, dbUser, dbPass, dbName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	DB = db
	fmt.Println("Database connected")

	Migration(DB)
}

func Migration(conn *gorm.DB) {
	conn.Debug().AutoMigrate(
		&models.Users{},
		&models.Posts{},
		&models.PostLikes{},
		&models.PostCommets{},
	)
}
