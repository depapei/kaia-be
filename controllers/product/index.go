package product

import (
	DataAccess "KAIA-BE/db"
	"KAIA-BE/model"
	res "KAIA-BE/responses"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Slices struct {
	ID     string  `json:"id"`
	Slices string  `json:"slices"`
	Price  float64 `json:"price"`
}

type ProductResponse struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Price    float64  `json:"price"`
	Category string   `json:"category"`
	Desc     string   `json:"desc"`
	LongDesc string   `json:"longDesc"`
	Image    string   `json:"image"`
	Slices   []Slices `json:"sliceOptions"`
}

func Index(c *gin.Context) {

	var products []model.Product

	result := DataAccess.DB.Limit(50).
		Preload("ProductSlices").
		Find(&products)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, res.Fail{
			Message: result.Error.Error(),
		})
		return
	}

	var response []ProductResponse
	for _, product := range products {
		var slcs []Slices
		for _, product_slices := range product.ProductSlices {
			slcs = append(slcs, Slices{
				ID:     product_slices.ID,
				Slices: product_slices.Slice,
				Price:  product_slices.Price,
			})
		}

		response = append(response, ProductResponse{
			ID:       product.ID,
			Name:     product.Name,
			Price:    float64(product.Price),
			Desc:     product.Description,
			LongDesc: product.Longdesc,
			Category: product.Category,
			Image:    product.Image,
			Slices:   slcs,
		})
	}

	c.JSON(http.StatusOK, res.Success{
		Success: true,
		Data:    response,
	})
}
