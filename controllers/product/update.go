package product

import (
	DataAccess "KAIA-BE/db"
	"KAIA-BE/model"
	res "KAIA-BE/responses"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Update(c *gin.Context) {
	id := c.Param("product_id")

	var input ProductInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, res.Fail{
			Message: err.Error(),
		})
		return
	}

	var product model.Product

	DA := DataAccess.DB.Begin()

	if err := DA.First(&product, `"id" = ?`, id).Error; err != nil {
		c.JSON(http.StatusNotFound, res.Fail{
			Message: "Data tidak ditemukan atau sudah dihapus",
		})
		return
	}

	product = model.Product{
		ID:          id,
		Name:        input.Name,
		Price:       int32(input.Price),
		Category:    input.Category,
		Description: input.Desc,
		Longdesc:    input.LongDesc,
		Image:       input.Image,
		CreatedBy:   input.CreatedBy,
	}

	if err := DA.Model(&product).Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, res.Fail{
			Message: "Gagal menyimpan produk",
		})
		DA.Rollback()
		return
	}

	var slices []model.Productslice
	if err := DA.Where(`"product_id" = ?`, id).Delete(&slices).Error; err != nil {
		c.JSON(http.StatusInternalServerError, res.Fail{
			Message: "Gagal menghapus product slices",
		})
		DA.Rollback()
		return
	}

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

	c.JSON(http.StatusOK, res.Success{
		Success: true,
		Data:    input,
	})
}
