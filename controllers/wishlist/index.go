package wishlist

import (
	DataAccess "KAIA-BE/db"
	"KAIA-BE/model"
	res "KAIA-BE/responses"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SliceOpt struct {
	Slices string  `json:"slices"`
	Price  float64 `json:"price"`
}

type WishlistResponse struct {
	ID              string     `json:"id"`
	Name            string     `json:"name"`
	Price           float64    `json:"price"`
	Category        string     `json:"category"`
	Description     string     `json:"desc"`
	LongDescription string     `json:"longDesc"`
	Image           string     `json:"image"`
	SlicesOptions   []SliceOpt `json:"sliceOptions"`
}

func Index(c *gin.Context) {
	var wishlists []model.Wishlist
	var user_id = c.Param("user_id")

	if user_id == "" {
		c.JSON(http.StatusBadRequest, res.Fail{
			Message: "User id tidak boleh kosong",
		})
		return
	}

	result := DataAccess.DB.Preload("Product.ProductSlices").Find(&wishlists).Where("user_id = ?", user_id).Limit(50)

	var response []WishlistResponse
	for _, wishlist := range wishlists {

		var sliceOpt []SliceOpt
		for _, productSlice := range wishlist.Product.ProductSlices {
			sliceOpt = append(sliceOpt, SliceOpt{
				Slices: productSlice.Slice,
				Price:  productSlice.Price,
			})
		}

		response = append(response, WishlistResponse{
			Name:            wishlist.Product.Name,
			ID:              wishlist.ProductID,
			Price:           float64(wishlist.Product.Price),
			Description:     wishlist.Product.Description,
			LongDescription: wishlist.Product.Longdesc,
			Image:           wishlist.Product.Image,
			Category:        wishlist.Product.Category,
			SlicesOptions:   sliceOpt,
		})
	}

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, res.Fail{
			Message: result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, res.Success{
		Success: true,
		Data:    response,
	})
}
