package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"

	"lifedev/reastshop/database"
	"lifedev/reastshop/models"
)

type TransactionController struct {
	// declare variables
	Db *gorm.DB
}

type BayarForm struct {
	// declare variables
	Bayar float64 `form:"inputbayar" validate:"required"`
}

func InitTransactionController() *TransactionController {
	db := database.InitDb()
	// gorm
	db.AutoMigrate(&models.Transaction{})

	return &TransactionController{Db: db}
}

// routing
// GET /transactions
func (controller *TransactionController) DashboardTransaction(c *fiber.Ctx) error {
	// // load all products
	// user := c.Locals("user").(*jwt.Token)
	// claims := user.Claims.(jwt.MapClaims)
	// iduser := claims["id"].(float64)
	// var idu int = int(iduser)
	var transactions []models.Transaction
	err := models.ReadTransaction(controller.Db, &transactions)
	if err != nil {
		return c.SendStatus(500) // http 500 internal server error
	}
	return c.JSON(transactions)
}

// POST /products/create
func (controller *TransactionController) AddPostedTransaction(c *fiber.Ctx) error {
	//myform := new(models.Product)
	id := c.Params("id")
	idn, _ := strconv.Atoi(id)
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	iduser := claims["id"].(float64)
	var idu int = int(iduser)

	var product models.Product
	err := models.ReadProductById(controller.Db, &product, idn)
	if err != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	var myform models.Transaction

	if err := c.BodyParser(&myform); err != nil {
		return c.SendStatus(400)
	}
	jumlahQuantity := product.Quantity - myform.Quantity
	if jumlahQuantity <= 0 {
		return c.JSON(fiber.Map{
			"message": "Stok Habis atau stok kecil",
		})
	}
	myform.IdUser = idu
	myform.IdProduck = product.Id
	myform.Name = product.Name
	myform.Image = product.Image
	myform.Price = product.Price
	myform.Total = float32(product.Price) * float32(myform.Quantity)
	myform.Status = "Belum Bayar"
	// save product
	errr := models.CreateTransaction(controller.Db, &myform)
	if errr != nil {
		return c.SendStatus(500)
	}
	// if succeed
	return c.JSON(myform)
}

// POST /products/create
func (controller *TransactionController) BayarTransaction(c *fiber.Ctx) error {
	//myform := new(models.Product)

	id := c.Params("id")
	idn, _ := strconv.Atoi(id)
	var bayar BayarForm
	var transaction models.Transaction
	err := models.ReadTransactionById(controller.Db, &transaction, idn)
	if err != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	var myform models.Transaction

	if err := c.BodyParser(&myform); err != nil {
		return c.SendStatus(400)
	}

	var bayarhasil float64 = float64(bayar.Bayar) - float64(transaction.Total)

	if bayarhasil <= 0 {
		return c.JSON(fiber.Map{
			"message": "Pembayaran tidak berhasil pastikan uanga anda cukup",
		})
	}
	if transaction.Status == "Sudah Bayar" {
		return c.JSON(fiber.Map{
			"message": "Anda Sudah Membayar",
		})
	}
	// if jumlah <= 0 {
	// 	return c.JSON(fiber.Map{
	// 		"message": "Saldo Anda Kurang",
	// 	})
	// }

	// jika berhasil bayar
	var product models.Product
	errProduct := models.ReadProductById(controller.Db, &product, transaction.IdProduck)
	if errProduct != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	jumlah := product.Quantity - transaction.Quantity
	if jumlah <= 0 {
		return c.JSON(fiber.Map{
			"message": "Stok Habis",
		})
	}
	product.Quantity = jumlah

	// save product
	models.UpdateProduct(controller.Db, &product)

	transaction.Status = "Sudah Bayar"

	// save product
	errr := models.UpdateTransaction(controller.Db, &transaction)
	if errr != nil {
		return c.SendStatus(500)
	}
	// if succeed
	return c.JSON(transaction)
}

// / GET /products/deleteproduct/xx
func (controller *TransactionController) DeleteTransactionById(c *fiber.Ctx) error {
	id := c.Params("id")
	idn, _ := strconv.Atoi(id)

	var transactions models.Transaction
	models.DeleteTransactionById(controller.Db, &transactions, idn)
	return c.JSON(fiber.Map{
		"message": "data was deleted",
	})
}
