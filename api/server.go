package api

import (
	db "github.com/badermezzi/KubeGoBank/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	v, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		v.RegisterValidation("currency", validCurrency)
	}

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccount)
	router.DELETE("/accounts/:id", server.deleteAccount)
	router.PATCH("/accounts/:id", server.updateAccount)

	router.POST("/transfers", server.createTransfer)

	server.router = router

	return server

}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponce(err error) gin.H {
	return gin.H{"error": err.Error()}
}
