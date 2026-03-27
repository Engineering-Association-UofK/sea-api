package response

import "github.com/gin-gonic/gin"

type TransactionResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	ID      int64  `json:"id"`
}

func NewTransactionResponse(status int, message string, id int64, c *gin.Context) {
	c.JSON(status, TransactionResponse{
		Status:  status,
		Message: message,
		ID:      id,
	})
}
