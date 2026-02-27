package authentication

import (
	DataAccess "KAIA-BE/db"
	"KAIA-BE/model"
	res "KAIA-BE/responses"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type ValidateRegister struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Name     string `json:"name" binding:"required"`
}

type CreatedResponse struct {
	Email     string `json:"email"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

func Register(c *gin.Context) {
	var input ValidateRegister

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, res.Fail{
			Message: "Gagal melakukan registrasi akun, silahkan cek kembali data anda",
		})
		return
	}

	var data model.User

	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, res.Fail{
			Message: "Gagal melakukan enkripsi password, silahkan dicoba lagi",
		})
		return
	}

	data = model.User{
		Email:    input.Email,
		Name:     input.Name,
		Password: string(hashed),
	}

	if err := DataAccess.DB.Create(&data).Error; err != nil {
		if err.Error() == "duplicated key not allowed" {
			errMessage := fmt.Sprintf(`Email %s sudah tersedia!`, input.Email)
			c.JSON(http.StatusInternalServerError, res.Fail{
				Message: errMessage,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, res.Fail{
			Message: "Gagal menyimpan data",
		})
		return
	}

	c.JSON(http.StatusCreated, res.Success{
		Success: true,
		Data: CreatedResponse{
			Email:     data.Email,
			Name:      data.Name,
			CreatedAt: time.Time(data.CreatedAt).Format("2 Jan 2006"),
		},
	})
}
