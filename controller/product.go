package controller

import (
	"database1/database"
	"database1/model"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ProductSerializer struct {
	ID         uint   `json:"id" gorm:"primaryKey"`
	NAME       string `json:"name"`
	IMAGE      string `json:"image"`
	CATEGORYID uint   `json:"categoryId"`
	CATEGORY   string `json:"category"`
}

func CreateResponseProduct(productModel model.Product) ProductSerializer {
	return ProductSerializer{
		ID:         productModel.ID,
		NAME:       productModel.Name,
		IMAGE:      productModel.Image,
		CATEGORYID: productModel.CategoryId,
		CATEGORY:   productModel.Category.Name}
}

// GetProduct godoc
// @Security BearerAuth
// @Summary Get All Product
// @Description Get All Product
// @Tags Product
// @Router /product/ [get]
func GetProduct(c *fiber.Ctx) error {
	// token := c.Locals("token").(string)
	// fmt.Println(token)
	Product := []model.Product{}

	database.Database.DB.Preload("Category").Limit(6).Offset(4).Find(&Product) //offset as skip
	responseCategories := []ProductSerializer{}
	for _, Product := range Product {
		responseProduct := CreateResponseProduct(Product)
		responseCategories = append(responseCategories, responseProduct)
	}
	return c.Status(200).JSON(responseCategories)
}

// AddProduct godoc
// @Security BearerAuth
// @Summary Add a Product
// @Description Add a new product
// @Tags Product
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "Product Name"
// @Param categoryId formData int true "Category ID"
// @Param image formData file true "Product Image"
// @Router /product/ [post]
func AddProduct(c *fiber.Ctx) error {
	var Product model.Product
	validate := validator.New()

	e := validate.Struct(Product)
	if e != nil {
		// Validation failed, handle the error
		errors := e.(validator.ValidationErrors)
		return c.Status(500).JSON(fiber.Map{"status": 500, "message": "validation Failed", "error": errors.Error()})
	}

	file, err := c.FormFile("image")
	if err != nil {
		log.Println("Error in uploading Image : ", err)
		return c.Status(500).JSON(fiber.Map{"status": 500, "message": "Upload Image"})
	}
	uniqueId := uuid.New()

	filename := strings.Replace(uniqueId.String(), "-", "", -1)

	fileExt := filepath.Ext(file.Filename)

	image := fmt.Sprintf("%s%s", filename, fileExt)

	err = c.SaveFile(file, fmt.Sprintf("./image/%s", image))

	if err != nil {
		log.Println("Error in saving Image :", err)
		return c.JSON(fiber.Map{"status": 500, "message": "Server error", "data": nil})
	}

	imageUrl := fmt.Sprintf("image/%s", image)
	Product.Image = imageUrl
	if err := c.BodyParser(&Product); err != nil {
		return c.Status(500).JSON(err.Error())
	}

	database.Database.DB.Create(&Product)
	ProductData := CreateResponseProduct(Product)
	return c.Status(200).JSON(fiber.Map{"message": "Product Added successfully", "data": ProductData})
}

// UpdateProduct godoc
// @Security BearerAuth
// @Summary Update a Product
// @Description Update an existing product by ID
// @Tags Product
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Product ID to update"
// @Param name formData string true "Product Name"
// @Param categoryId formData int true "Category ID"
// @Param image formData file true "Product Image"
// @Router /product/{id} [put]
func UpdateProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	var Product model.Product

	validate := validator.New()

	e := validate.Struct(Product)
	if e != nil {
		// Validation failed, handle the error
		errors := e.(validator.ValidationErrors)
		return c.Status(500).JSON(fiber.Map{"status": 500, "message": "validation Failed", "error": errors.Error()})
	}

	cat := database.Database.DB.First(&Product, id)
	if cat.Error != nil {
		return c.Status(500).JSON(cat.Error)
	}

	var UpdateProduct model.Product
	if err := c.BodyParser(&UpdateProduct); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	if UpdateProduct.Name != "" {
		Product.Name = UpdateProduct.Name
	}

	if UpdateProduct.CategoryId != 0 {
		Product.CategoryId = UpdateProduct.CategoryId
	}

	database.Database.DB.Save(&Product)
	ProductData := CreateResponseProduct(Product)
	return c.Status(200).JSON(fiber.Map{"message": "Product Updated successfully", "data": ProductData})
}

// DeleteProduct godoc
// @Security BearerAuth
// @Summary Delete a Product
// @Description Delete a product by ID
// @Tags Product
// @Produce json
// @Param id path string true "Product ID to delete"
// @Router /product/:{id} [delete]
func DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	var Product model.Product
	cat := database.Database.DB.First(&Product, id)

	if cat.Error != nil {
		return c.Status(500).JSON(cat.Error)
	}

	database.Database.DB.Delete(&Product, id)
	return c.Status(200).JSON(fiber.Map{"message": "Product Deleted successfully", "data": id})
}
