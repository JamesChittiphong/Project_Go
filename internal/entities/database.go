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
	Status     string `gorm:"default:'pending';type:varchar(20)" json:"status"` // pending, approved, suspended
	IsApproved bool   `gorm:"default:false" json:"is_approved"`                 // Keep for backward compatibility or remove later

	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

type Car struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

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
	Status       string  `gorm:"default:'pending';type:varchar(20)" json:"status"` // pending, approved, rejected, delete_requested, selling, sold
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

	Car Car `gorm:"foreignKey:CarID" json:"-"`
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
	UserID uint `gorm:"uniqueIndex:idx_fav_user_car" json:"user_id"`
	CarID  uint `gorm:"uniqueIndex:idx_fav_user_car" json:"car_id"`

	Car Car `gorm:"foreignKey:CarID" json:"car"`
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

type Conversation struct {
	gorm.Model
	UserID            uint      `gorm:"uniqueIndex:idx_user_dealer;index" json:"user_id"`   // Customer User ID
	DealerID          uint      `gorm:"uniqueIndex:idx_user_dealer;index" json:"dealer_id"` // Dealer ID (from Dealer table)
	CarID             *uint     `json:"car_id"`
	Topic             string    `json:"topic"`
	UnreadCountUser   int       `gorm:"default:0" json:"unread_count_user"`
	UnreadCountDealer int       `gorm:"default:0" json:"unread_count_dealer"`
	LastMessageID     *uint     `gorm:"index" json:"last_message_id"`
	LastMessage       string    `gorm:"type:text" json:"last_message"` // Cache content for preview
	UpdatedAt         time.Time `json:"updated_at"`

	User   User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Dealer Dealer `gorm:"foreignKey:DealerID" json:"dealer,omitempty"`
	Car    *Car   `gorm:"foreignKey:CarID" json:"car,omitempty"`
}

type Message struct {
	gorm.Model
	ConversationID uint   `gorm:"index"`
	SenderID       uint   `gorm:"index"` // UserID of sender
	Content        string `gorm:"type:text"`
	MsgType        string `gorm:"type:varchar(20);default:'text'"` // text, image
	IsRead         bool   `gorm:"default:false"`

	Conversation Conversation `gorm:"foreignKey:ConversationID"`
}
