package data

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/haapjari/glass/pkg/models"
	"gorm.io/gorm"
)

// Handler handle's the communication between application and database.

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

func (h *Handler) HandleGetEntities() []models.Entity {
	var e []models.Entity
	h.DatabaseConnection.Find(&e)

	return e
}

func (h *Handler) HandleGetEntityById() models.Entity {
	var e models.Entity

	if err := h.DatabaseConnection.Where("id = ?", h.Context.Param("id")).First(&e).Error; err != nil {
		h.Context.JSON(http.StatusBadRequest, gin.H{"error": "entity not found!"})

		return e
	}

	return e
}
