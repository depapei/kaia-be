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

	result := DataAccess.DB.Find(&transactions).Limit(50)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, res.Fail{
			Message: result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, res.Success{
		Success: true,
		Data:    transactions,
	})
}
