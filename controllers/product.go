package controllers

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"lifedev/reastshop/database"
	"lifedev/reastshop/models"
)

// type ProductForm struct {
// 	Email string `form:"email" validate:"required"`
// 	Address string `form:"address" validate:"required"`
// }

type ProductController struct {
	// declare variables
	Db *gorm.DB
}

func InitProductController() *ProductController {
	db := database.InitDb()
	// gorm
	db.AutoMigrate(&models.Product{})

	return &ProductController{Db: db}
}

// routing
// GET /products
func (controller *ProductController) HomeProduct(c *fiber.Ctx) error {
	// load all products
	var products []models.Product
	err := models.ReadProducts(controller.Db, &products)
	if err != nil {
		return c.SendStatus(500) // http 500 internal server error
	}
	return c.JSON(products)
}

// routing
// GET /products
func (controller *ProductController) DashboardProduct(c *fiber.Ctx) error {
	// load all products
	var products []models.Product
	err := models.ReadProducts(controller.Db, &products)
	if err != nil {
		return c.SendStatus(500) // http 500 internal server error
	}
	return c.JSON(products)
}

// POST /products/create
func (controller *ProductController) AddPostedProduct(c *fiber.Ctx) error {
	//myform := new(models.Product)
	var myform models.Product

	file, errFile := c.FormFile("image")
	if errFile != nil {
		fmt.Println("Error File =", errFile)
	}
	var filename string = file.Filename
	if file != nil {

		errSaveFile := c.SaveFile(file, fmt.Sprintf("./public/images/%s", filename))
		if errSaveFile != nil {
			fmt.Println("Fail to store file into public/ikmages directory.")
		}
	} else {
		fmt.Println("Nothing file to uploading.")
	}

	if err := c.BodyParser(&myform); err != nil {
		return c.Redirect("/products")
	}

	myform.Image = filename
	// save product
	errr := models.CreateProduct(controller.Db, &myform)
	if errr != nil {
		return c.Redirect("/products")
	}
	// if succeed
	return c.JSON(myform)
}

// GET /products/detail/xxx
func (controller *ProductController) GetDetailProduct2(c *fiber.Ctx) error {
	id := c.Params("id")
	idn, _ := strconv.Atoi(id)

	var product models.Product
	err := models.ReadProductById(controller.Db, &product, idn)
	if err != nil {
		return c.SendStatus(500) // http 500 internal server error
	}
	return c.JSON(product)
}

// / POST products/editproduct/xx
func (controller *ProductController) EditlPostedProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	idn, _ := strconv.Atoi(id)

	var product models.Product
	err := models.ReadProductById(controller.Db, &product, idn)
	if err != nil {
		return c.SendStatus(500) // http 500 internal server error
	}
	var myform models.Product

	if err := c.BodyParser(&myform); err != nil {
		return c.SendStatus(400)
	}

	file, errFile := c.FormFile("image")
	if errFile != nil {
		fmt.Println("Error File =", errFile)
	}
	var filename string = file.Filename
	if file != nil {

		errSaveFile := c.SaveFile(file, fmt.Sprintf("./public/images/%s", filename))
		if errSaveFile != nil {
			fmt.Println("Fail to store file into public/ikmages directory.")
		}
	} else {
		fmt.Println("Nothing file to uploading.")
	}
	myform.Image = filename
	product.Name = myform.Name
	product.Image = myform.Image
	product.Description = myform.Description
	product.Quantity = myform.Quantity
	product.Price = myform.Price
	// save product
	models.UpdateProduct(controller.Db, &product)

	return c.JSON(product)

}

// / GET /products/deleteproduct/xx
func (controller *ProductController) DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	idn, _ := strconv.Atoi(id)

	var product models.Product
	models.DeleteProductById(controller.Db, &product, idn)
	return c.JSON(fiber.Map{
		"message": "Product deleted",
	})
}
