package chat

import (
	"Backend_Go/internal/entities"
	"Backend_Go/internal/repositories"
	"Backend_Go/internal/ws"
)

type ChatUsecase struct {
	ChatRepo   *repositories.ChatRepository
	DealerRepo *repositories.DealerRepository
	Hub        *ws.Hub
}

func (u *ChatUsecase) SendMessageToDealer(userID uint, dealerID uint, carID *uint, content string) error {
	conv, err := u.ChatRepo.GetOrCreateConversation(userID, dealerID, carID)
	if err != nil {
		return err
	}

	msg := &entities.Message{
		ConversationID: conv.ID,
		SenderID:       userID,
		Content:        content,
		MsgType:        "text",
	}
	if err := u.ChatRepo.CreateMessage(msg); err != nil {
		return err
	}

	// Broadcast to Dealer
	// Need to find Dealer's UserID
	var dealer entities.Dealer
	if err := u.DealerRepo.FindByID(dealerID, &dealer); err == nil {
		u.Hub.BroadcastToUser(dealer.UserID, map[string]interface{}{
			"type":            "new_message",
			"conversation_id": conv.ID,
			"message":         msg,
		})
	}

	return nil
}

func (u *ChatUsecase) ReplyToCustomer(dealerUserID uint, convID uint, content string) error {
	// Verify dealer owns this conversation
	msg := &entities.Message{
		ConversationID: convID,
		SenderID:       dealerUserID,
		Content:        content,
		MsgType:        "text",
	}
	if err := u.ChatRepo.CreateMessage(msg); err != nil {
		return err
	}

	// Fetch conversation to find Customer's UserID
	conv, err := u.ChatRepo.GetConversation(convID)
	if err == nil {
		u.Hub.BroadcastToUser(conv.UserID, map[string]interface{}{
			"type":            "new_message",
			"conversation_id": conv.ID,
			"message":         msg,
		})
	}

	return nil
}

func (u *ChatUsecase) GetTotalUnreadCount(userID uint, role string) (int, error) {
	if role == "dealer" {
		var dealer entities.Dealer
		if err := u.DealerRepo.FindByUserID(userID, &dealer); err != nil {
			return 0, err
		}
		return u.ChatRepo.GetTotalUnreadForDealer(dealer.ID)
	}
	return u.ChatRepo.GetTotalUnreadForUser(userID)
}

func (u *ChatUsecase) GetConversations(userID uint, role string) ([]*entities.Conversation, error) {
	if role == "dealer" {
		// Find dealer profile for this user
		var dealer entities.Dealer
		if err := u.DealerRepo.FindByUserID(userID, &dealer); err != nil {
			return nil, err
		}
		return u.ChatRepo.GetConversationsByDealer(dealer.ID)
	}
	return u.ChatRepo.GetConversationsByUser(userID)
}

func (u *ChatUsecase) GetMessages(convID uint, userID uint, role string) ([]*entities.Message, error) {
	// Verify access
	// Fetch messages
	msgs, err := u.ChatRepo.GetMessages(convID)
	if err != nil {
		return nil, err
	}

	// Mark read
	isDealer := (role == "dealer")
	u.ChatRepo.MarkReadByRole(convID, isDealer)

	return msgs, nil
}
