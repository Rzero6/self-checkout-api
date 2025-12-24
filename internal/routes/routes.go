package routes

import (
	"github.com/Rzero6/self-checkout-api/internal/controllers"
	"github.com/Rzero6/self-checkout-api/internal/middlewares"
	"github.com/Rzero6/self-checkout-api/internal/services"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, paymentService *services.PaymentService) {
	api := app.Group("/api")

	api.Post("/login", controllers.Login)

	product := api.Group("/products")
	product.Get("/", controllers.GetAllProducts)
	product.Get("/donations", controllers.GetDonationsProduct)
	product.Get("/search", controllers.GetProductByBarcode)
	product.Get("/random", controllers.GetProductRandom)
	product.Post("/", middlewares.AuthMiddleware, middlewares.AdminOnly, controllers.CreateProduct)

	cart := api.Group("/cart")
	cart.Post("/", controllers.StartCart)
	cart.Get("/", controllers.GetCurrentCart)
	cart.Delete("/:cart_id", controllers.DeleteAllProductFromCart)
	cart.Get("/details", controllers.GetCartDetailsBySessionID)
	cartDetail := cart.Group("/detail")
	cartDetail.Post("/", controllers.AddProductsToCart)
	cartDetail.Patch("/", controllers.UpdateDetailInCart)
	cartDetail.Delete("/:detail_id", controllers.DeleteDetailFromCart)

	paymentController := controllers.NewPaymentController(paymentService)
	payment := api.Group("/transaction")
	payment.Post("/", paymentController.CreateTransaction)
	payment.Get("/:order_id", paymentController.GetTransactionStatus)
	payment.Get("/:order_id/details", controllers.GetTransactionDetails)
	payment.Post("/:order_id/cancel", paymentController.CancelTransaction)
	payment.Post("/notification", paymentController.MidtransNotification)

}
