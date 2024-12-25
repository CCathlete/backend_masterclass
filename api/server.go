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
	authRouter gin.IRoutes
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

	// We wrap our router with a middleware for all routes in the "all routes group".
	s.authRouter = s.router.Group("/").Use(authMiddleware(s.tokenMaker))
	// ------------------No Auth Middleware--------------------
	routeUser(s)
	// ------------------With Auth Middleware------------------
	routeAccount(s)
	routeTransfer(s)

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
	s.authRouter.POST("/accounts", s.createAccount)
	s.authRouter.POST("/accounts/updbalance", s.updateAccountBalance)
	s.authRouter.POST("/accounts/setbalance", s.updateAccount)
	// GET Requests:
	s.authRouter.GET("/accounts/", s.listAccounts)
	s.authRouter.GET("/accounts/:id", s.getAccount)
	s.authRouter.GET("/accounts/forupdate/:id", s.getAccountForUpdate)
	s.authRouter.GET("/accounts/delete/:id", s.deleteAccount)
}

func routeTransfer(s *Server) {
	// POST requests:
	s.authRouter.POST("/transfers", s.createTransfer)
	s.authRouter.POST("/transfers/updamount", s.updateTransfer)
	s.authRouter.POST("/transfers/fromaccount", s.getTransfersFromAccount)
	// GET Requests:
	s.authRouter.GET("/transfers", s.listTransfers)
	s.authRouter.GET("/transfers/:id", s.getTransfer)
	s.authRouter.GET("/transfers/delete/:id", s.deleteTransfer)
}

func routeUser(s *Server) {
	// POST requests:
	s.router.POST("/users", s.createUser)
	s.router.POST("/users/login", s.loginUser)
	s.router.POST("/users/updusername", s.updateUser)
	// GET Requests:
	s.router.GET("/users", s.listUsers)
	s.router.GET("/users/:username", s.getUser)
	s.router.GET("/users/delete/:username", s.deleteUser)
}
