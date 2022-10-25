package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"lifedev/reastshop/database"
	"lifedev/reastshop/models"
)

// type ProductForm struct {
// 	Email string `form:"email" validate:"required"`
// 	Address string `form:"address" validate:"required"`
// }

type LoginController struct {
	// declare variables
	Db *gorm.DB
}

type LoginForm struct {
	// declare variables
	Username string `form:"username" validate:"required"`
	Password string `form:"password" validate:"required"`
}

func InitAuthController() *LoginController {
	db := database.InitDb()
	// gorm
	db.AutoMigrate(&models.User{})
	return &LoginController{Db: db}
}

// POST /register
func (controller *LoginController) AddPostedRegister(c *fiber.Ctx) error {
	// load all products
	var myform models.User
	var convertpass LoginForm

	if err := c.BodyParser(&myform); err != nil {
		return c.SendStatus(400)
	}
	comvertpassword, _ := bcrypt.GenerateFromPassword([]byte(convertpass.Password), 10)
	sHash := string(comvertpassword)

	myform.Password = sHash

	// save product
	err := models.CreateUser(controller.Db, &myform)
	if err != nil {
		return c.SendStatus(500)
	}
	// if succeed
	return c.JSON(myform)
}

// POST /login
func (controller *LoginController) LoginPosted(c *fiber.Ctx) error {
	// load all products

	var user models.User
	var myform LoginForm
	if err := c.BodyParser(&myform); err != nil {
		return c.SendStatus(400)
	}

	er := models.FindByUsername(controller.Db, &user, myform.Username)
	if er != nil {
		return c.SendStatus(500) // http 500 internal server error
	}

	// hardcode auth
	mycompare := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(myform.Password))
	if mycompare != nil {
		exp := time.Now().Add(time.Hour * 72)
		claims := jwt.MapClaims{
			"id":    user.Id,
			"name":  user.Name,
			"admin": true,
			"exp":   exp.Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte("mysecretpassword"))
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.JSON(fiber.Map{
			"token":   t,
			"expired": exp.Format("2006-01-02 15:04:05"),
		})
	}
	return c.SendStatus(fiber.StatusUnauthorized)

}

// // /profile

// func (controller *LoginController) Profile(c *fiber.Ctx) error {
// 	sess, err := controller.store.Get(c)
// 	if err != nil {
// 		panic(err)
// 	}
// 	val := sess.Get("username")

// 	return c.JSON(fiber.Map{
// 		"username": val,
// 	})
// }

func (controller *LoginController) AllUser(c *fiber.Ctx) error {
	var users []models.User
	err := models.ReadUser(controller.Db, &users)
	if err != nil {
		return c.SendStatus(500) // http 500 internal server error
	}
	return c.JSON(users)
}

// // /logout
// func (controller *LoginController) Logout(c *fiber.Ctx) error {
// 	sess, err := controller.store.Get(c)
// 	if err != nil {
// 		panic(err)
// 	}
// 	sess.Destroy()

// 	return c.JSON(fiber.Map{
// 		"message": "User Logout",
// 	})
// }
