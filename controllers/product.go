package controllers

import (
	"devcode-pos/helpers"
	"devcode-pos/models"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func ProductList(c *fiber.Ctx) error {
	var model models.Products
	var meta models.Pagination
	var result []models.Products
	models.DB.
		Scopes(helpers.Paginate(c, &meta, model)).
		Preload("ProductCategory").
		Preload("ProductDiscount").
		Find(&result)

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data": map[string]interface{}{
			"products": result,
			"meta":     meta,
		},
	})
}

func ProductAdd(c *fiber.Ctx) error {
	var err error

	type discountInput struct {
		Qty       uint16  `json:"qty"`
		Type      string  `json:"type"`
		Result    float32 `json:"result"`
		ExpiredAt int64   `json:"expiredAt"`
	}

	var input = struct {
		CategoryId uint32        `validate:"required" json:"categoryId"`
		Name       string        `validate:"required" json:"name"`
		Image      string        `validate:"required" json:"image"`
		Price      float32       `validate:"required" json:"price"`
		Stock      int64         `validate:"required" json:"stock"`
		Discount   discountInput `validate:"omitempty" json:"discount"`
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

	createData := &models.Products{
		ProductCategoryId: input.CategoryId,
		ProductName:       input.Name,
		ProductImage:      input.Image,
		ProductStock:      input.Stock,
		ProductPrice:      input.Price,
	}

	models.DB.WithContext(c.Context()).Create(&createData)

	createData.ProductSku = fmt.Sprintf("%06d", createData.ProductId)

	if input.Discount != (discountInput{}) {
		discountString := ""

		switch input.Discount.Type {
		case "BUY_N":
			discountString = fmt.Sprintf("Buy %d only Rp. %g", input.Discount.Qty, input.Discount.Result)
		case "PERCENT":
			finalPrice := fmt.Sprintf("%.2f", input.Price-(input.Price*input.Discount.Result/100))
			discountString = fmt.Sprintf("Discount %g%s Rp. %s", input.Discount.Result, "%", finalPrice)
		}

		discountData := &models.Discounts{
			DiscountQty:     input.Discount.Qty,
			DiscountType:    input.Discount.Type,
			DiscountResult:  input.Discount.Result,
			ExpiredAt:       time.Unix(input.Discount.ExpiredAt, 0),
			StringFormat:    discountString,
			ExpiredAtFormat: time.Unix(input.Discount.ExpiredAt, 0).Format("01 Jan 2006"),
		}

		models.DB.WithContext(c.Context()).Create(&discountData)

		createData.ProductDiscountId = discountData.DiscountId
	}

	go models.DB.Save(createData)

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data":    createData,
	})
}

func ProductDelete(c *fiber.Ctx) error {
	id, _ := strconv.ParseUint(c.Params("id"), 10, 32)
	var model models.Products

	res := models.DB.WithContext(c.Context()).
		Where(&models.Products{ProductId: uint32(id)}).
		Delete(&model)

	if res.RowsAffected > 0 {
		return c.Status(200).JSON(fiber.Map{
			"success": true,
			"message": "Success",
		})
	} else {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Product Not Found",
			"error":   make(map[string]string),
		})
	}
}

func ProductUpdate(c *fiber.Ctx) error {
	var err error
	id, _ := strconv.ParseUint(c.Params("id"), 10, 32)

	var model models.Products
	type Input struct {
		CategoryId uint32  `json:"categoryId"`
		Name       string  `json:"name"`
		Image      string  `validate:"omitempty,uri" json:"image"`
		Price      float32 `validate:"omitempty,number" json:"price"`
		Stock      int64   `validate:"omitempty,number" json:"stock"`
	}
	var input Input

	res := models.DB.WithContext(c.Context()).
		Where(&models.Products{ProductId: uint32(id)}).
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

	model.ProductCategoryId = input.CategoryId
	model.ProductName = input.Name
	model.ProductImage = input.Image
	model.ProductPrice = input.Price
	model.ProductStock = input.Stock
	go models.DB.Save(&model)

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Success",
	})
}

func ProductDetail(c *fiber.Ctx) error {
	id, _ := strconv.ParseUint(c.Params("id"), 10, 32)
	var model models.Products

	res := models.DB.WithContext(c.Context()).
		Where(&models.Products{ProductId: uint32(id)}).
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
