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

type deleteAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) deleteAccount(context *gin.Context) {
	var req deleteAccountRequest

	err := context.ShouldBindUri(&req)
	if err != nil {
		context.JSON(http.StatusBadRequest, errorResponce(err))
		return
	}

	// checking if id valid and account exist
	_, err = server.store.GetAccount(context, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			context.JSON(http.StatusNotFound, errorResponce(err))
			return
		}

		context.JSON(http.StatusInternalServerError, errorResponce(err))
		return
	}

	// delete account
	err = server.store.DeleteAccount(context, req.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponce(err))
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "account deleted"})
}

type updateAccountRequest struct {
	ID      int64 `json:"id" binding:"required,min=1"`
	Balance int64 `json:"balance" binding:"required,min=1"`
}

func (server *Server) updateAccount(context *gin.Context) {
	var req updateAccountRequest

	err := context.ShouldBindJSON(&req)
	if err != nil {
		context.JSON(http.StatusBadRequest, errorResponce(err))
		return
	}

	arg := db.UpdateAccountParams{
		ID:      req.ID,
		Balance: req.Balance,
	}

	account, err := server.store.UpdateAccount(context, arg)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponce(err))
		return
	}

	context.JSON(http.StatusOK, account)
}
