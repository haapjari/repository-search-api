package srv

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/haapjari/repository-metadata-aggregator/api"
	"github.com/haapjari/repository-metadata-aggregator/internal/pkg/svc"
	"github.com/haapjari/repository-metadata-aggregator/internal/pkg/utils"
)

type ErrorResponse struct {
	Error string `json:"error,omitempty"`
}

func (s *Server) GetApiV1RepositoriesSearch(ctx *gin.Context, params api.GetApiV1RepositoriesSearchParams) {
	var (
		stars             = params.Stars
		language          = params.Language
		lastCreationDate  = params.LastCreationDate
		firstCreationDate = params.FirstCreationDate
		order             = fmt.Sprintf("%v", *params.Order)
		token, err        = utils.ParseBearer(ctx.GetHeader("Authorization"))
	)

	if err != nil {
		s.log.Warnf("request provided malformed or empty authorization header: %s", err.Error())
	}

	service, err := svc.NewRepositorySearchService(s.log, s.config, token)
	if err != nil {
		s.log.Errorf("error, while creating the repository search service: %s", err.Error())
		ctx.IndentedJSON(500, &ErrorResponse{Error: err.Error()})
		return
	}

	repositories, status, err := service.Search(&svc.SearchOptions{
		Stars:             stars,
		Language:          language,
		LastCreationDate:  lastCreationDate,
		FirstCreationDate: firstCreationDate,
		Order:             order,
	})
	if err != nil {
		s.log.Errorf("error, while executing 'search': %s", err.Error())
		ctx.IndentedJSON(status, &ErrorResponse{Error: err.Error()})
		return
	}

	res := &api.GetApiV1RepositoriesSearchResponse{
		JSON200: &struct {
			Items      []api.Repository `json:"items"`
			TotalCount int              `json:"total_count"`
		}{
			TotalCount: len(repositories),
			Items:      repositories,
		},
	}

	ctx.IndentedJSON(status, res.JSON200)
}

func (s *Server) GetApiV1RepositoriesSearchFirstCreationDate(ctx *gin.Context,
	params api.GetApiV1RepositoriesSearchFirstCreationDateParams) {
	var (
		stars      = params.Stars
		language   = params.Language
		token, err = utils.ParseBearer(ctx.GetHeader("Authorization"))
	)

	if err != nil {
		s.log.Warnf("request provided malformed or empty authorization header: %s", err.Error())
	}

	service, err := svc.NewRepositorySearchService(s.log, s.config, token)
	if err != nil {
		s.log.Errorf("error, while creating the repository search service: %s", err.Error())
		ctx.IndentedJSON(500, &ErrorResponse{Error: err.Error()})
		return
	}

	firstCreationDate, status, err := service.FirstCreationDate(language, stars)
	if err != nil {
		s.log.Errorf("error, while querying for first creation date: %s", err.Error())
		ctx.IndentedJSON(status, &ErrorResponse{Error: err.Error()})
		return
	}

	res := &api.GetApiV1RepositoriesSearchFirstCreationDateResponse{
		JSON200: &struct {
			FirstCreationDate *string `json:"firstCreationDate,omitempty"`
		}{
			FirstCreationDate: &firstCreationDate,
		},
	}

	ctx.IndentedJSON(status, res.JSON200)
}

func (s *Server) GetApiV1RepositoriesSearchLastCreationDate(ctx *gin.Context, params api.GetApiV1RepositoriesSearchLastCreationDateParams) {
	var (
		stars      = params.Stars
		language   = params.Language
		token, err = utils.ParseBearer(ctx.GetHeader("Authorization"))
	)

	if err != nil {
		s.log.Warnf("request provided malformed or empty authorization header: %s", err.Error())
	}

	service, err := svc.NewRepositorySearchService(s.log, s.config, token)
	if err != nil {
		s.log.Errorf("error, while creating the repository search service: %s", err.Error())
		ctx.IndentedJSON(500, &ErrorResponse{Error: err.Error()})
		return
	}

	lastCreationDate, status, err := service.LastCreationDate(language, stars)
	if err != nil {
		s.log.Errorf("error, while querying for first creation date: %s", err.Error())
		ctx.IndentedJSON(status, &ErrorResponse{Error: err.Error()})
		return
	}

	res := &api.GetApiV1RepositoriesSearchLastCreationDateResponse{
		JSON200: &struct {
			LastCreationDate *string `json:"lastCreationDate,omitempty"`
		}{
			LastCreationDate: &lastCreationDate,
		},
	}

	ctx.IndentedJSON(status, res.JSON200)
}
