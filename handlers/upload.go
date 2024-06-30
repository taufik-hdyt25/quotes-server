package handlers

import (
	"context"
	"net/http"
	"os"
	"strconv"

	"go-photo-upload/models"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type Handler struct {
	db *sqlx.DB
}

func NewHandler(db *sqlx.DB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) UploadPhoto(c *gin.Context) {
	// Retrieve the file from the form-data
	file, err := c.FormFile("photo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file"})
		return
	}

	// Open the file
	fileContent, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer fileContent.Close()

	// Initialize Cloudinary
	cld, err := cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize Cloudinary"})
		return
	}

	// Upload to Cloudinary in the "quotes" folder
	ctx := context.Background()
	resp, err := cld.Upload.Upload(ctx, fileContent, uploader.UploadParams{Folder: "quotes"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload to Cloudinary"})
		return
	}

	// Save the URL to the database
	photo := models.Photo{
		URL: resp.SecureURL,
	}

	_, err = h.db.NamedExec(`INSERT INTO photos (url) VALUES (:url)`, &photo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save to database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": photo.URL})
}

func (h *Handler) GetAllPhotos(c *gin.Context) {
	var photos []models.Photo

	// Default limit and page values
	limit := 10
	page := 1

	// Get limit and page from query parameters
	if l, err := strconv.Atoi(c.Query("limit")); err == nil {
		limit = l
	}
	if p, err := strconv.Atoi(c.Query("page")); err == nil {
		page = p
	}

	offset := (page - 1) * limit

	// Query photos from the database with limit and offset for pagination
	err := h.db.Select(&photos, "SELECT id, url FROM photos ORDER BY RANDOM() LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch photos"})
		return
	}

	// Check if there are more photos for the next page
	var nextPhotos []models.Photo
	err = h.db.Select(&nextPhotos, "SELECT id, url FROM photos ORDER BY RANDOM() LIMIT 1 OFFSET $1", offset+limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check for next page"})
		return
	}

	nextPage := 0
	if len(nextPhotos) > 0 {
		nextPage = page + 1
	}

	c.JSON(http.StatusOK, gin.H{
		"current_page": page,
		"next_page":    nextPage,
		"photos":       photos,
	})
}

func (h *Handler) DeletePhoto(c *gin.Context) {
	id := c.Param("id")

	// Execute delete query
	result, err := h.db.Exec("DELETE FROM photos WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete photo"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve affected rows"})
		return
	}

	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Photo not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Photo deleted successfully"})
}
