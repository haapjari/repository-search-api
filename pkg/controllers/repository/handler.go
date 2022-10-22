package repository

import (
	"fmt"
	"net/http"
	"strconv"

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

func (h *Handler) HandleGetRepositories() {
	var e []models.Repository

	h.DatabaseConnection.Find(&e)

	h.Context.JSON(http.StatusOK, gin.H{"data": e})
}

func (h *Handler) HandleGetRepositoryById() {
	var e models.Repository

	if err := h.DatabaseConnection.Where("id = ?", h.Context.Param("id")).First(&e).Error; err != nil {
		h.Context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	h.Context.JSON(http.StatusOK, gin.H{"data": e})
}

func (h *Handler) HandleCreateRepository() {
	var i models.Repository

	if err := h.Context.ShouldBindJSON(&i); err != nil {
		h.Context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	r := models.Repository{RepositoryName: i.RepositoryName, RepositoryUrl: i.RepositoryUrl, OpenIssueCount: i.OpenIssueCount, ClosedIssueCount: i.ClosedIssueCount, OriginalCodebaseSize: i.OriginalCodebaseSize, LibraryCodebaseSize: i.LibraryCodebaseSize, RepositoryType: i.RepositoryType, PrimaryLanguage: i.PrimaryLanguage}

	h.DatabaseConnection.Create(&r)

	h.Context.JSON(http.StatusOK, gin.H{"data": r})
}

func (h *Handler) HandleDeleteRepositoryById() {
	var r models.Repository

	if err := h.DatabaseConnection.Where("id = ?", h.Context.Param("id")).First(&r).Error; err != nil {
		h.Context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	h.DatabaseConnection.Delete(&r)

	h.Context.JSON(http.StatusOK, gin.H{"succeed": r})
}

func (h *Handler) HandleUpdateRepositoryById() {
	var r models.Repository

	if err := h.DatabaseConnection.Where("id = ?", h.Context.Param("id")).First(&r).Error; err != nil {
		h.Context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	var i models.Repository

	if err := h.Context.ShouldBindJSON(&i); err != nil {
		h.Context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	h.DatabaseConnection.Model(&r).Updates(i)

	h.Context.JSON(http.StatusOK, gin.H{"data": r})
}

// TODO
func (h *Handler) HandleFetchRepositories() {
	count, _ := strconv.Atoi(h.Context.Query("amount"))
	lang := h.Context.Query("lang")

	// Query SourceGraph GraphQL API -> Fetch Metadata of the Repositories -> Save this to Database.
	// Filter out the Repositories, that are not GitHub -projects.
	// Encrich the data with GitHub GraphQL API.
	// Export the data as a .CSV

	// graphQlQueryStr := `query {
	// 	search(query: "(lang:javascript OR lang:typescript) AND select:repo AND repohasfile:package.json AND count:1000", version: V2) {
	// 	  results {
	// 		  repositories {
	// 			  name
	// 			  stars
	// 		  }
	// 	  }
	// 	}
	//   }`

	fmt.Println("")
	fmt.Println(count, lang)
}
