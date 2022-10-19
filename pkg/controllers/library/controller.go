package library

import (
	"net/http"

	"github.com/haapjari/glass/pkg/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GET /api/glass/libraries
func GetLibraries(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var libraries []models.Library
	db.Find(&libraries)

	c.JSON(http.StatusOK, gin.H{"data": libraries})
}

// POST /api/glass/libraries
func CreateLibrary(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// Validate Input
	var input models.CreateLibraryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create Model
	library := models.Library{Name: input.Name, Uri: input.Uri, Liblin: input.Liblin}
	db.Create(&library)
	c.JSON(http.StatusOK, gin.H{"data": library})
}

// GET /api/glass/libraries/:id
func GetLibraryById(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// Attempt to Fetch Models
	var library models.Library
	if err := db.Where("id = ?", c.Param("id")).First(&library).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": library})
}

func DeleteLibraryById(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// Attempt to Fetch Models
	var library models.Library
	if err := db.Where("id = ?", c.Param("id")).First(&library).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	db.Delete(&library)
	c.JSON(http.StatusOK, gin.H{"data": true})
}

func UpdateLibraryById(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// Attempt to Fetch Models
	var library models.Library
	if err := db.Where("id = ?", c.Param("id")).First(&library).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	// Validate Input
	var input models.UpdateLibraryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Model(&library).Updates(input)
	c.JSON(http.StatusOK, gin.H{"data": library})
}
