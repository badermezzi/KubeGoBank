package api

import (
	"database/sql"
	"fmt"
	"net/http"

	db "github.com/badermezzi/KubeGoBank/db/sqlc"
	"github.com/gin-gonic/gin"
)

type transferRequest struct {
	FromAccountID int64 `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64 `json:"to_account_id" binding:"required,min=1"`
	Amount        int64 `json:"amount" binding:"required,gt=0"`
	// Currency      string `json:"currency" binding:"required,oneof=USD EUR"`
	Currency string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(context *gin.Context) {
	var req transferRequest

	err := context.ShouldBindJSON(&req)
	if err != nil {
		context.JSON(http.StatusBadRequest, errorResponce(err))
		return
	}

	if !server.validAccount(context, req.FromAccountID, req.Currency) {
		return
	}

	if !server.validAccount(context, req.ToAccountID, req.Currency) {
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	result, err := server.store.TransferTx(context, arg)
	if err != nil {
		context.JSON(http.StatusInternalServerError, errorResponce(err))
		return
	}

	context.JSON(http.StatusOK, result)

}

func (server *Server) validAccount(context *gin.Context, accountID int64, currency string) bool {
	account, err := server.store.GetAccount(context, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			context.JSON(http.StatusNotFound, errorResponce(err))
			return false
		}

		context.JSON(http.StatusInternalServerError, errorResponce(err))
		return false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", account.ID, account.Currency, currency)
		context.JSON(http.StatusBadRequest, errorResponce(err))
		return false
	}

	return true
}
