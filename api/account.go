package api

import (
	"database/sql"
	"errors"
	"net/http"

	db "github.com/badermezzi/KubeGoBank/db/sqlc"
	"github.com/badermezzi/KubeGoBank/token"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createAccountRequest struct {
	Currency string `json:"currency" binding:"required,currency"`
}

func (server *Server) createAccount(context *gin.Context) {
	var req createAccountRequest

	err := context.ShouldBindJSON(&req)
	if err != nil {
		context.JSON(http.StatusBadRequest, errorResponce(err))
		return
	}

	authPayload := context.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.CreateAccountParams{
		Owner:    authPayload.Username,
		Balance:  0,
		Currency: req.Currency,
	}

	account, err := server.store.CreateAccount(context, arg)
	if err != nil {
		pqError, ok := err.(*pq.Error)
		if ok {
			switch pqError.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				context.JSON(http.StatusForbidden, errorResponce(err))
				return
			}
		}
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

	authPayload := context.MustGet(authorizationPayloadKey).(*token.Payload)

	if account.Owner != authPayload.Username {
		err := errors.New("account doesn't belong to the authenticated user")
		context.JSON(http.StatusUnauthorized, errorResponce(err))
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
