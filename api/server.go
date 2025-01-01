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
	Store sqlc.Store
	// Responsible for creating the context for each route.
	// It will automatically send the context to the handler functions.
	Router     *gin.Engine
	AuthRouter gin.IRoutes
	TokenMaker token.Maker
	Config     u.Config
}

func NewServer(store sqlc.Store, config u.Config, maker token.Maker,
) (s *Server) {
	s = &Server{
		Store:      store,
		Config:     config,
		TokenMaker: maker,
		Router:     gin.Default(),
	}
	if validationEngine, ok := binding.Validator.Engine().(*validator.Validate); ok {
		validationEngine.RegisterValidation("validcurrency", validCurrency)
	}

	// We wrap our router with a middleware for all routes in the "all routes group".
	s.AuthRouter = s.Router.Group("/").Use(authMiddleware(s.TokenMaker))
	// ------------------No Auth Middleware--------------------
	routeUser(s)
	// ------------------With Auth Middleware------------------
	routeAccount(s)
	routeTransfer(s)

	return
}

func (server *Server) Start(address string) (err error) {
	return server.Router.Run(address)
}

func routeAccount(s *Server) {
	// POST requests:
	s.AuthRouter.POST("/accounts", s.createAccount)
	// TODO: Change to PATCH.
	s.AuthRouter.POST("/accounts/updbalance", s.updateAccountBalance)
	// s.AuthRouter.POST("/accounts/setbalance", s.updateAccount)

	// GET Requests:
	s.AuthRouter.GET("/accounts", s.listAccounts)
	s.AuthRouter.GET("/accounts/:id", s.getAccount)
	s.AuthRouter.GET("/accounts/forupdate/:id", s.getAccountForUpdate)
	// TODO: Change to DELETE.
	s.AuthRouter.GET("/accounts/delete/:id", s.deleteAccount)
}

func routeTransfer(s *Server) {
	// POST requests:
	s.AuthRouter.POST("/transfers", s.createTransfer)
	// TODO: Change to PATCH.
	s.AuthRouter.POST("/transfers/updamount", s.updateTransfer)
	s.AuthRouter.POST("/transfers/fromaccount", s.getTransfersFromAccount)

	// GET Requests:
	s.AuthRouter.GET("/transfers", s.listTransfers)
	s.AuthRouter.GET("/transfers/:id", s.getTransfer)
	// TODO: Change to DELETE.
	s.AuthRouter.GET("/transfers/delete/:id", s.deleteTransfer)
}

func routeUser(s *Server) {
	// POST requests:
	s.Router.POST("/users", s.createUser)
	s.Router.POST("/users/login", s.loginUser)
	// TODO: Change to PATCH.
	s.Router.POST("/users/updusername", s.updateUser)
	// GET Requests:
	s.Router.GET("/users", s.listUsers)
	s.Router.GET("/users/:username", s.getUser)
	// TODO: Change to DELETE.
	s.Router.GET("/users/delete/:username", s.deleteUser)
}
