package repositories

import (
	"Backend_Go/internal/entities"

	"gorm.io/gorm"
)

type ChatRepository struct{ DB *gorm.DB }

// Get or Create Conversation
func (r *ChatRepository) GetConversation(convID uint) (*entities.Conversation, error) {
	var conv entities.Conversation
	err := r.DB.First(&conv, convID).Error
	return &conv, err
}

func (r *ChatRepository) GetOrCreateConversation(userID, dealerID uint, carID *uint) (*entities.Conversation, error) {
	var conv entities.Conversation
	// Try to find existing conversation between user and dealer
	// Ideally distinct per car? Or per dealer?
	// Project requirements usually imply per dealer, but maybe context of car matters initially.
	// Let's stick to 1 conv per User-Dealer pair for simplicity, similar to FB Marketplace.
	// CarID is just context for the *start*.

	err := r.DB.Where("user_id = ? AND dealer_id = ?", userID, dealerID).First(&conv).Error
	if err == nil {
		return &conv, nil
	}
	if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// Create new
	conv = entities.Conversation{
		UserID:   userID,
		DealerID: dealerID,
		CarID:    carID,
	}
	if err := r.DB.Create(&conv).Error; err != nil {
		return nil, err
	}
	return &conv, nil
}

func (r *ChatRepository) CreateMessage(msg *entities.Message) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(msg).Error; err != nil {
			return err
		}
		// Update Conversation
		updates := map[string]interface{}{
			"last_message_id": msg.ID,
			"last_message":    msg.Content,
			"updated_at":      msg.CreatedAt,
		}

		var conv entities.Conversation
		if err := tx.First(&conv, msg.ConversationID).Error; err != nil {
			return err
		}

		if msg.SenderID == conv.UserID {
			// Sender is Customer -> Dealer gets unread
			updates["unread_count_dealer"] = gorm.Expr("unread_count_dealer + 1")
		} else {
			// Sender is Dealer -> Customer gets unread
			updates["unread_count_user"] = gorm.Expr("unread_count_user + 1")
		}

		return tx.Model(&entities.Conversation{}).Where("id = ?", msg.ConversationID).Updates(updates).Error
	})
}

func (r *ChatRepository) GetConversationsByUser(userID uint) ([]*entities.Conversation, error) {
	var convs []*entities.Conversation
	err := r.DB.
		Preload("Dealer").
		Preload("Car"). // Context car
		Where("user_id = ?", userID).
		Order("updated_at DESC").
		Find(&convs).Error
	return convs, err
}

func (r *ChatRepository) GetConversationsByDealer(dealerID uint) ([]*entities.Conversation, error) {
	var convs []*entities.Conversation
	err := r.DB.
		Preload("User").
		Preload("Car").
		Where("dealer_id = ?", dealerID).
		Order("updated_at DESC").
		Find(&convs).Error
	return convs, err
}

func (r *ChatRepository) GetMessages(convID uint) ([]*entities.Message, error) {
	var msgs []*entities.Message
	err := r.DB.Where("conversation_id = ?", convID).Order("created_at ASC").Find(&msgs).Error
	return msgs, err
}

func (r *ChatRepository) MarkReadByRole(convID uint, isDealer bool) error {
	updates := map[string]interface{}{}
	if isDealer {
		updates["unread_count_dealer"] = 0
	} else {
		updates["unread_count_user"] = 0
	}
	return r.DB.Model(&entities.Conversation{}).Where("id = ?", convID).Updates(updates).Error
}

// Better MarkRead with role
func (r *ChatRepository) GetTotalUnreadForUser(userID uint) (int, error) {
	var count int64
	err := r.DB.Model(&entities.Conversation{}).Where("user_id = ?", userID).Select("COALESCE(SUM(unread_count_user), 0)").Scan(&count).Error
	return int(count), err
}

func (r *ChatRepository) GetTotalUnreadForDealer(dealerID uint) (int, error) {
	var count int64
	err := r.DB.Model(&entities.Conversation{}).Where("dealer_id = ?", dealerID).Select("COALESCE(SUM(unread_count_dealer), 0)").Scan(&count).Error
	return int(count), err
}
