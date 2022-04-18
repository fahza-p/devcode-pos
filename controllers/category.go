package controllers

import (
	"devcode-pos/helpers"
	"devcode-pos/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func CategoryList(c *fiber.Ctx) error {
	var categoryModel models.Categories
	var meta models.Pagination
	var result []models.Categories

	models.DB.WithContext(c.Context()).
		Scopes(helpers.Paginate(c, &meta, categoryModel)).Find(&result)

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data": map[string]interface{}{
			"categories": result,
			"meta":       meta,
		},
	})
}

func CategoryDetail(c *fiber.Ctx) error {
	id, _ := strconv.ParseUint(c.Params("id"), 10, 32)
	var model models.Categories

	res := models.DB.WithContext(c.Context()).
		Where(&models.Categories{CategoryId: uint32(id)}).
		First(&model)

	if res.Error != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Category Not Found",
			"error":   make(map[string]string),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data":    model,
	})
}

func CategoryAdd(c *fiber.Ctx) error {
	var err error

	var input = struct {
		Name string `validate:"required" json:"name"`
	}{}

	// Parsing Input
	if err = c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Bad Request",
			"error":   make([]int, 0),
		})
	}

	// Run Validation
	if msg, err := helpers.RunValidate(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "test",
			"error":   msg,
		})
	}

	createData := &models.Categories{
		CategoryName: input.Name,
	}

	models.DB.WithContext(c.Context()).Create(&createData)

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data":    createData,
	})
}

func CategoryUpdate(c *fiber.Ctx) error {
	var err error
	id, _ := strconv.ParseUint(c.Params("id"), 10, 32)

	var model models.Categories

	var input = struct {
		Name string `validate:"required" json:"name"`
	}{}

	res := models.DB.WithContext(c.Context()).
		Where(&models.Categories{CategoryId: uint32(id)}).
		First(&model)

	if res.Error != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Category Not Found",
			"error":   make(map[string]string),
		})
	}

	// Parsing Input
	if err = c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Bad Request",
			"error":   make([]int, 0),
		})
	}

	// Run Validation
	if msg, err := helpers.RunValidate(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "test",
			"error":   msg,
		})
	}

	model.CategoryName = input.Name
	go models.DB.Save(&model)

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Success",
	})
}

func CategoryDelete(c *fiber.Ctx) error {
	id, _ := strconv.ParseUint(c.Params("id"), 10, 32)
	var cashierModel models.Categories

	res := models.DB.WithContext(c.Context()).
		Where(&models.Categories{CategoryId: uint32(id)}).
		Delete(&cashierModel)

	if res.RowsAffected > 0 {
		return c.Status(200).JSON(fiber.Map{
			"success": true,
			"message": "Success",
		})
	} else {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Category Not Found",
			"error":   make(map[string]string),
		})
	}
}
