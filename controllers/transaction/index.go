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

	da := DataAccess.DB.Model(&model.HeaderTransaction{})
	result := da.Preload("DetailTransaction.ProductSlice.Product").
		Limit(200).
		Find(&transactions)

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
			PostalCode:    transaction.Postalcode,
			Address:       transaction.Address,
			City:          transaction.City,
			CustomerEmail: transaction.Customeremail,
			CustomerName:  transaction.Customername,
			TotalPrice:    float64(transaction.Totalprice),
			Items:         trx_items,
			CreatedAt:     transaction.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	c.JSON(http.StatusOK, res.Success{
		Success: true,
		Data:    trx_res,
	})
}
