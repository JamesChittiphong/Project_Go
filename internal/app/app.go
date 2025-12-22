package app

import (
	"Backend_Go/internal/controller/deliveries/http"
	"Backend_Go/internal/repositroies"
	"Backend_Go/internal/routes"
	"Backend_Go/internal/usecases"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func NewApp(db *gorm.DB) *fiber.App {
	app := fiber.New()

	// ========== REPOSITORIES ==========
	carRepo := &repositroies.CarRepository{DB: db}
	carImageRepo := &repositroies.CarImageRepository{DB: db}
	dealerRepo := &repositroies.DealerRepository{DB: db}
	leadRepo := &repositroies.LeadRepository{DB: db}
	favoriteRepo := &repositroies.FavoriteRepository{DB: db}
	reviewRepo := &repositroies.ReviewRepository{DB: db}
	userRepo := repositroies.NewUserRepository(db)
	reportRepo := &repositroies.ReportRepository{DB: db}
	refreshTokenRepo := repositroies.NewRefreshTokenRepository(db)

	// ========== USECASES ==========
	authUsecase := usecases.NewAuthUsecase(userRepo, dealerRepo, refreshTokenRepo)

	carUsecase := &usecases.CarUsecase{
		CarRepo:      carRepo,
		ImageRepo:    carImageRepo,
		DealerRepo:   dealerRepo,
		LeadRepo:     leadRepo,
		FavoriteRepo: favoriteRepo,
	}

	leadUsecase := &usecases.LeadUsecase{
		LeadRepo:   leadRepo,
		CarRepo:    carRepo,
		DealerRepo: dealerRepo,
	}

	favoriteUsecase := &usecases.FavoriteUsecase{
		FavoriteRepo: favoriteRepo,
		CarRepo:      carRepo,
	}

	reviewUsecase := &usecases.ReviewUsecase{
		ReviewRepo: reviewRepo,
		DealerRepo: dealerRepo,
	}

	adminUsecase := &usecases.AdminUsecase{
		UserRepo:   userRepo,
		DealerRepo: dealerRepo,
		ReportRepo: reportRepo,
		CarRepo:    carRepo,
	}

	// ========== HANDLERS ==========
	carHandler := &http.CarHandler{Usecase: carUsecase}
	leadHandler := &http.LeadHandler{Usecase: leadUsecase}
	favoriteHandler := &http.FavoriteHandler{Usecase: favoriteUsecase}
	reviewHandler := &http.ReviewHandler{Usecase: reviewUsecase}

	userUsecase := &usecases.UserUsecase{UserRepo: userRepo}
	userHandler := &http.UserHandler{Usecase: userUsecase}
	dealerHandler := &http.DealerHandler{Repo: dealerRepo}
	carImageHandler := &http.CarImageHandler{Repo: carImageRepo}
	adminHandler := &http.AdminHandler{Usecase: adminUsecase}
	authHandler := &http.AuthHandler{Usecase: authUsecase}

	// ========== ROUTES ==========
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
	)

	return app
}
