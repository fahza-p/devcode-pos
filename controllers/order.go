package controllers

import (
	"devcode-pos/helpers"
	"devcode-pos/models"
	"fmt"
	"math"
	"sort"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func OrderSubTotal(c *fiber.Ctx) error {
	var err error

	type inputProduct struct {
		ProductId uint16 `json:"productId"`
		Qty       int32  `json:"qty"`
	}

	var input []inputProduct
	var productModel []models.Subtotal
	var subtotal float64

	// Parsing Input
	if err = c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Bad Request",
			"error":   make([]int, 0),
		})
	}

	// Run Validation
	if len(input) == 0 {
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

	// Sort ASC
	sort.Slice(input, func(i, j int) bool {
		return input[i].ProductId < input[j].ProductId
	})

	var productId []uint16
	for _, el := range input {
		productId = append(productId, el.ProductId)
	}

	models.DB.
		Table("products").
		Where("product_id IN ?", productId).
		Preload("ProductDiscount").
		Order("product_id asc").
		Find(&productModel)

	for index := range productModel {
		productModel[index].Qty = input[index].Qty
		productModel[index].TotalNormalPrice = productModel[index].ProductPrice * float32(productModel[index].Qty)
		switch productModel[index].ProductDiscount.DiscountType {
		case "BUY_N":
			productModel[index].TotalFinalPrice = (productModel[index].ProductDiscount.DiscountResult * float32(productModel[index].Qty/int32(productModel[index].ProductDiscount.DiscountQty))) + (productModel[index].ProductPrice * float32(productModel[index].Qty%int32(productModel[index].ProductDiscount.DiscountQty)))
		case "PERCENT":
			productModel[index].TotalFinalPrice = (productModel[index].ProductPrice - (productModel[index].ProductPrice * productModel[index].ProductDiscount.DiscountResult / 100)) * float32(productModel[index].Qty)
		default:
			productModel[index].TotalFinalPrice = productModel[index].ProductPrice * float32(productModel[index].Qty)
		}

		subtotal += math.Round(float64(productModel[index].TotalFinalPrice)*100) / 100
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data": map[string]interface{}{
			"products": productModel,
			"subtotal": subtotal,
		},
	})
}

func OrderAdd(c *fiber.Ctx) error {
	var err error
	var subtotalModel []models.Subtotal
	var orderModel models.Orders
	var subtotal float64
	var createDetailData []models.OrderDetails

	// Get & Decode Token
	// getToken := c.Locals("user").(*jwt.Token)
	// user := getToken.Claims.(jwt.MapClaims)

	type inputProduct struct {
		ProductId uint16 `json:"productId"`
		Qty       int32  `json:"qty"`
	}

	var input = struct {
		PaymentId uint32         `validate:"required" json:"paymentId"`
		TotalPaid float32        `validate:"required" json:"totalPaid"`
		Products  []inputProduct `validate:"required" json:"products"`
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

	// Sort ASC
	sort.Slice(input.Products, func(i, j int) bool {
		return input.Products[i].ProductId < input.Products[j].ProductId
	})

	var productId []uint16
	for _, el := range input.Products {
		productId = append(productId, el.ProductId)
	}

	models.DB.
		Table("products").
		Where("product_id IN ?", productId).
		Preload("ProductDiscount").
		Order("product_id asc").
		Find(&subtotalModel)

	models.DB.Last(&orderModel)

	for index := range subtotalModel {
		subtotalModel[index].Qty = input.Products[index].Qty
		subtotalModel[index].TotalNormalPrice = subtotalModel[index].ProductPrice * float32(subtotalModel[index].Qty)
		switch subtotalModel[index].ProductDiscount.DiscountType {
		case "BUY_N":
			subtotalModel[index].TotalFinalPrice = (subtotalModel[index].ProductDiscount.DiscountResult * float32(subtotalModel[index].Qty/int32(subtotalModel[index].ProductDiscount.DiscountQty))) + (subtotalModel[index].ProductPrice * float32(subtotalModel[index].Qty%int32(subtotalModel[index].ProductDiscount.DiscountQty)))
		case "PERCENT":
			subtotalModel[index].TotalFinalPrice = (subtotalModel[index].ProductPrice - (subtotalModel[index].ProductPrice * subtotalModel[index].ProductDiscount.DiscountResult / 100)) * float32(subtotalModel[index].Qty)
		default:
			subtotalModel[index].TotalFinalPrice = subtotalModel[index].ProductPrice * float32(subtotalModel[index].Qty)
		}

		subtotal += math.Round(float64(subtotalModel[index].TotalFinalPrice)*100) / 100

		createDetailData = append(createDetailData, models.OrderDetails{
			DetailProductId:           uint32(subtotalModel[index].ProductId),
			DetailOrderId:             orderModel.OrderId + 1,
			DetailProductName:         subtotalModel[index].ProductName,
			DetailDiscountId:          subtotalModel[index].ProductDiscount.DiscountId,
			DetailDiscountQty:         subtotalModel[index].Qty,
			DetailDiscountPrice:       subtotalModel[index].TotalNormalPrice - subtotalModel[index].TotalFinalPrice,
			DetailDiscountNormalPrice: subtotalModel[index].TotalNormalPrice,
			DetailDiscountFinalPrice:  subtotalModel[index].TotalFinalPrice,
		})
	}
	go models.DB.WithContext(c.Context()).Create(&createDetailData)

	createData := &models.Orders{
		OrderId:          orderModel.OrderId + 1,
		OrderCashiersId:  1,
		OrderPaymentId:   input.PaymentId,
		OrderTotalPrice:  float32(subtotal),
		OrderTotalPaid:   input.TotalPaid,
		OrderTotalReturn: float32(input.TotalPaid) - float32(subtotal),
		OrderRecipeId:    fmt.Sprintf("S1%06d", orderModel.OrderId+1),
	}

	models.DB.WithContext(c.Context()).Create(&createData)

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data": map[string]interface{}{
			"order":    createData,
			"products": subtotalModel,
		},
	})
}

func OrderList(c *fiber.Ctx) error {
	var model models.Orders
	var meta models.Pagination
	var result []models.Orders
	models.DB.
		Scopes(helpers.Paginate(c, &meta, model)).
		Preload("Cashiers").
		Preload("PaymentType").
		Find(&result)

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data": map[string]interface{}{
			"orders": result,
			"meta":   meta,
		},
	})
}

func OrderDetail(c *fiber.Ctx) error {
	id, _ := strconv.ParseUint(c.Params("id"), 10, 32)
	var model models.Orders

	// type ProductId struct {
	// 	DetailProductId uint32 `json:"product_id"`
	// }

	var arrProduct []models.Products

	res := models.DB.WithContext(c.Context()).
		Where(&models.Orders{OrderId: uint32(id)}).
		Preload("Cashiers").
		Preload("PaymentType").
		First(&model)

	if res.Error != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "Order Not Found",
			"error":   make(map[string]string),
		})
	}

	models.DB.WithContext(c.Context()).
		Table("products").
		Joins("JOIN order_details ON detail_product_id = product_id").
		Where("detail_order_id = ?", uint32(id)).
		Preload("ProductDiscount").
		Find(&arrProduct)

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data": map[string]interface{}{
			"order":    model,
			"products": arrProduct,
		},
	})
}
