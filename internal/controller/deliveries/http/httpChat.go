package http

import (
	"Backend_Go/internal/repositories"
	"Backend_Go/internal/usecases/chat"
	"Backend_Go/internal/ws"

	"github.com/gofiber/fiber/v2"
	websocket "github.com/gofiber/websocket/v2"
)

type ChatHandler struct {
	Usecase    *chat.ChatUsecase
	DealerRepo *repositories.DealerRepository
	Hub        *ws.Hub
}

// POST /api/chat/send
// Payload: { dealer_id: uint, car_id?: uint, content: string }
func (h *ChatHandler) SendMessage(c *fiber.Ctx) error {
	uid := c.Locals("user_id")
	if uid == nil {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}
	userID := uid.(uint)

	var req struct {
		DealerID uint   `json:"dealer_id"`
		CarID    *uint  `json:"car_id"`
		Content  string `json:"content"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.Usecase.SendMessageToDealer(userID, req.DealerID, req.CarID, req.Content); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "sent"})
}

// POST /api/chat/reply/:id
// Payload: { content: string }
func (h *ChatHandler) Reply(c *fiber.Ctx) error {
	uid := c.Locals("user_id")
	if uid == nil {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}
	userID := uid.(uint)
	convID, _ := c.ParamsInt("id")

	var req struct {
		Content string `json:"content"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// This handler is primarily for DEALERS replying, but customers could use it too
	// if we had a generic Reply.
	// Check Role
	role := c.Locals("role").(string)

	if role == "dealer" {
		if err := h.Usecase.ReplyToCustomer(userID, uint(convID), req.Content); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
	} else {
		// Customer replying to existing thread
		// We'd need to know the dealer ID from the conversation to reuse SendMessageToDealer?
		// Or Generic Reply Usecase.
		// For now, let's assume 'SendMessage' is main way for customers.
		// If customer calls this, it might fail or we implement generic Reply.
		return c.Status(400).JSON(fiber.Map{"error": "Use /api/chat/send for customers"})
	}

	return c.JSON(fiber.Map{"message": "replied"})
}

// GET /api/chat/conversations
func (h *ChatHandler) GetConversations(c *fiber.Ctx) error {
	uid := c.Locals("user_id")
	if uid == nil {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}
	userID := uid.(uint)
	role := c.Locals("role").(string)

	convs, err := h.Usecase.GetConversations(userID, role)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(convs)
}

// GET /api/chat/conversations/:id/messages
func (h *ChatHandler) GetMessages(c *fiber.Ctx) error {
	uid := c.Locals("user_id")
	if uid == nil {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}
	userID := uid.(uint)
	role := c.Locals("role").(string)
	convID, _ := c.ParamsInt("id")

	msgs, err := h.Usecase.GetMessages(uint(convID), userID, role)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(msgs)
}

// GET /api/chat/unread-count
func (h *ChatHandler) GetUnreadCount(c *fiber.Ctx) error {
	uid := c.Locals("user_id")
	if uid == nil {
		return c.Status(401).JSON(fiber.Map{"error": "unauthorized"})
	}
	userID := uid.(uint)
	role := c.Locals("role").(string)

	count, err := h.Usecase.GetTotalUnreadCount(userID, role)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"unread_count": count})
}

// GET /ws/chat
func (h *ChatHandler) HandleWebSocket(c *fiber.Ctx) error {
	// Handled by fiber-websocket middleware in routes
	return nil
}

func (h *ChatHandler) WebSocketUpgrade(c *websocket.Conn) {
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		c.Close()
		return
	}

	// Convert string to uint
	var userID uint
	for _, ch := range userIDStr {
		if ch >= '0' && ch <= '9' {
			userID = userID*10 + uint(ch-'0')
		}
	}

	if userID == 0 {
		c.Close()
		return
	}

	h.Hub.Register(userID, c)
	defer h.Hub.Unregister(userID)

	for {
		_, _, err := c.ReadMessage()
		if err != nil {
			break
		}
	}
}
