package controllers

import (
	"devcode-pos/models"

	"github.com/gofiber/fiber/v2"
)

func Revenues(c *fiber.Ctx) error {
	var model models.Revenues

	models.DB.WithContext(c.Context()).
		Table("orders").
		Select("payment_id, payment_name, payment_logo, SUM(order_total_price) AS payment_total_amount").
		Joins("JOIN payments on payment_id = order_payment_id").
		Group("order_payment_id").Find(&model.PaymentTypes)

	for _, v := range model.PaymentTypes {
		model.TotalRevenue += v.PaymentTotalAmount
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data":    model,
	})
}

func Solds(c *fiber.Ctx) error {
	var model models.Solds

	models.DB.WithContext(c.Context()).
		Table("order_details").
		Select("detail_product_id AS product_id, detail_product_name AS product_name, SUM(detail_discount_qty) AS product_total_qty_sold, SUM(detail_discount_final_price) AS product_total_amount").
		Group("detail_product_id").Find(&model.SoldsOrderProduct)

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data":    model,
	})
}
