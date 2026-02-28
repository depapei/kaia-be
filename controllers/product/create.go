package product

import (
	DataAccess "KAIA-BE/db"
	"KAIA-BE/model"
	res "KAIA-BE/responses"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SliceOptsInput struct {
	Slices string  `json:"slices" binding:"required"`
	Price  float64 `json:"price" binding:"required"`
}

type ProductInput struct {
	ID           string           `json:"id,omitempty"`
	Name         string           `json:"name" binding:"required"`
	Price        float64          `json:"price" binding:"required"`
	Category     string           `json:"category" binding:"required"`
	Desc         string           `json:"desc" binding:"required"`
	LongDesc     string           `json:"longDesc" binding:"required"`
	Image        string           `json:"image" binding:"required"`
	CreatedBy    string           `json:"createdBy" binding:"required"`
	SliceOptions []SliceOptsInput `json:"sliceOptions" binding:"required"`
}

func Create(c *gin.Context) {

	var input ProductInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, res.Fail{
			Message: err.Error(),
		})
		return
	}

	product := model.Product{
		Name:        input.Name,
		Price:       int32(input.Price),
		Category:    input.Category,
		Description: input.Desc,
		Longdesc:    input.LongDesc,
		Image:       input.Image,
		CreatedBy:   input.CreatedBy,
	}

	DA := DataAccess.DB.Begin()

	if err := DA.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, res.Fail{
			Message: "Gagal menyimpan produk",
		})
		DA.Rollback()
		return
	}

	var slices []model.Productslice
	for _, slice := range input.SliceOptions {
		slices = append(slices, model.Productslice{
			ProductID: product.ID,
			Slice:     slice.Slices,
			Price:     slice.Price,
		})
	}

	if err := DA.Create(&slices).Error; err != nil {
		c.JSON(http.StatusInternalServerError, res.Fail{
			Message: "Gagal menyimpan potongan produk",
		})
		DA.Rollback()
		return
	}

	DA.Commit()
	input.ID = product.ID

	c.JSON(http.StatusOK, res.Success{
		Success: true,
		Data:    input,
	})
}
