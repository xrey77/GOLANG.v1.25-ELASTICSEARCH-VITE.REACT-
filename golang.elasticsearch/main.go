package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	_ "golang.elasticsearch/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	dbconfig "golang.elasticsearch/dbconfig"
	"golang.elasticsearch/middleware"
	auth "golang.elasticsearch/middleware/auth"
	prods "golang.elasticsearch/middleware/prods"
	users "golang.elasticsearch/middleware/users"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	} else {
		dbconfig.Connection()
	}
}

// @title BARCLAYS BANK API Management
// @version 1.0

// @tag.name Auth
// @tag.description Authentication and Authorization

// @tag.name User
// @tag.description Users Management

// @tag.name Products
// @tag.description Products Management

// @tag.name MultiFactor Authenticator
// @tag.description Time-Based One-Time Password (TOTP)

// @description REST API Documentation Gin server. \n Reynald Marquez-Gragasin \n rey107@gmail.com
// @host localhost:5000
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and your token.
func main() {

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Static("/assets", "./assets")

	// router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.InstanceName("swagger"),
		ginSwagger.DefaultModelsExpandDepth(-1),
	))

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.POST("/auth/signin", auth.Login)
	router.POST("/auth/signup", auth.Register)
	router.POST("/addproduct", prods.AddProduct)
	router.GET("/products/list/:page", prods.GetProductList)
	router.GET("/products/search/:page/:key", prods.ProductSearch)
	router.GET("/productreport", prods.ProductPDFReport)
	router.GET("/sales/barchart", prods.GetSalesChart)
	router.GET("/sales/piechart", prods.GetLineChart)
	router.POST("/addsalesdata", prods.AddSalesData)

	authGuard := router.Group("/api")
	authGuard.Use(middleware.AuthMiddleware())
	{
		authGuard.GET("/getallusers", users.GetAllUsers)
		authGuard.GET("/getuserbyid/:id", users.GetUserid)
		authGuard.PATCH("/mfa/activate/:id", auth.MfaActivate)
		authGuard.PATCH("/mfa/verifytotp/:id", auth.MfaVerifyotp)
		authGuard.PATCH("/changepassword/:id", users.ChangePassword)
		authGuard.PATCH("/updateprofile/:id", users.UpdateProfile)
		authGuard.PATCH("/uploadpicture/:id", users.UploadPicture)
		authGuard.DELETE("/deleteuserbyid/:id", users.DeleteUserid)
	}

	host := "localhost"
	port := "5000"
	address := fmt.Sprintf("%s:%s", host, port)
	log.Print("Listening to ", address)
	log.Fatal(http.ListenAndServe(":5000", router))
}
