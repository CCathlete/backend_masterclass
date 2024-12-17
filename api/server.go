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
	s.router.POST("/accounts/updbalance", s.updateAccountBalance)
	s.router.GET("/accounts/:id", s.getAccount)
	s.router.GET("/accounts/delete/:id", s.deleteAccount)
	s.router.GET("/accounts/", s.listAccounts)

	return
}

func (server *Server) Start(address string) (err error) {
	return server.router.Run(address)
}

func errorResponse(err error) (resBody gin.H) {
	return gin.H{
		"error": err.Error(),
	}
}
