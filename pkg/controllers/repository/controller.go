package repository

import (
	"net/http"

	"github.com/haapjari/glass/pkg/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RepositoryController struct {
	Handler        *Handler
	GitHubApiToken string
	GitHubName     string
}

func NewRepositoryController() *RepositoryController {
	return new(RepositoryController)
}

func GetRepositories(c *gin.Context) {
	var (
		h = NewHandler(c)
		e = h.HandleGetRepositories()
	)

	c.JSON(http.StatusOK, gin.H{"data": e})
}

func GetRepositoryById(c *gin.Context) {
	var (
		h = NewHandler(c)
		e = h.HandleGetRepositoryById()
	)

	c.JSON(http.StatusOK, gin.H{"data": e})
}

func CreateRepository(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// Validate Input
	var input models.CreateEntityInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create Model
	dataEntity := models.Entity{Name: input.Name, Uri: input.Uri, Lin: input.Lin, Liblin: input.Liblin, Issct: input.Issct}
	db.Create(&dataEntity)
	c.JSON(http.StatusOK, gin.H{"data": dataEntity})
}

func DeleteRepositoryById(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// Attempt to Fetch Models
	var dataEntity models.Entity
	if err := db.Where("id = ?", c.Param("id")).First(&dataEntity).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	db.Delete(&dataEntity)
	c.JSON(http.StatusOK, gin.H{"data": true})
}

func UpdateRepositoryById(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// Attempt to Fetch Models
	var dataEntity models.Entity
	if err := db.Where("id = ?", c.Param("id")).First(&dataEntity).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	// Validate Input
	var input models.UpdateEntityInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Model(&dataEntity).Updates(input)
	c.JSON(http.StatusOK, gin.H{"data": dataEntity})
}
