package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Server implements
type Server struct{}

func NewServer() *Server {
	return &Server{}
}

// GetHello handles the GET request at /hello endpoint.
func (s *Server) GetApiV1Hello(c *gin.Context) {
	// Implement the logic here. For example, return a hello message:
	c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
}
