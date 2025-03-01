package api

import (
	"database/sql"
	"net/http"

	db "github.com/badermezzi/KubeGoBank/db/sqlc"
	"github.com/gin-gonic/gin"
)

type createAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR"`
}

func (server *Server) createAccount(context *gin.Context) {
	var req createAccountRequest

	err := context.ShouldBindJSON(&req)
	if err != nil {
		context.JSON(http.StatusBadRequest, errorResponce(err))
		return
	}

	arg := db.CreateAccountParams{
		Owner:    req.Owner,
		Balance:  0,
		Currency: req.Currency,
	}

	account, err := server.store.CreateAccount(context, arg)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponce(err))
		return
	}

	context.JSON(http.StatusOK, account)

}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(context *gin.Context) {
	var req getAccountRequest

	err := context.ShouldBindUri(&req)
	if err != nil {
		context.JSON(http.StatusBadRequest, errorResponce(err))
		return
	}

	account, err := server.store.GetAccount(context, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			context.JSON(http.StatusNotFound, errorResponce(err))
			return
		}

		context.JSON(http.StatusInternalServerError, errorResponce(err))
		return
	}

	context.JSON(http.StatusOK, account)

}

type listAccountRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listAccount(context *gin.Context) {
	var req listAccountRequest

	err := context.ShouldBindQuery(&req)
	if err != nil {
		context.JSON(http.StatusBadRequest, errorResponce(err))
		return
	}

	arg := db.ListAccountsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	accounts, err := server.store.ListAccounts(context, arg)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponce(err))
		return
	}

	context.JSON(http.StatusOK, accounts)

}
