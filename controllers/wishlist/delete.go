package wishlist

import (
	DataAcces "KAIA-BE/db"
	"KAIA-BE/model"
	res "KAIA-BE/responses"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Delete(c *gin.Context) {

	var userId = c.Param("user_id")
	var productId = c.Param("product_id")

	if userId == "" || productId == "" {
		c.JSON(http.StatusBadRequest, res.Fail{
			Message: "User id atau produk id tidak boleh koson!",
		})
		return
	}

	data := model.Wishlist{
		UserID:    userId,
		ProductID: productId,
	}

	if err := DataAcces.DB.Model(&model.Wishlist{}).Delete(&data).Error; err != nil {
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
