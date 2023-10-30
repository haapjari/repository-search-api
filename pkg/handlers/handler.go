package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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

func (h *Handler) HandleGetRepositoryMetadata() {
	var count, err = strconv.Atoi(h.Context.Query("count"))
	utils.CheckErr(err)

	if h.Context.Query("type") == "go" {
		plugin := goplg.NewGoPlugin(h.Database)
		processedRepositories := plugin.GenerateRepositoryData(count)

		if processedRepositories != nil {
			h.Context.IndentedJSON(http.StatusOK, processedRepositories)
		}
	}
}
