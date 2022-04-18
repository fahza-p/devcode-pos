package controllers

import (
	"devcode-pos/helpers"
	"devcode-pos/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func PaymentList(c *fiber.Ctx) error {
	var model models.Payments
	var meta models.Pagination
	var result []models.Payments

	models.DB.WithContext(c.Context()).
		Scopes(helpers.Paginate(c, &meta, model)).Find(&result)

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data": map[string]interface{}{
			"payments": result,
			"meta":     meta,
		},
	})
}

func PaymentDetail(c *fiber.Ctx) error {
	id, _ := strconv.ParseUint(c.Params("id"), 10, 32)
	var model models.Payments

	res := models.DB.WithContext(c.Context()).
		Where(&models.Payments{PaymentId: uint32(id)}).
		First(&model)

	if res.Error != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Payment Not Found",
			"error":   make(map[string]string),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data":    model,
	})
}

func PaymentAdd(c *fiber.Ctx) error {
	var err error

	var input = struct {
		Name string `validate:"required" json:"name"`
		Type string `validate:"required" json:"type"`
		Logo string `json:"logo"`
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

	createData := &models.Payments{
		PaymentName: input.Name,
		PaymentType: input.Type,
		PaymentLogo: input.Logo,
	}

	models.DB.WithContext(c.Context()).Create(&createData)

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data":    createData,
	})
}

func PaymentUpdate(c *fiber.Ctx) error {
	var err error
	id, _ := strconv.ParseUint(c.Params("id"), 10, 32)

	var model models.Payments
	type Input struct {
		Name string `json:"name"`
		Type string `validate:"omitempty,enum=CASH;E-WALLET;EDC" json:"type"`
		Logo string `validate:"omitempty,uri" json:"logo"`
	}
	var input Input

	res := models.DB.WithContext(c.Context()).
		Where(&models.Payments{PaymentId: uint32(id)}).
		First(&model)

	if res.Error != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Payment Not Found",
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

	model.PaymentName = input.Name
	model.PaymentType = input.Type
	model.PaymentLogo = input.Logo
	go models.DB.Save(&model)

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Success",
	})
}

func PaymentDelete(c *fiber.Ctx) error {
	id, _ := strconv.ParseUint(c.Params("id"), 10, 32)
	var model models.Payments

	res := models.DB.WithContext(c.Context()).
		Where(&models.Payments{PaymentId: uint32(id)}).
		Delete(&model)

	if res.RowsAffected > 0 {
		return c.Status(200).JSON(fiber.Map{
			"success": true,
			"message": "Success",
		})
	} else {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Payment Not Found",
			"error":   make(map[string]string),
		})
	}
}
