package entities

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       uint   `gorm:"primaryKey"`
	Name     string `json:"name"`
	Email    string `gorm:"uniqueIndex" json:"email"`
	Phone    string `gorm:"uniqueIndex" json:"phone"`
	Password string `json:"password"`
	Role     string `gorm:"type:varchar(20)" json:"role"`
	IsActive bool   `gorm:"default:true" json:"is_active"`
}

type Dealer struct {
	gorm.Model
	ID         uint   `gorm:"primaryKey" json:"id"`
	UserID     uint   `gorm:"uniqueIndex" json:"user_id"`
	ShopName   string `json:"shop_name"`
	Phone      string `json:"phone"`
	LineID     string `json:"line_id"`
	Address    string `json:"address"`
	Province   string `json:"province"`
	Latitude   string `json:"latitude"`
	Longitude  string `json:"longitude"`
	IsApproved bool   `gorm:"default:false" json:"is_approved"`

	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

type Car struct {
	gorm.Model
	DealerID     uint    `gorm:"index" json:"dealer_id"`
	Brand        string  `json:"brand"`
	ModelName    string  `json:"model_name"`
	Year         int     `json:"year"`
	Mileage      int     `json:"mileage"`
	Price        float64 `json:"price"`
	CarType      string  `json:"car_type"`
	FuelType     string  `json:"fuel_type"`
	Transmission string  `json:"transmission"`
	Color        string  `json:"color"`
	Description  string  `gorm:"type:text" json:"description"`
	Status       string  `gorm:"type:varchar(20)" json:"status"`
	Views        int     `gorm:"default:0" json:"views"`
	IsFeatured   bool    `gorm:"default:false" json:"is_featured"`
	// Contact / promotion statistics
	CallCount     int        `gorm:"default:0" json:"call_count"`
	LineCount     int        `gorm:"default:0" json:"line_count"`
	LeadCount     int        `gorm:"default:0" json:"lead_count"`
	IsPromoted    bool       `gorm:"default:false" json:"is_promoted"`
	PromotedUntil *time.Time `json:"promoted_until"`
	// Admin moderation
	IsHidden        bool   `gorm:"default:false" json:"is_hidden"`
	Flagged         bool   `gorm:"default:false" json:"flagged"`
	ViolationReason string `gorm:"type:text" json:"violation_reason"`

	CarImages []CarImage `gorm:"foreignKey:CarID" json:"car_images"`
	Dealer    Dealer     `gorm:"foreignKey:DealerID" json:"dealer,omitempty"`
}

type CarImage struct {
	gorm.Model
	CarID     uint   `gorm:"index" json:"car_id"`
	ImageURL  string `json:"image_url"`
	SortOrder int    `json:"sort_order"`

	Car Car `gorm:"foreignKey:CarID" json:"car,omitempty"`
}

type Lead struct {
	gorm.Model
	CarID      uint `gorm:"index"`
	DealerID   uint `gorm:"index"`
	CustomerID *uint
	ContactVia string `gorm:"type:varchar(20)"`
}

type Favorite struct {
	gorm.Model
	UserID uint `gorm:"index"`
	CarID  uint `gorm:"index"`
}

type Review struct {
	gorm.Model
	DealerID uint `gorm:"index"`
	UserID   uint `gorm:"index"`
	Rating   int
	Comment  string
}

type Report struct {
	gorm.Model
	CarID  uint `gorm:"index"`
	UserID *uint
	Reason string
}
