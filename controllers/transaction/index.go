package transaction

import (
	DataAccess "KAIA-BE/db"
	"KAIA-BE/model"
	res "KAIA-BE/responses"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {

	var transactions []model.HeaderTransaction

	result := DataAccess.DB.Preload("DetailTransaction.ProductSlice.Product").Find(&transactions).Limit(50)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, res.Fail{
			Message: result.Error.Error(),
		})
		return
	}

	var trx_res []TransactionInput
	for _, transaction := range transactions {
		var trx_items []ItemsInput
		for _, dt_trx := range transaction.DetailTransaction {
			trx_items = append(trx_items, ItemsInput{
				ID:       dt_trx.ProductsliceID,
				Name:     dt_trx.ProductSlice.Product.Name,
				Quantity: int64(dt_trx.Quantity),
				Price:    dt_trx.ProductSlice.Price,
				Slices:   dt_trx.ProductSlice.Slice,
			})
		}

		trx_res = append(trx_res, TransactionInput{
			ID:            transaction.ID,
			UserID:        transaction.CreatedBy.String(),
			PostalCode:    transaction.Postalcode,
			Address:       transaction.Address,
			City:          transaction.City,
			CustomerEmail: transaction.Customeremail,
			CustomerName:  transaction.Customername,
			TotalPrice:    float64(transaction.Totalprice),
			Items:         trx_items,
		})
	}

	c.JSON(http.StatusOK, res.Success{
		Success: true,
		Data:    trx_res,
	})
}
