package repository

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/haapjari/glass/pkg/models"
	"github.com/haapjari/glass/pkg/plugins/goplg"
	"github.com/haapjari/glass/pkg/utils"
	"gorm.io/gorm"
)

type Handler struct {
	Context  *gin.Context
	Database *gorm.DB
}

func NewHandler(c *gin.Context) *Handler {
	h := new(Handler)

	h.Context = c
	h.Database = c.MustGet("db").(*gorm.DB)

	return h
}

func (h *Handler) HandleGetRepositories() {
	var e []models.Repository

	h.Database.Find(&e)

	h.Context.JSON(http.StatusOK, gin.H{"data": e})
}

func (h *Handler) HandleGetRepositoryById() {
	var e models.Repository

	if err := h.Database.Where("id = ?", h.Context.Param("id")).First(&e).Error; err != nil {
		h.Context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	h.Context.JSON(http.StatusOK, gin.H{"data": e})
}

func (h *Handler) HandleCreateRepository() {
	var i models.Repository

	if err := h.Context.ShouldBindJSON(&i); err != nil {
		h.Context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	r := models.Repository{LatestRelease: i.LatestRelease, RepositoryName: i.RepositoryName, RepositoryUrl: i.RepositoryUrl, CommitCount: i.CommitCount, OpenIssueCount: i.OpenIssueCount, ClosedIssueCount: i.ClosedIssueCount, OriginalCodebaseSize: i.OriginalCodebaseSize, LibraryCodebaseSize: i.LibraryCodebaseSize, RepositoryType: i.RepositoryType, PrimaryLanguage: i.PrimaryLanguage, CreationDate: i.CreationDate, StargazerCount: i.StargazerCount, LicenseInfo: i.LicenseInfo}

	h.Database.Create(&r)

	h.Context.JSON(http.StatusOK, gin.H{"data": r})
}

func (h *Handler) HandleDeleteRepositoryById() {
	var r models.Repository

	if err := h.Database.Where("id = ?", h.Context.Param("id")).First(&r).Error; err != nil {
		h.Context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	h.Database.Delete(&r)

	h.Context.JSON(http.StatusOK, gin.H{"succeed": r})
}

func (h *Handler) HandleUpdateRepositoryById() {
	var r models.Repository

	if err := h.Database.Where("id = ?", h.Context.Param("id")).First(&r).Error; err != nil {
		h.Context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	var i models.Repository

	if err := h.Context.ShouldBindJSON(&i); err != nil {
		h.Context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	h.Database.Model(&r).Updates(i)

	h.Context.JSON(http.StatusOK, gin.H{"data": r})
}

func (h *Handler) FetchRepositoryMetadata() {
	var count, err = strconv.Atoi(h.Context.Query("count"))
	utils.CheckErr(err)

	if h.Context.Query("type") == "go" {
		plugin := goplg.NewGoPlugin(h.Database)

		plugin.GetRepositoryMetadata(count)
	}
}
