package routes

import (
	"Backend_Go/internal/controller/deliveries/http"
	"Backend_Go/internal/middleware"
	"Backend_Go/internal/repositories"

	"github.com/gofiber/fiber/v2"
	websocket "github.com/gofiber/websocket/v2"
)

// SetupRoutes รวม route ทั้งระบบ
// Updated for Production Security & Role Separation
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
	chatHandler *http.ChatHandler, // New
	dealerRepo *repositories.DealerRepository,
) {
	// ... (Previous middleware setup) ...

	api := app.Group("/api")

	// ... (Auth routes) ...
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

	api.Post("/cars/:id/contact", middleware.RequireAuth(), carHandler.RecordContact)

	// ==================== USER (Protected) ====================
	// Users can access their own profile and favorites
	users := api.Group("/users", middleware.RequireAuth())

	users.Get("/me", userHandler.GetMe)
	users.Put("/me", userHandler.UpdateMe)

	users.Get("/:id", userHandler.GetUser)

	// Favorites (Aligned with request)
	favorites := api.Group("/favorites", middleware.RequireAuth())
	favorites.Get("/", favoriteHandler.GetMyFavorites)
	favorites.Post("/:car_id", favoriteHandler.AddFavoriteMe)
	favorites.Delete("/:car_id", favoriteHandler.RemoveFavoriteMe)

	// Reviews (User writes review)
	api.Post("/reviews", middleware.RequireAuth(), reviewHandler.CreateReview)

	// ==================== CHAT (Protected) ====================
	chat := api.Group("/chat", middleware.RequireAuth())
	chat.Post("/send", chatHandler.SendMessage)
	chat.Post("/reply/:id", chatHandler.Reply)
	chat.Get("/conversations", chatHandler.GetConversations)
	chat.Get("/conversations/:id/messages", chatHandler.GetMessages)
	chat.Get("/unread-count", chatHandler.GetUnreadCount)

	// WebSocket Upgrade
	app.Get("/ws", websocket.New(chatHandler.WebSocketUpgrade))

	// ==================== DEALER (Protected) ====================
	// Must be a dealer role AND have an approved dealer profile
	dealer := api.Group("/dealer", middleware.RequireRole("dealer"), middleware.RequireActiveDealer(dealerRepo))

	dealer.Get("/me", dealerHandler.GetMyDealer)
	dealer.Get("/cars", dealerHandler.GetMyCars)
	dealer.Get("/leads", dealerHandler.GetMyLeads)

	// Secure Dealer Actions
	api.Post("/cars", middleware.RequireRole("dealer"), middleware.RequireActiveDealer(dealerRepo), carHandler.CreateCar)

	dealerCars := api.Group("/cars", middleware.RequireRole("dealer"), middleware.RequireActiveDealer(dealerRepo))
	dealerCars.Put("/:id", carHandler.UpdateCar)
	dealerCars.Delete("/:id", carHandler.DeleteCar)
	dealerCars.Patch("/:id/status", carHandler.SetStatus)
	dealerCars.Patch("/:id/sold", carHandler.SetSold)
	dealerCars.Patch("/:id/unpublish", carHandler.SetUnpublish)
	dealerCars.Post("/:id/promote", carHandler.PromoteCar)
	dealerCars.Post("/:id/images", carImageHandler.AddImages)

	adminMiddleware := middleware.AdminOnly(adminHandler.Usecase)
	admin := api.Group("/admin", middleware.RequireAuth(), adminMiddleware)

	admin.Get("/users", adminHandler.GetUsers)
	admin.Get("/dealers", adminHandler.GetDealers)
	admin.Get("/reports", adminHandler.GetReports)

	admin.Post("/dealers/:id/approve", adminHandler.ApproveDealer)
	admin.Patch("/dealers/:id/suspend", adminHandler.SuspendDealer)
	admin.Post("/dealers/:id/reject", adminHandler.RejectDealer)

	admin.Get("/cars", adminHandler.GetCars)

	admin.Post("/cars/:id/approve", adminHandler.ApproveCar)
	admin.Post("/cars/:id/reject", adminHandler.RejectCar)
	admin.Post("/cars/:id/hide", adminHandler.HideCar)
	admin.Post("/cars/:id/flag", adminHandler.FlagCar)
	admin.Delete("/cars/:id", adminHandler.DeleteCar)

	admin.Patch("/users/:id/ban", adminHandler.BanUser)
	admin.Patch("/users/:id/unban", adminHandler.UnbanUser)

	api.Delete("/users/:id/favorites/:car_id", middleware.RequireAuth(), favoriteHandler.RemoveFavorite)
}
