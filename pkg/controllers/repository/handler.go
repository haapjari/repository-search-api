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

	e := models.Repository{RepositoryName: i.RepositoryName, RepositoryUrl: i.RepositoryUrl, OpenIssueCount: i.OpenIssueCount, ClosedIssueCount: i.ClosedIssueCount, OriginalCodebaseSize: i.OriginalCodebaseSize, LibraryCodebaseSize: i.LibraryCodebaseSize, RepositoryType: i.RepositoryType, PrimaryLanguage: i.PrimaryLanguage}

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

	var input models.Repository
	if err := h.Context.ShouldBindJSON(&input); err != nil {
		h.Context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return dataEntity
	}

	h.DatabaseConnection.Model(&dataEntity).Updates(input)

	return dataEntity
}

// TODO
func (h *Handler) HandleFetchRepositories() bool {
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

	return true
}
