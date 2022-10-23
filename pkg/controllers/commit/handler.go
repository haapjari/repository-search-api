package commit

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

func (h *Handler) HandleGetCommits() {
	var c []models.Commit

	h.DatabaseConnection.Find(&c)

	h.Context.JSON(http.StatusOK, gin.H{"data": c})
}

func (h *Handler) HandleGetCommitById() {
	var c models.Commit

	if err := h.DatabaseConnection.Where("id = ?", h.Context.Param("id")).First(&c).Error; err != nil {
		h.Context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	h.Context.JSON(http.StatusOK, gin.H{"data": c})
}

func (h *Handler) HandleCreateCommit() {
	var i models.Commit

	if err := h.Context.ShouldBindJSON(&i); err != nil {
		h.Context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	var c models.Commit

	c.RepositoryName = i.RepositoryName
	c.CommitDate = i.CommitDate
	c.CommitUser = i.CommitUser

	h.DatabaseConnection.Create(&c)

	h.Context.JSON(http.StatusOK, gin.H{"data": c})
}

func (h *Handler) HandleDeleteCommitById() {
	var c models.Commit

	if err := h.DatabaseConnection.Where("id = ?", h.Context.Param("id")).First(&c).Error; err != nil {
		h.Context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	h.DatabaseConnection.Delete(&c)

	h.Context.JSON(http.StatusOK, gin.H{"succeed": c})
}

func (h *Handler) HandleUpdateCommitById() {
	var c models.Commit

	if err := h.DatabaseConnection.Where("id = ?", h.Context.Param("id")).First(&c).Error; err != nil {
		h.Context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	var i models.Commit

	if err := h.Context.ShouldBindJSON(&i); err != nil {
		h.Context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	h.DatabaseConnection.Model(&c).Updates(i)

	h.Context.JSON(http.StatusOK, gin.H{"data": c})
}
