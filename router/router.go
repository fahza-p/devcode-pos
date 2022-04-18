package router

import (
	"devcode-pos/controllers"
	"devcode-pos/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).SendString("hello world !!")
	})

	// Cashier
	app.Post("/cashiers", controllers.CashierAdd)
	app.Get("/cashiers/:id", controllers.CashierDetail)
	app.Get("/cashiers", controllers.CashierList)
	app.Delete("/cashiers/:id", controllers.CashierDelete)
	app.Put("/cashiers/:id", controllers.CashierUpdate)

	// Passcode
	app.Get("/cashiers/:id/passcode", controllers.GetPasscode)
	app.Post("/cashiers/:id/login", controllers.VerifyLogin)
	app.Post("/cashiers/:id/logout", controllers.VerifyLogout)

	// Category
	app.Get("/categories", middleware.Protected(), controllers.CategoryList)
	app.Get("/categories/:id", middleware.Protected(), controllers.CategoryDetail)
	app.Post("/categories", controllers.CategoryAdd)
	app.Put("/categories/:id", controllers.CategoryUpdate)
	app.Delete("/categories/:id", controllers.CategoryDelete)

	// Payment
	app.Get("/payments", middleware.Protected(), controllers.PaymentList)
	app.Get("/payments/:id", middleware.Protected(), controllers.PaymentDetail)
	app.Post("/payments", controllers.PaymentAdd)
	app.Put("/payments/:id", controllers.PaymentUpdate)
	app.Delete("/payments/:id", controllers.PaymentDelete)

	// Product
	app.Get("/products", middleware.Protected(), controllers.ProductList)
	app.Post("/products", controllers.ProductAdd)
	app.Delete("/products/:id", controllers.ProductDelete)
	app.Put("/products/:id", controllers.ProductUpdate)
	app.Get("/products/:id", middleware.Protected(), controllers.ProductDetail)

	// Order
	app.Post("/orders/subtotal", middleware.Protected(), controllers.OrderSubTotal)
	app.Post("/orders", middleware.Protected(), controllers.OrderAdd)
	app.Get("/orders", middleware.Protected(), controllers.OrderList)
	app.Get("/orders/:id", middleware.Protected(), controllers.OrderDetail)

	// Report
	app.Get("/revenues", middleware.Protected(), controllers.Revenues)
	app.Get("/solds", middleware.Protected(), controllers.Solds)

	// File
	app.Get("/orders/:id/download", middleware.Protected(), controllers.OrderPdf)
	app.Get("/orders/:id/check-download", middleware.Protected(), controllers.OrderPdfCheck)
}
