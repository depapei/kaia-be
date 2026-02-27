package authentication

import (
	DataAccess "KAIA-BE/db"
	"KAIA-BE/model"
	res "KAIA-BE/responses"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type ValidateInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type DataResponse struct {
	Email string `json:"email"`
	ID    string `json:"user_id"`
	Token string `json:"token"`
}

type JWTClaim struct {
	UserID    string `json:"user_id"`
	UserEmail string `json:"user_email"`
	Sub       string `json:"sub"`
	jwt.RegisteredClaims
}

var jwtSecret = []byte(os.Getenv("SECRET_KEY"))

func Login(c *gin.Context) {
	var input ValidateInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, res.Fail{
			Message: "Login gagal, silahkan periksa data anda",
		})
		return
	}

	var user model.User
	if err := DataAccess.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, res.Fail{
			Message: "User tidak ditermukan",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, res.Fail{
			Message: "Password salah!",
		})
		return
	}

	expTime := time.Now().Add(365 * time.Hour)
	claims := &JWTClaim{
		UserID:    user.ID,
		UserEmail: user.Email,
		Sub:       (user.Email) + "-" + (user.ID),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "kaia-login-system-service",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		fmt.Print(err.Error())
		fmt.Print(jwtSecret)
		c.JSON(http.StatusUnauthorized, res.Fail{
			Message: "Failed to generate token",
		})
		return
	}

	c.JSON(http.StatusOK, res.Success{
		Success: true,
		Data: DataResponse{
			ID:    user.ID,
			Email: user.Email,
			Token: tokenString,
		},
	})

}
