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
		c.JSON(http.StatusUnauthorized, res.Fail{
			Message: "Failed to generate token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"token":   tokenString,
	})
}

type ValidateAdminInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func LoginAdmin(c *gin.Context) {
	var input ValidateAdminInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, res.Fail{
			Message: "Login gagal, silahkan periksa data anda",
		})
		return
	}

	var user model.Admin
	if err := DataAccess.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
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

	expTime := time.Now().Add(24 * time.Hour)
	claims := &JWTClaim{
		UserID: user.ID,
		Sub:    user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "kaia-login-admin-system-service",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		fmt.Print(err.Error())
		c.JSON(http.StatusUnauthorized, res.Fail{
			Message: "Failed to generate token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"token":   tokenString,
	})
}
