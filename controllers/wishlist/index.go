package wishlist

import (
	DataAccess "KAIA-BE/db"
	"KAIA-BE/model"
	res "KAIA-BE/responses"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	var wishlists []model.Wishlist
	var user_id = c.Query("user_id")

	if user_id == "" {
		c.JSON(http.StatusBadRequest, res.Fail{
			Message: "user_id cannot be empty",
		})
		return
	}

	result := DataAccess.DB.Find(&wishlists).Where("user_id = ?", user_id).Limit(50)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, res.Fail{
			Message: result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, res.Success{
		Success: true,
		Data:    wishlists,
	})
}
