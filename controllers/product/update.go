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
	if err := DA.Model(&slices).Where(`"product_id" = ?`, id).Find(&slices).Error; err != nil {
		c.JSON(http.StatusInternalServerError, res.Fail{
			Message: "Gagal menemukan product slices",
		})
		DA.Rollback()
		return
	}

	if err := DA.Model(&slices).Where(`"product_id" = ?`, id).Update("is_deleted", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, res.Fail{
			Message: "Gagal menghapus product slices",
		})
		DA.Rollback()
		return
	}

	var sliceIDs []string
	for _, s := range slices {
		sliceIDs = append(sliceIDs, s.ID)
	}

	var usedSliceIDs []string
	if err := DA.
		Model(&model.DetailTransaction{}).
		Where(`productslice_id IN ?`, sliceIDs).
		Distinct().
		Pluck("productslice_id", &usedSliceIDs).
		Error; err != nil {

		c.JSON(http.StatusInternalServerError, res.Fail{
			Message: "Gagal mengambil slice yang dipakai transaksi",
		})
		DA.Rollback()
		return
	}

	usedMap := map[string]bool{}
	for _, id := range usedSliceIDs {
		usedMap[id] = true
	}

	var deleteIDs []string
	for _, s := range slices {
		if !usedMap[s.ID] {
			deleteIDs = append(deleteIDs, s.ID)
		}
	}

	if len(deleteIDs) > 0 {
		if err := DA.
			Where(`id IN ?`, deleteIDs).
			Delete(&model.Productslice{}).
			Error; err != nil {

			c.JSON(http.StatusInternalServerError, res.Fail{
				Message: "Gagal menghapus slice",
			})
			DA.Rollback()
			return
		}
	}

	var newSlices []model.Productslice
	for _, slice := range input.SliceOptions {
		newSlices = append(newSlices, model.Productslice{
			ProductID: product.ID,
			Slice:     slice.Slices,
			Price:     slice.Price,
		})
	}

	if err := DA.Create(&newSlices).Error; err != nil {
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
