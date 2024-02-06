package srv

import (
	"github.com/gin-gonic/gin"
	"github.com/haapjari/glass-api/api"
	"github.com/haapjari/glass-api/internal/pkg/svc"
)

// TODO
func (s *Server) GetApiV1RepositoriesSearchPages(ctx *gin.Context, params api.GetApiV1RepositoriesSearchPagesParams) {
	ctx.JSON(200, "Succesful Query to Repositories Search Pages")

	service := svc.NewRepositoryService(s.log)
	service.SearchPages(params.Language, params.Stars)
}

// TODO
func (s *Server) GetApiV1RepositoriesSearchPageNumber(ctx *gin.Context, pageNumber int,
	params api.GetApiV1RepositoriesSearchPageNumberParams) {
	ctx.JSON(200, "TODO")
}
