package data

import (
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
