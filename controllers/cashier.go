package controllers

import (
	"devcode-pos/helpers"
	"devcode-pos/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func CashierList(c *fiber.Ctx) error {
	var cashierModel models.Cashiers
	var meta models.Pagination
	var result []models.Cashiers

	models.DB.WithContext(c.Context()).
		Scopes(helpers.Paginate(c, &meta, cashierModel)).Find(&result)

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data": map[string]interface{}{
			"cashiers": result,
			"meta":     meta,
		},
	})
}

func CashierAdd(c *fiber.Ctx) error {
	var err error

	var input = struct {
		Name string `validate:"required" json:"name"`
		Pass string `validate:"required" json:"passcode"`
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

	createData := &models.Cashiers{
		CashierName:     input.Name,
		CashierPasscode: input.Pass,
	}

	models.DB.WithContext(c.Context()).Create(&createData)

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data":    createData,
	})
}

func CashierDetail(c *fiber.Ctx) error {
	id, _ := strconv.ParseUint(c.Params("id"), 10, 32)
	var cashierModel models.Cashiers

	res := models.DB.WithContext(c.Context()).
		Where(&models.Cashiers{CashierId: uint32(id)}).
		First(&cashierModel)

	if res.Error != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Cashier Not Found",
			"error":   make(map[string]string),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data":    cashierModel,
	})
}

func CashierDelete(c *fiber.Ctx) error {
	id, _ := strconv.ParseUint(c.Params("id"), 10, 32)
	var cashierModel models.Cashiers

	res := models.DB.WithContext(c.Context()).
		Where(&models.Cashiers{CashierId: uint32(id)}).
		Delete(&cashierModel)

	if res.RowsAffected > 0 {
		return c.Status(200).JSON(fiber.Map{
			"success": true,
			"message": "Success",
		})
	} else {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Cashier Not Found",
			"error":   make(map[string]string),
		})
	}
}

func CashierUpdate(c *fiber.Ctx) error {
	var err error
	id, _ := strconv.ParseUint(c.Params("id"), 10, 32)

	var cashierModel models.Cashiers
	type Input struct {
		Name string `json:"name"`
		Pass string `validate:"omitempty,len=6" json:"passcode"`
	}
	var input Input

	res := models.DB.WithContext(c.Context()).
		Where(&models.Cashiers{CashierId: uint32(id)}).
		First(&cashierModel)

	if res.Error != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Cashier Not Found",
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

	if input == (Input{}) {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "test",
			"error": map[string]interface{}{
				"message": "value must be an array",
				"path":    make([]int, 0),
				"type":    "array.base",
				"context": map[string]interface{}{
					"label": "value",
					"value": make(map[string]string),
				},
			},
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

	cashierModel.CashierName = input.Name
	cashierModel.CashierPasscode = input.Pass
	go models.DB.Save(&cashierModel)

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Success",
	})
}
