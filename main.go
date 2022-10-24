package main

import (
	"lifedev/reastshop/controllers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	jwtware "github.com/gofiber/jwt/v3"
)

func main() {
	// session
	store := session.New()

	// load template engine
	app := fiber.New()

	// static
	app.Static("/public", "./public")

	// controllers
	// helloController := controllers.InitHelloController(store)
	prodController := controllers.InitProductController()
	tranController := controllers.InitTransactionController()
	authController := controllers.InitAuthController(store)

	prod := app.Group("/products")
	prod.Get("/", prodController.HomeProduct)
	prod.Get("/dashboard", prodController.DashboardProduct)
	prod.Post("/create", prodController.AddPostedProduct)
	prod.Get("/detail/:id", prodController.GetDetailProduct2)
	prod.Put("/editproduct/:id", prodController.EditlPostedProduct)
	prod.Delete("/deleteproduct/:id", prodController.DeleteProduct)

	app.Post("/register", authController.AddPostedRegister)
	app.Get("/logout", authController.Logout)
	app.Get("/profile", authController.Profile)
	app.Get("/users", authController.AllUser)
	// JWT Middleware
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte("mysecretpassword"),
	}))

	tran := app.Group("/transactions")
	// tran.Get("/", func(c *fiber.Ctx) error {
	// 	sess, _ := store.Get(c)
	// 	val := sess.Get("username")
	// 	if val != nil {
	// 		return c.Next()
	// 	}
	// 	return c.Redirect("/login")
	// }, tranController.DashboardTransaction)
	tran.Get("/", tranController.DashboardTransaction)
	tran.Post("/create/:id", tranController.AddPostedTransaction)
	tran.Put("/bayar/:id", tranController.BayarTransaction)
	tran.Delete("/delete/:id", tranController.DeleteTransactionById)

	app.Post("/login", authController.LoginPosted)
	// tran.Post("/create", func(c *fiber.Ctx) error {
	// 	sess, _ := store.Get(c)
	// 	val := sess.Get("username")
	// 	if val != nil {
	// 		return c.Next()
	// 	}

	// 	return c.Redirect("/login")

	// }, tranController.AddPostedTransaction)
	// tran.Get("/delete/:id", func(c *fiber.Ctx) error {
	// 	sess, _ := store.Get(c)
	// 	val := sess.Get("username")
	// 	if val != nil {
	// 		return c.Next()
	// 	}

	// 	return c.Redirect("/login")

	// }, tranController.DeleteTransactionById)

	// app.Use("/profile", func(c *fiber.Ctx) error {
	// 	sess,_ := store.Get(c)
	// 	val := sess.Get("username")
	// 	if val != nil {
	// 		return c.Next()
	// 	}

	// 	return c.Redirect("/login")

	// })
	app.Get("/profile", func(c *fiber.Ctx) error {
		sess, _ := store.Get(c)
		val := sess.Get("username")
		if val != nil {
			return c.Next()
		}

		return c.Redirect("/login")

	}, authController.Profile)

	app.Listen(":3000")
}
