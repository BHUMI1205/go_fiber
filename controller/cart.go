package controller

import (
	"database1/database"
	"database1/model"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"os"
)

type CartSerializer struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	PRODUCTID uint   `json:"productId"`
	PRODUCT   string `json:"product"`
}

func CreateResponseCart(cartModel model.Cart) CartSerializer {
	return CartSerializer{
		ID:        cartModel.ID,
		PRODUCTID: cartModel.ProductId,
		PRODUCT:   cartModel.Product.Name,
	}
}

func AddProductToCart(c *fiber.Ctx) error {
	var Cart model.Cart

	if err := c.BodyParser(&Cart); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}

	tokenStr := c.Locals("token").(string)

	token, err := jwt.ParseWithClaims(tokenStr, &model.Token{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid token",
		})
	}

	if claims, ok := token.Claims.(*model.Token); ok && token.Valid {
		Cart.UserId = claims.Id
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid token claims",
		})
	}

	var product model.Product
	if err := database.Database.DB.First(&product, Cart.ProductId).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Product not found",
		})
	}

	if err := database.Database.DB.Create(&Cart).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to add product to cart",
		})
	}

	CartData := CreateResponseCart(Cart)
	return c.Status(200).JSON(fiber.Map{"message": "Product Added successfully to cart", "data": CartData})
}
