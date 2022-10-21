package commit

import (
	"net/http"

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
	e := h.HandleGetCommits()

	c.JSON(http.StatusOK, gin.H{"data": e})
}

func GetCommitById(c *gin.Context) {
	h := NewHandler(c)
	e := h.HandleGetCommits()

	c.JSON(http.StatusOK, gin.H{"data": e})
}

func CreateCommit(c *gin.Context) {
	h := NewHandler(c)
	b := h.HandleCreateCommit()

	c.JSON(http.StatusOK, gin.H{"data: ": b})
}

func DeleteCommitById(c *gin.Context) {
	h := NewHandler(c)
	b := h.HandleDeleteCommitById()

	c.JSON(http.StatusOK, gin.H{"data ": b})
}

func UpdateCommitById(c *gin.Context) {
	h := NewHandler(c)
	e := h.HandleUpdateCommitById()

	c.JSON(http.StatusOK, gin.H{"data": e})
}
