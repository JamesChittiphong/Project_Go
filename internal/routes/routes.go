package routes

import (
	"Backend_Go/internal/controller/deliveries/http"
	"Backend_Go/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes รวม route ทั้งระบบ
// Updated for Production Security & Role Separation
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
	// ==================== GLOBAL MIDDLEWARE ====================
	// Logger, recoverable, cors likely setup in main.go

	api := app.Group("/api")

	// ==================== PUBLIC ====================
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/register-dealer", authHandler.RegisterDealer)
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.Refresh)
	auth.Post("/logout", authHandler.Logout)

	// Public Resources (Read Only)
	api.Get("/cars", carHandler.GetCars)
	api.Get("/cars/:id", carHandler.GetCarDetail)
	api.Get("/cars/:id/images", carImageHandler.GetImages)

	api.Get("/dealers", dealerHandler.GetDealers)
	api.Get("/dealers/:id", dealerHandler.GetDealer)
	api.Get("/dealers/:id/stats", dealerHandler.GetDealerStats)
	api.Get("/dealers/:id/reviews", reviewHandler.GetReviewsByDealer)

	api.Post("/cars/:id/contact", carHandler.RecordContact)

	// ==================== USER (Protected) ====================
	// Users can access their own profile and favorites
	users := api.Group("/users", middleware.RequireAuth())

	users.Get("/me", userHandler.GetMe)
	users.Put("/me", userHandler.UpdateMe)

	users.Get("/:id", userHandler.GetUser)

	// Favorites (Secure Implementation)
	favorites := api.Group("/favorites", middleware.RequireAuth())
	favorites.Get("/me", favoriteHandler.GetMyFavorites)
	favorites.Post("/:car_id", favoriteHandler.AddFavoriteMe)
	favorites.Delete("/:car_id", favoriteHandler.RemoveFavoriteMe)

	// Reviews (User writes review)
	api.Post("/reviews", middleware.RequireAuth(), reviewHandler.CreateReview)
	// api.Delete("/reviews/:id", middleware.RequireAuth(), reviewHandler.DeleteReview) // Owner Check needed

	// ==================== DEALER (Protected) ====================
	dealer := api.Group("/dealer", middleware.RequireRole("dealer"))

	dealer.Get("/me", dealerHandler.GetMyDealer)
	dealer.Get("/cars", dealerHandler.GetMyCars)   // /dealer/cars -> My cars
	dealer.Get("/leads", dealerHandler.GetMyLeads) // /dealer/leads -> My leads (To be implemented)

	// Secure Dealer Actions
	api.Post("/cars", middleware.RequireRole("dealer"), carHandler.CreateCar)

	dealerCars := api.Group("/cars", middleware.RequireRole("dealer"))
	dealerCars.Put("/:id", carHandler.UpdateCar)    // Handler should check ownership
	dealerCars.Delete("/:id", carHandler.DeleteCar) // Handler checks ownership
	dealerCars.Patch("/:id/status", carHandler.SetStatus)
	dealerCars.Patch("/:id/sold", carHandler.SetSold)           // New
	dealerCars.Patch("/:id/unpublish", carHandler.SetUnpublish) // New
	dealerCars.Post("/:id/promote", carHandler.PromoteCar)
	dealerCars.Post("/:id/images", carImageHandler.AddImages)

	adminMiddleware := middleware.AdminOnly(adminHandler.Usecase)
	admin := api.Group("/admin", middleware.RequireAuth(), adminMiddleware)

	admin.Get("/users", adminHandler.GetUsers)
	admin.Get("/dealers", adminHandler.GetDealers)
	admin.Get("/reports", adminHandler.GetReports)

	admin.Post("/dealers/:id/approve", adminHandler.ApproveDealer)
	admin.Patch("/dealers/:id/suspend", func(c *fiber.Ctx) error { return c.SendStatus(501) }) // Placeholder

	admin.Get("/cars", func(c *fiber.Ctx) error { return carHandler.GetCars(c) })     // Reuse
	admin.Get("/cars/flagged", func(c *fiber.Ctx) error { return c.SendStatus(501) }) // Placeholder

	admin.Post("/cars/:id/hide", adminHandler.HideCar)
	admin.Post("/cars/:id/flag", adminHandler.FlagCar)
	admin.Delete("/cars/:id", adminHandler.DeleteCar)

	admin.Patch("/users/:id/ban", func(c *fiber.Ctx) error { return c.SendStatus(501) })
	admin.Patch("/users/:id/unban", func(c *fiber.Ctx) error { return c.SendStatus(501) })

	api.Delete("/users/:id/favorites/:car_id", middleware.RequireAuth(), favoriteHandler.RemoveFavorite) // Verify ID matches token?
}
