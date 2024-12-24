package api

import (
	"backend-masterclass/db/sqlc"
	"backend-masterclass/token"
	u "backend-masterclass/util"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Serves all HTTP requests for our banking service.

type Server struct {
	store sqlc.Store
	// Responsible for creating the context for each route.
	// It will automatically send the context to the handler functions.
	router     *gin.Engine
	tokenMaker token.Maker
	config     u.Config
}

func NewServer(store sqlc.Store, config u.Config, maker token.Maker,
) (s *Server) {
	s = &Server{
		store:      store,
		config:     config,
		tokenMaker: maker,
		router:     gin.Default(),
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
	return server.router.Run(address)
}

func errorResponse(err error) (resBody gin.H) {
	return gin.H{
		"error": err.Error(),
	}
}

func routeAccount(s *Server) {
	// POST requests:
	s.router.POST("/accounts", s.createAccount)
	s.router.POST("/accounts/updbalance", s.updateAccountBalance)
	s.router.POST("/accounts/setbalance", s.updateAccount)
	// GET Requests:
	s.router.GET("/accounts/", s.listAccounts)
	s.router.GET("/accounts/:id", s.getAccount)
	s.router.GET("/accounts/forupdate/:id", s.getAccountForUpdate)
	s.router.GET("/accounts/delete/:id", s.deleteAccount)
}

func routeTransfer(s *Server) {
	// POST requests:
	s.router.POST("/transfers", s.createTransfer)
	s.router.POST("/transfers/updamount", s.updateTransfer)
	// GET Requests:
	s.router.GET("/transfers", s.listTransfers)
	s.router.GET("/transfers/:id", s.getTransfer)
	s.router.GET("/transfers/delete/:id", s.deleteTransfer)
}

func routeUser(s *Server) {
	// POST requests:
	s.router.POST("/users", s.createUser)
	s.router.POST("/users/updusername", s.updateUser)
	// GET Requests:
	s.router.GET("/users", s.listUsers)
	s.router.GET("/users/:username", s.getUser)
	s.router.GET("/users/delete/:username", s.deleteUser)
}
