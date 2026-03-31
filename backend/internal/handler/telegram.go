package handler

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"kanban/internal/service"
)

type TelegramHandler struct {
	authService *service.AuthService
	botToken    string
}

func NewTelegramHandler(authService *service.AuthService, botToken string) *TelegramHandler {
	return &TelegramHandler{
		authService: authService,
		botToken:    botToken,
	}
}

type TelegramUpdate struct {
	UpdateID int              `json:"update_id"`
	Message  *TelegramMessage `json:"message"`
}

type TelegramMessage struct {
	MessageID int           `json:"message_id"`
	From      *TelegramUser `json:"from"`
	Chat      *TelegramChat `json:"chat"`
	Text      string        `json:"text"`
}

type TelegramChat struct {
	ID int64 `json:"id"`
}

func (h *TelegramHandler) Webhook(w http.ResponseWriter, r *http.Request) {
	var update TelegramUpdate
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if update.Message == nil || update.Message.Text == "" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if update.Message.From == nil {
		w.WriteHeader(http.StatusOK)
		return
	}

	text := update.Message.Text
	if len(text) >= 7 && text[:7] == "/start " {
		token := text[7:]
		go h.handleStartCommand(token, update.Message.From)
	}

	w.WriteHeader(http.StatusOK)
}

func (h *TelegramHandler) handleStartCommand(token string, user *TelegramUser) {
	verifyReq := VerifyRequest{
		Token: token,
		User:  *user,
	}

	reqBody, _ := json.Marshal(verifyReq)

	resp, err := http.Post("http://localhost:8080/api/auth/verify", "application/json", bytes.NewReader(reqBody))
	if err != nil {
		log.Printf("Failed to verify token: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Verify failed with status: %d", resp.StatusCode)
		return
	}

	log.Printf("User %s authenticated successfully", user.Username)
}
