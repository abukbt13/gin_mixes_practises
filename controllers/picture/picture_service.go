package picture

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"practise/config"
	"practise/logs"
	"practise/models"
	"strings"
)

func SavePictures(c *gin.Context) {
	var req models.Picture
	if err := c.ShouldBind(&req); err != nil {
		logs.LogToFile("Error binding request: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if strings.TrimSpace(req.Description) == "" || strings.TrimSpace(req.PictureName) == "" {
		// Log the error to a file
		logs.LogToFile("All fields are required")
		// Return a response with status Bad Request (400)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": "All fields are required"})
		return
	}

	if err := config.DB.Create(&req).Error; err != nil {
		logs.LogToFile("Error saving picture: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save picture"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"picture": req,
	})
}
