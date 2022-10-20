// Handler handle's the communication between application and database.
package repository

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/haapjari/glass/pkg/models"
	"gorm.io/gorm"
)

type Handler struct {
	Context            *gin.Context
	DatabaseConnection *gorm.DB
}

func NewHandler(c *gin.Context) *Handler {
	h := new(Handler)

	h.Context = c
	h.DatabaseConnection = c.MustGet("db").(*gorm.DB)

	return h
}

func (h *Handler) HandleGetRepositories() []models.Repository {
	var e []models.Repository
	h.DatabaseConnection.Find(&e)

	return e
}

func (h *Handler) HandleGetRepositoryById() models.Repository {
	var e models.Repository

	if err := h.DatabaseConnection.Where("id = ?", h.Context.Param("id")).First(&e).Error; err != nil {
		h.Context.JSON(http.StatusBadRequest, gin.H{"error": "entity not found!"})

		return e
	}

	return e
}

func (h *Handler) HandleCreateRepository() bool {
	var i models.Repository

	if err := h.Context.ShouldBindJSON(&i); err != nil {
		h.Context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return false
	}

	e := models.Repository{Name: i.Name, Uri: i.Uri, CodeLines: i.CodeLines, IssueCount: i.IssueCount, LibraryCodeLines: i.LibraryCodeLines, RepositoryType: i.RepositoryType}

	h.DatabaseConnection.Create(&e)

	return true
}

func (h *Handler) HandleDeleteRepositoryById() bool {
	var r models.Repository

	if err := h.DatabaseConnection.Where("id = ?", h.Context.Param("id")).First(&r).Error; err != nil {
		h.Context.JSON(http.StatusBadRequest, gin.H{"error": "record not found!"})
		return false
	}

	h.DatabaseConnection.Delete(&r)

	return true
}

func (h *Handler) HandleUpdateRepositoryById() models.Repository {
	var dataEntity models.Repository

	if err := h.DatabaseConnection.Where("id = ?", h.Context.Param("id")).First(&dataEntity).Error; err != nil {
		h.Context.JSON(http.StatusBadRequest, gin.H{"error": "record not found!"})

		return dataEntity
	}

	// Validate Input
	var input models.Repository
	if err := h.Context.ShouldBindJSON(&input); err != nil {
		h.Context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return dataEntity
	}

	h.DatabaseConnection.Model(&dataEntity).Updates(input)

	return dataEntity
}
