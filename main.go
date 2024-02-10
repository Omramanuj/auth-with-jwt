package main

import (
	"auth-with-jwt/handlers"
	"auth-with-jwt/middleware"
	"auth-with-jwt/models"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func main() {

	r := gin.Default()

	dsn := "host=localhost user=om password=1234 dbname=om port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panic("Failed to connect to database:", err)
	}

	err = db.AutoMigrate(&models.User{},&models.Product{})
	if err != nil {
		log.Panic("Failed to auto migrate database:", err)
	}

	authHandlers := handlers.NewAuthHandlers(db)

	r.POST("/signup", authHandlers.Signup)
	r.POST("/login", authHandlers.Login)

	productRoutes := r.Group("/products")
	productRoutes.Use(middleware.JWTAuthenticationMiddleware()) // Apply JWT middleware
	{
		productRoutes.GET("/", authHandlers.GetAllProducts)
		productRoutes.POST("Add/", authHandlers.CreateProduct)
		productRoutes.GET("Products/:id", authHandlers.GetProduct)
		productRoutes.PUT("Update/:id", authHandlers.UpdateProduct)
		productRoutes.DELETE("Delete/:id", authHandlers.DeleteProduct)
	}
	port := "8080"

	r.Run(":" + port)
}
