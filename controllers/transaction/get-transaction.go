package transaction

import (
	DataAccess "KAIA-BE/db"
	"KAIA-BE/model"
	res "KAIA-BE/responses"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetTransaction(c *gin.Context) {
	var data []model.HeaderTransaction
	userId := c.Query("userId")

	da := DataAccess.DB.Model(&model.HeaderTransaction{})

	result := da.Where("created_by = ?", userId).Find(&data)
	if err := result.Error; err != nil {
		c.JSON(http.StatusInternalServerError, res.Fail{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, res.Success{
		Success: true,
		Data:    result,
	})
}
