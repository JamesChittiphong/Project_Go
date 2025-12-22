package routes

import (
	"Backend_Go/internal/controller/deliveries/http"
	"Backend_Go/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes รวม route ทั้งระบบ
func SetupRoutes(app *fiber.App,
	carHandler *http.CarHandler,
	leadHandler *http.LeadHandler,
	favoriteHandler *http.FavoriteHandler,
	reviewHandler *http.ReviewHandler,
	userHandler *http.UserHandler,
	dealerHandler *http.DealerHandler,
	carImageHandler *http.CarImageHandler,
	adminHandler *http.AdminHandler,
	authHandler *http.AuthHandler,
) {
	// ==================== PUBLIC ====================
	api := app.Group("/api")

	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.Refresh)
	auth.Post("/logout", authHandler.Logout)

	auth.Post("/register/dealer", authHandler.RegisterDealer)
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.Refresh)
	auth.Post("/logout", authHandler.Logout)

	// ---------- Cars ----------
	api.Get("/cars", carHandler.GetCars)
	api.Get("/cars/:id", carHandler.GetCarDetail)
	api.Put("/cars/:id", middleware.RequireRole("dealer"), carHandler.UpdateCar)
	api.Patch("/cars/:id/status", middleware.RequireRole("dealer"), carHandler.SetStatus)
	api.Post("/cars/:id/contact", carHandler.RecordContact)
	api.Post("/cars/:id/promote", middleware.RequireRole("dealer"), carHandler.PromoteCar)
	api.Get("/cars/:id/stats", middleware.RequireRole("dealer"), carHandler.GetStats)

	// ---------- Dealers ----------
	api.Get("/dealers", dealerHandler.GetDealers)
	api.Get("/dealers/:id", dealerHandler.GetDealer)

	// ---------- Reviews ----------
	api.Get("/dealers/:id/reviews", reviewHandler.GetReviewsByDealer)

	// ==================== USER ====================
	user := api.Group("/users")

	user.Post("/", userHandler.CreateUser)
	user.Get("/:id", userHandler.GetUser)
	user.Put("/:id", userHandler.UpdateUser)
	user.Delete("/:id", userHandler.DeleteUser)

	// ---------- Favorites ----------
	user.Get("/:id/favorites", favoriteHandler.GetFavoritesByUser)
	api.Post("/favorites", favoriteHandler.AddFavorite)

	// ==================== DEALER ====================
	dealer := api.Group("/dealers", middleware.RequireRole("dealer"))

	dealer.Post("/", dealerHandler.CreateDealer)
	dealer.Put("/:id", dealerHandler.UpdateDealer)

	// ---------- Cars (Dealer) ----------
	api.Post("/cars", middleware.RequireRole("dealer"), carHandler.CreateCar)
	api.Delete("/cars/:id", middleware.RequireRole("dealer"), carHandler.DeleteCar)

	// ---------- Leads ----------
	api.Post("/leads", leadHandler.CreateLead)
	api.Get("/dealers/:id/leads", middleware.RequireRole("dealer"), leadHandler.GetLeadsByDealer)

	// ---------- Car Images ----------
	api.Post("/cars/:id/images", middleware.RequireRole("dealer"), carImageHandler.AddImage)
	api.Get("/cars/:id/images", carImageHandler.GetImages)
	api.Delete("/images/:id", carImageHandler.DeleteImage)

	// ---------- Reviews ----------
	api.Post("/reviews", reviewHandler.CreateReview)

	// ==================== ADMIN ====================
	admin := api.Group("/admin", middleware.RequireRole("admin"))

	admin.Get("/users", adminHandler.GetUsers)
	admin.Get("/dealers", adminHandler.GetDealers)
	admin.Get("/reports", adminHandler.GetReports)
	admin.Post("/dealers/:id/approve", adminHandler.ApproveDealer)
	admin.Post("/cars/:id/hide", adminHandler.HideCar)
	admin.Post("/cars/:id/flag", adminHandler.FlagCar)
	admin.Delete("/cars/:id", adminHandler.DeleteCar)
}
