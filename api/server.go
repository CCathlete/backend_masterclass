package api

import (
	"backend-masterclass/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Serves all HTTP requests for our banking service.

type Server struct {
	store sqlc.Store
	// Responsible for creating the context for each route.
	// It will automatically send the context to the handler functions.
	Router *gin.Engine
}

func NewServer(store sqlc.Store) (s *Server) {
	s = &Server{
		store:  store,
		Router: gin.Default(),
	}
	if validationEngine, ok := binding.Validator.Engine().(*validator.Validate); ok {
		validationEngine.RegisterValidation("validcurrency", validCurrency)
	}

	routeAccount(s)
	routeTransfer(s)
	routeUser(s)

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
	// POST requests:
	s.Router.POST("/transfers", s.createTransfer)
	s.Router.POST("/transfers/updamount", s.updateTransfer)
	// GET Requests:
	s.Router.GET("/transfers", s.listTransfers)
	s.Router.GET("/transfers/:id", s.getTransfer)
	s.Router.GET("/transfers/delete/:id", s.deleteTransfer)
}

func routeUser(s *Server) {
	// POST requests:
	s.Router.POST("/users", s.createUser)
	s.Router.POST("/users/updusername", s.updateUser)
	// GET Requests:
	s.Router.GET("/users", s.listUsers)
	s.Router.GET("/users/:username", s.getUser)
	s.Router.GET("/users/delete/:username", s.deleteUser)
}
