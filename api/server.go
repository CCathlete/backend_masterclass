package api

import (
	"backend-masterclass/db/sqlc"

	"github.com/gin-gonic/gin"
)

// Serves all HTTP requests for our banking service.

type Server struct {
	store  *sqlc.Store
	router *gin.Engine
}

func NewServer(store *sqlc.Store) (s *Server) {
	s = &Server{
		store:  store,
		router: gin.Default(),
	}

	s.router.POST("/accounts", s.createAccount)

	return
}

func errorResponse(err error) (resBody gin.H) {
	resBody = gin.H{
		"error": err.Error(),
	}

	return
}
