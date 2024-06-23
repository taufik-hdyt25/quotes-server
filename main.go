package main

import (
	"log"
	"time"

	"go-photo-upload/handlers"
	"go-photo-upload/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	db := utils.SetupDB()
	handler := handlers.NewHandler(db)

	router.POST("/upload", handler.UploadPhoto)
	router.GET("/photos", handler.GetAllPhotos)

	router.Run(":8080")
}
