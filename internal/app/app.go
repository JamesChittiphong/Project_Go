package app

import (
	"Backend_Go/internal/controller/deliveries/http"
	"Backend_Go/internal/repositories"
	"Backend_Go/internal/routes"
	"Backend_Go/internal/ws"

	adminUC "Backend_Go/internal/usecases/admin"
	authUC "Backend_Go/internal/usecases/auth"
	carUC "Backend_Go/internal/usecases/car"
	carImageUC "Backend_Go/internal/usecases/car_image"
	"Backend_Go/internal/usecases/chat"
	dealerUC "Backend_Go/internal/usecases/dealer"
	favoriteUC "Backend_Go/internal/usecases/favorite"
	lendUC "Backend_Go/internal/usecases/lend"
	reviewUC "Backend_Go/internal/usecases/review"
	userUC "Backend_Go/internal/usecases/user"
	_ "Backend_Go/internal/ws"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/gorm"
)

func NewApp(db *gorm.DB) *fiber.App {
	app := fiber.New()

	// =====================================================
	// ✅ STATIC FILES (ต้องอยู่บนสุด ไม่โดน middleware)
	// =====================================================

	// =====================================================
	// CORS
	// =====================================================
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET, POST, HEAD, PUT, DELETE, PATCH, OPTIONS",
		AllowHeaders:     "*",
		ExposeHeaders:    "Content-Length",
		AllowCredentials: false,
		MaxAge:           86400,
	}))

	// =====================================================
	// REPOSITORIES
	// =====================================================
	carRepo := &repositories.CarRepository{DB: db}
	carImageRepo := &repositories.CarImageRepository{DB: db}
	dealerRepo := &repositories.DealerRepository{DB: db}
	leadRepo := &repositories.LeadRepository{DB: db}
	favoriteRepo := &repositories.FavoriteRepository{DB: db}
	reviewRepo := &repositories.ReviewRepository{DB: db}
	reportRepo := &repositories.ReportRepository{DB: db}

	userRepo := repositories.NewUserRepository(db)
	refreshTokenRepo := repositories.NewRefreshTokenRepository(db)

	// =====================================================
	// USECASES
	// =====================================================
	authUsecase := authUC.NewAuthUsecase(
		userRepo,
		dealerRepo,
		refreshTokenRepo,
	)

	carUsecase := &carUC.CarUsecase{
		CarRepo:      carRepo,
		DealerRepo:   dealerRepo,
		LeadRepo:     leadRepo,
		FavoriteRepo: favoriteRepo,
	}

	carImageUsecase := &carImageUC.CarImageUsecase{
		CarImageRepo: carImageRepo,
		CarRepo:      carRepo,
	}

	leadUsecase := &lendUC.LeadUsecase{
		LeadRepo:   leadRepo,
		CarRepo:    carRepo,
		DealerRepo: dealerRepo,
	}

	favoriteUsecase := &favoriteUC.FavoriteUsecase{
		FavoriteRepo: favoriteRepo,
		CarRepo:      carRepo,
	}

	reviewUsecase := &reviewUC.ReviewUsecase{
		ReviewRepo: reviewRepo,
		DealerRepo: dealerRepo,
	}

	adminUsecase := &adminUC.AdminUsecase{
		UserRepo:   userRepo,
		DealerRepo: dealerRepo,
		ReportRepo: reportRepo,
		CarRepo:    carRepo,
	}

	dealerUsecase := &dealerUC.DealerUsecase{
		DealerRepo: dealerRepo,
		CarRepo:    carRepo,
		ReviewRepo: reviewRepo,
	}

	userUsecase := &userUC.UserUsecase{
		UserRepo: userRepo,
	}

	// =====================================================
	// HANDLERS
	// =====================================================
	carHandler := &http.CarHandler{Usecase: carUsecase}
	carImageHandler := &http.CarImageHandler{Usecase: carImageUsecase}
	leadHandler := &http.LeadHandler{Usecase: leadUsecase}
	favoriteHandler := &http.FavoriteHandler{Usecase: favoriteUsecase}
	reviewHandler := &http.ReviewHandler{Usecase: reviewUsecase}
	userHandler := &http.UserHandler{Usecase: userUsecase}
	dealerHandler := &http.DealerHandler{Usecase: dealerUsecase}
	adminHandler := &http.AdminHandler{Usecase: adminUsecase}
	authHandler := &http.AuthHandler{Usecase: authUsecase}

	// =====================================================
	// ROUTES
	// =====================================================
	// Chat Initialization
	chatHub := ws.NewHub()
	go chatHub.Run()

	chatRepo := &repositories.ChatRepository{DB: db}
	chatUsecase := &chat.ChatUsecase{
		ChatRepo:   chatRepo,
		DealerRepo: dealerRepo,
		Hub:        chatHub,
	}
	chatHandler := &http.ChatHandler{
		Usecase:    chatUsecase,
		DealerRepo: dealerRepo,
		Hub:        chatHub,
	}

	// Setup Routes
	routes.SetupRoutes(
		app,
		carHandler,
		leadHandler,
		favoriteHandler,
		reviewHandler,
		userHandler,
		dealerHandler,
		carImageHandler,
		adminHandler,
		authHandler,
		chatHandler,
		dealerRepo,
	)

	return app
}
