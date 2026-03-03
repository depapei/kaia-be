package wishlist

import (
	DataAcces "KAIA-BE/db"
	"KAIA-BE/model"
	res "KAIA-BE/responses"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WishlistInput struct {
	ProductID string `json:"productId" binding:"required"`
	UserId    string `json:"userId" binding:"required"`
}

func Create(c *gin.Context) {
	var input WishlistInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, res.Fail{
			Message: err.Error(),
		})
		return
	}

	data := model.Wishlist{
		UserID:    input.UserId,
		ProductID: input.ProductID,
	}
	if err := DataAcces.DB.Model(&model.Wishlist{}).Create(&data).Error; err != nil {
		c.JSON(http.StatusInternalServerError, res.Fail{
			Message: "Terdapat kesalahan server",
		})
		return
	}

	c.JSON(http.StatusOK, res.Success{
		Success: true,
		Data:    data,
	})
}
