// Handler handle's the communication between application and database.
package repository

import (
	"fmt"
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

func (h *Handler) HandleGetRepositories() []models.Entity {
	var e []models.Entity
	h.DatabaseConnection.Find(&e)

	return e
}

func (h *Handler) HandleGetRepositoryById() models.Entity {
	var e models.Entity

	if err := h.DatabaseConnection.Where("id = ?", h.Context.Param("id")).First(&e).Error; err != nil {
		h.Context.JSON(http.StatusBadRequest, gin.H{"error": "entity not found!"})

		return e
	}

	return e
}

func (h *Handler) HandleCreateRepository() {
	fmt.Println("TODO")
}

func (h *Handler) HandleDeleteRepositoryById() {
	fmt.Println("TODO")
}

func (h *Handler) HandleUpdateRepositoryById() {
	fmt.Println("TODO")
}
