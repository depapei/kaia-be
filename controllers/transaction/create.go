package transaction

import (
	DataAccess "KAIA-BE/db"
	"KAIA-BE/model"
	res "KAIA-BE/responses"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ItemsInput struct {
	ID       string  `json:"id" binding:"required"`
	Quantity int64   `json:"quantity" binding:"required"`
	Price    float64 `json:"price" binding:"required"`
}

type TransactionInput struct {
	UserID        string       `json:"userId,omitempty"`
	PostalCode    string       `json:"postalCode" binding:"required"`
	Address       string       `json:"address" binding:"required"`
	City          string       `json:"city" binding:"required"`
	CustomerEmail string       `json:"customerEmail" binding:"required"`
	CustomerName  string       `json:"customerName" binding:"required"`
	Items         []ItemsInput `json:"items" binding:"required"`
}

func Create(c *gin.Context) {
	var input TransactionInput

	if err := c.ShouldBindJSON(&input); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, res.Fail{
			Message: err.Error(),
		})
	}

	totalPrice := sumPrice(input.Items)

	var trx_h model.HeaderTransaction
	trx_h = model.HeaderTransaction{
		Customername:  input.CustomerName,
		Customeremail: input.CustomerEmail,
		Address:       input.Address,
		City:          input.City,
		Postalcode:    input.PostalCode,
		Totalprice:    int32(totalPrice),
	}

	if len(input.UserID) > 0 && input.UserID != "null" {
		parsedId, err := uuid.Parse(input.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, res.Fail{
				Message: "Failed to input user_id",
			})
			return
		}
		trx_h.CreatedBy = &parsedId
	}

	dbTrx := DataAccess.DB.Begin()
	if err := dbTrx.Error; err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, res.Fail{
			Message: "Failed to begin transaction",
		})
		return
	}

	if err := dbTrx.Create(&trx_h).Error; err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, res.Fail{
			Message: "Failed to create Transaction",
		})
		return
	}

	var trx_d []model.DetailTransaction
	for _, item := range input.Items {
		trx_d = append(trx_d, model.DetailTransaction{
			TransactionID:  trx_h.ID,
			ProductsliceID: item.ID,
			Quantity:       float64(item.Quantity),
		})
	}

	if len(trx_d) > 0 {
		if err := dbTrx.Create(&trx_d).Error; err != nil {
			dbTrx.Rollback()
			fmt.Println(err.Error())
			c.JSON(http.StatusInternalServerError, res.Fail{
				Message: "Failed to create Detail Transaction",
			})
			return
		}
	}

	if err := dbTrx.Commit().Error; err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, res.Fail{
			Message: "Failed to commit Transaction",
		})
		return
	}

	c.JSON(http.StatusCreated, res.Success{
		Success: true,
		Data:    input,
	})
}

func sumPrice(items []ItemsInput) float64 {
	var sum float64

	for _, items := range items {
		sum += items.Price
	}

	return sum
}
