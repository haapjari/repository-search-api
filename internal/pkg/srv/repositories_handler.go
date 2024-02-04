package srv

import (
	"github.com/gin-gonic/gin"
	"github.com/haapjari/glass-api/api"
)

// TODO
func (s *Server) GetApiV1RepositoriesSearchPages(ctx *gin.Context, params api.GetApiV1RepositoriesSearchPagesParams) {
	ctx.JSON(200, "TODO")
}

// TODO
func (s *Server) GetApiV1RepositoriesSearchPageNumber(ctx *gin.Context, pageNumber int,
	params api.GetApiV1RepositoriesSearchPageNumberParams) {
	ctx.JSON(200, "TODO")
}
