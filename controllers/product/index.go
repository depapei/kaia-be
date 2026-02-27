package product

import (
	DataAccess "KAIA-BE/db"
	"KAIA-BE/model"
	res "KAIA-BE/responses"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {

	var products []model.Product

	result := DataAccess.DB.Limit(50).
		Preload("Admin").
		Find(&products)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, res.Fail{
			Message: result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, res.Success{
		Success: true,
		Data:    products,
	})
}
