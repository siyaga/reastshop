package main

import (
	"lifedev/reastshop/controllers"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

func main() {

	// load template engine
	app := fiber.New()

	// static
	app.Static("/public", "./public")

	// controllers
	// helloController := controllers.InitHelloController(store)
	prodController := controllers.InitProductController()
	tranController := controllers.InitTransactionController()
	authController := controllers.InitAuthController()

	app.Post("/login", authController.LoginPosted)
	prod := app.Group("/products")
	prod.Get("/", prodController.HomeProduct)
	prod.Get("/dashboard", prodController.DashboardProduct)
	prod.Post("/create", prodController.AddPostedProduct)
	prod.Get("/detail/:id", prodController.GetDetailProduct2)
	prod.Put("/editproduct/:id", prodController.EditlPostedProduct)
	prod.Delete("/deleteproduct/:id", prodController.DeleteProduct)

	app.Post("/register", authController.AddPostedRegister)
	app.Get("/users", authController.AllUser)
	// JWT Middleware
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte("mysecretpassword"),
	}))
	tran := app.Group("/transactions")
	tran.Get("/", tranController.DashboardTransaction)
	tran.Post("/create/:id", tranController.AddPostedTransaction)
	tran.Put("/bayar/:id", tranController.BayarTransaction)
	tran.Delete("/delete/:id", tranController.DeleteTransactionById)

	app.Listen(":3000")
}
