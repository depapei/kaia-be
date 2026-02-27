package middleware

import (
	res "KAIA-BE/responses"
	"KAIA-BE/utils"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}

		auth_header := c.GetHeader("Authorization")
		if auth_header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, res.Fail{
				Message: "Please login first!",
			})
			return
		}

		parts := strings.Split(auth_header, "")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusBadRequest, res.Fail{
				Message: "Invalid token format",
			})
			return
		}

		claims, err := utils.ParseJWT(parts[1])
		if err != nil {
			fmt.Println(utils.ParseJWT(parts[1]))
			c.AbortWithStatusJSON(http.StatusUnauthorized, res.Fail{
				Message: "Invalid token!",
			})
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.UserEmail)

		c.Next()
	}
}
