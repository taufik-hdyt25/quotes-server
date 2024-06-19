package main

import (
	"log"

	"go-photo-upload/handlers"
	"go-photo-upload/utils"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	router := gin.Default()

	db := utils.SetupDB()
	handler := handlers.NewHandler(db)

	router.POST("/upload", handler.UploadPhoto)
	router.GET("/photos", handler.GetAllPhotos)

	router.Run(":8080")
}
