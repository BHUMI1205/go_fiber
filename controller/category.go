
package controller

import (
	"database1/database"
	"database1/model"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type CategorySerializer struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	NAME string `json:"name"`
}

func CreateResponseCategory(categoryModel model.Category) CategorySerializer {
	return CategorySerializer{ID: categoryModel.ID, NAME: categoryModel.Name}
}

// GetCategory godoc
// @Security BearerAuth
// @Summary Get All Category
// @Description Get All Category
// @Tags Category
// @Router /category/ [get]
func GetCategory(c *fiber.Ctx) error {
	category := []model.Category{}

	database.Database.DB.Find(&category)
	responseCategories := []CategorySerializer{}
	for _, category := range category {
		responseCategory := CreateResponseCategory(category)
		responseCategories = append(responseCategories, responseCategory)
	}
	return c.Status(200).JSON(responseCategories)
}

// AddCategory godoc
// @Security BearerAuth
// @Summary Add a Category
// @Description Add a new category
// @Tags Category
// @Accept json
// @Produce json
// @Param category body CategorySerializer true "Category object to add"
// @Router /category/ [post]
func AddCategory(c *fiber.Ctx) error {
	var category model.Category
	validate := validator.New()

	e := validate.Struct(category)
	if e != nil {
		// Validation failed, handle the error
		errors := e.(validator.ValidationErrors)
		return c.Status(500).JSON(fiber.Map{"status": 500, "message": "validation Failed", "error": errors.Error()})
	}
	if err := c.BodyParser(&category); err != nil {
		return c.Status(500).JSON(err.Error())
	}

	database.Database.DB.Create(&category)
	categoryData := CreateResponseCategory(category)
	return c.Status(200).JSON(fiber.Map{"message": "category Added successfully", "data": categoryData})
}

// UpdateCategory godoc
// @Security BearerAuth
// @Summary Update a Category
// @Description Update an existing category by ID
// @Tags Category
// @Accept json
// @Produce json
// @Param id path string true "Category ID to update"
// @Param category body CategorySerializer true "Updated category object"
// @Router /category/{id} [put]
func UpdateCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	var category model.Category

	validate := validator.New()

	e := validate.Struct(category)
	if e != nil {
		// Validation failed, handle the error
		errors := e.(validator.ValidationErrors)
		return c.Status(500).JSON(fiber.Map{"status": 500, "message": "validation Failed", "error": errors.Error()})
	}

	cat := database.Database.DB.First(&category, id)
	if cat.Error != nil {
		return c.Status(500).JSON(cat.Error)
	}

	var UpdateCategory model.Category
	if err := c.BodyParser(&UpdateCategory); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	if UpdateCategory.Name != "" {
		category.Name = UpdateCategory.Name
	}

	database.Database.DB.Save(&category)
	categoryData := CreateResponseCategory(category)
	return c.Status(200).JSON(fiber.Map{"message": "category Updated successfully", "data": categoryData})
}

// DeleteCategory godoc
// @Security BearerAuth
// @Summary Delete a Category
// @Description Delete a category by ID
// @Tags Category
// @Produce json
// @Param id path string true "Category ID to delete"
// @Router /category/{id} [delete]
func DeleteCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	var category model.Category
	cat := database.Database.DB.First(&category, id)

	if cat.Error != nil {
		return c.Status(500).JSON(cat.Error)
	}

	database.Database.DB.Delete(&category, id)
	return c.Status(200).JSON(fiber.Map{"message": "category Deleted successfully", "data": id})
}
