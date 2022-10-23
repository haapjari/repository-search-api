package commit

import (
	"github.com/gin-gonic/gin"
)

type CommitController struct {
	Handler        *Handler
	Context        *gin.Context
	GitHubApiToken string
	GitHubName     string
}

func GetCommits(c *gin.Context) {
	h := NewHandler(c)
	h.HandleGetCommits()
}

func GetCommitById(c *gin.Context) {
	h := NewHandler(c)
	h.HandleGetCommitById()
}

func CreateCommit(c *gin.Context) {
	h := NewHandler(c)
	h.HandleCreateCommit()
}

func DeleteCommitById(c *gin.Context) {
	h := NewHandler(c)
	h.HandleDeleteCommitById()
}

func UpdateCommitById(c *gin.Context) {
	h := NewHandler(c)
	h.HandleUpdateCommitById()
}
