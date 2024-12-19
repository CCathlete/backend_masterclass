package api

import (
	"backend-masterclass/db/sqlc"

	"github.com/gin-gonic/gin"
)

// Serves all HTTP requests for our banking service.

type Server struct {
	store  sqlc.Store
	Router *gin.Engine
}

func NewServer(store sqlc.Store) (s *Server) {
	s = &Server{
		store:  store,
		Router: gin.Default(),
	}
	routeAccount(s)
	routeTransfer(s)

	return
}

func (server *Server) Start(address string) (err error) {
	return server.Router.Run(address)
}

func errorResponse(err error) (resBody gin.H) {
	return gin.H{
		"error": err.Error(),
	}
}

func routeAccount(s *Server) {
	// POST requests:
	s.Router.POST("/accounts", s.createAccount)
	s.Router.POST("/accounts/updbalance", s.updateAccountBalance)
	s.Router.POST("/accounts/setbalance", s.updateAccount)
	// GET Requests:
	s.Router.GET("/accounts/", s.listAccounts)
	s.Router.GET("/accounts/:id", s.getAccount)
	s.Router.GET("/accounts/forupdate/:id", s.getAccountForUpdate)
	s.Router.GET("/accounts/delete/:id", s.deleteAccount)
}

func routeTransfer(s *Server) {
	// TODO: Add routes for other transfer operations.
	// POST requests:
	s.Router.POST("/transfers", s.createTransfer)
	// GET Requests:
	s.Router.GET("/transfers/:id", s.getTransfer)
}
