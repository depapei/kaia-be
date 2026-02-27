package main

import (
	authentication "KAIA-BE/controllers/auth"
	"KAIA-BE/controllers/product"
	"KAIA-BE/controllers/transaction"
	"KAIA-BE/controllers/wishlist"
	"KAIA-BE/db"
	"KAIA-BE/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"log"
)

func main() {

	// Load Env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed load .env file")
	}

	router := gin.Default()
	router.Use(Cors())

	db.Connect()
	// db.GenerateModels()

	// authentication
	auth := router.Group("/auth")
	{
		auth.POST("/login", authentication.Login)
		auth.POST("/register", authentication.Register)
	}

	// able to access by public
	p := router.Group("/products")
	{
		p.GET("/", product.Index)
	}
	trx := router.Group("/transactions")
	{
		trx.POST("/", transaction.Create)
	}

	// need login
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{

		trxP := protected.Group("/transactions")
		{
			trxP.GET("/", transaction.Index)
		}

		wl := protected.Group("/wishlists")
		{
			wl.GET("/", wishlist.Index)
		}
	}

	router.Run()
}

func Cors() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Timezone", "User", "X-Telegram-Auth-Date", "X-Telegram-Hash", "X-Telegram-Init-Data", "Service-Token", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "Accept", "Origin", "Cache-Control", "X-Requested-With"},
		AllowCredentials: false,
		ExposeHeaders:    []string{"Total-records", "Content-disposition"},
	})
}
