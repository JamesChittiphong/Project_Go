package entities

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       uint   `gorm:"primaryKey"`
	Name     string `json:"name"`
	Email    string `gorm:"uniqueIndex"`
	Password string
	Role     string `gorm:"type:varchar(20)"` // customer, dealer, admin
	IsActive bool   `gorm:"default:true"`
}

type Dealer struct {
	gorm.Model
	ID         uint `gorm:"primaryKey"`
	UserID     uint `gorm:"uniqueIndex"`
	ShopName   string
	Phone      string
	LineID     string
	Address    string
	Province   string
	Latitude   string
	Longitude  string
	IsApproved bool `gorm:"default:false"`

	User User `gorm:"foreignKey:UserID"`
}

type Car struct {
	gorm.Model
	DealerID     uint `gorm:"index"`
	Brand        string
	ModelName    string
	Year         int
	Mileage      int
	Price        float64
	CarType      string // sedan, suv, pickup
	FuelType     string
	Transmission string
	Color        string
	Description  string
	Status       string `gorm:"type:varchar(20)"` // available, contacted, sold
	Views        int    `gorm:"default:0"`
	IsFeatured   bool   `gorm:"default:false"`
	// Contact / promotion statistics
	CallCount     int  `gorm:"default:0"`
	LineCount     int  `gorm:"default:0"`
	LeadCount     int  `gorm:"default:0"`
	IsPromoted    bool `gorm:"default:false"`
	PromotedUntil *time.Time
	// Admin moderation
	IsHidden        bool   `gorm:"default:false"`
	Flagged         bool   `gorm:"default:false"`
	ViolationReason string `gorm:"type:text"`

	Dealer Dealer `gorm:"foreignKey:DealerID"`
}

type CarImage struct {
	gorm.Model
	CarID     uint `gorm:"index"`
	ImageURL  string
	SortOrder int

	Car Car `gorm:"foreignKey:CarID"`
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
