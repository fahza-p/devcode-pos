package router

import (
	"devcode-pos/controllers"

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
	app.Get("/categories", controllers.CategoryList)
	app.Get("/categories/:id", controllers.CategoryDetail)
	app.Post("/categories", controllers.CategoryAdd)
	app.Put("/categories/:id", controllers.CategoryUpdate)
	app.Delete("/categories/:id", controllers.CategoryDelete)

	// Payment
	app.Get("/payments", controllers.PaymentList)
	app.Get("/payments/:id", controllers.PaymentDetail)
	app.Post("/payments", controllers.PaymentAdd)
	app.Put("/payments/:id", controllers.PaymentUpdate)
	app.Delete("/payments/:id", controllers.PaymentDelete)

	// Product
	app.Get("/products", controllers.ProductList)
	app.Post("/products", controllers.ProductAdd)
	app.Delete("/products/:id", controllers.ProductDelete)
	app.Put("/products/:id", controllers.ProductUpdate)
	app.Get("/products/:id", controllers.ProductDetail)

	// Order
	app.Post("/orders/subtotal", controllers.OrderSubTotal)
	app.Post("/orders", controllers.OrderAdd)
	app.Get("/orders", controllers.OrderList)
	app.Get("/orders/:id", controllers.OrderDetail)

	// Report
	app.Get("/revenues", controllers.Revenues)
	app.Get("/solds", controllers.Solds)
}
