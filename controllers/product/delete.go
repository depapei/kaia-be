package product

import (
	DataAccess "KAIA-BE/db"
	"KAIA-BE/model"
	res "KAIA-BE/responses"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Delete(c *gin.Context) {
	id := c.Param("product_id")

	var product model.Product

	DA := DataAccess.DB.Begin()

	if err := DA.First(&product, `"id" = ?`, id).Error; err != nil {
		c.JSON(http.StatusNotFound, res.Fail{
			Message: "Data tidak ditemukan atau sudah dihapus",
		})
		return
	}

	if err := DA.Model(&product).Update(`"is_deleted"`, true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, res.Fail{
			Message: "Gagal menghapus produk",
		})
		DA.Rollback()
		return
	}

	DA.Commit()

	c.JSON(http.StatusOK, res.Success{
		Success: true,
		Data:    "Data berhasil di hapus",
	})
}
