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

func (h *Handler) HandleGetCommits() []models.Commit {
	var c []models.Commit
	h.DatabaseConnection.Find(&c)

	return c
}

func (h *Handler) HandleGetCommitById() models.Commit {
	var c models.Commit

	if err := h.DatabaseConnection.Where("id = ?", h.Context.Param("id")).First(&c).Error; err != nil {
		h.Context.JSON(http.StatusBadRequest, gin.H{"error": "entity not found!"})

		return c
	}

	return c
}

func (h *Handler) HandleCreateCommit() bool {
	var c models.Commit

	if err := h.Context.ShouldBindJSON(&c); err != nil {
		h.Context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return false
	}

	e := models.Commit{RepositoryId: c.RepositoryName, CommitDate: c.CommitDate, CommitUser: c.CommitDate}

	h.DatabaseConnection.Create(&e)

	return true
}

func (h *Handler) HandleDeleteCommitById() bool {
	var c models.Commit

	if err := h.DatabaseConnection.Where("id = ?", h.Context.Param("id")).First(&c).Error; err != nil {
		h.Context.JSON(http.StatusBadRequest, gin.H{"error": "record not found!"})
		return false
	}

	h.DatabaseConnection.Delete(&c)

	return true
}

func (h *Handler) HandleUpdateCommitById() models.Commit {
	var c models.Commit

	if err := h.DatabaseConnection.Where("id = ?", h.Context.Param("id")).First(&c).Error; err != nil {
		h.Context.JSON(http.StatusBadRequest, gin.H{"error": "record not found!"})

		return c
	}

	var i models.Commit

	if err := h.Context.ShouldBindJSON(&i); err != nil {
		h.Context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return c
	}

	h.DatabaseConnection.Model(&c).Updates(i)

	return c
}
