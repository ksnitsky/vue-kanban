package handler

import (
	"encoding/json"
	"log"
	"net"
	"net/http"

	"kanban/internal/model"
	"kanban/internal/service"
	"kanban/internal/ws"
)

type AuthHandler struct {
	authService *service.AuthService
	hub         *ws.Hub
}

func NewAuthHandler(authService *service.AuthService, hub *ws.Hub) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		hub:         hub,
	}
}

type TokenResponse struct {
	Token string `json:"token"`
}

func (h *AuthHandler) GetToken(w http.ResponseWriter, r *http.Request) {
	token, err := h.authService.GenerateToken(r.Context())
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(TokenResponse{Token: token.Token})
}

type VerifyRequest struct {
	Token string       `json:"token"`
	User  TelegramUser `json:"user"`
}

type TelegramUser struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	PhotoURL  string `json:"photo_url"`
}

type AuthSuccessMessage struct {
	Type string      `json:"type"`
	User *model.User `json:"user"`
}

func (h *AuthHandler) Verify(w http.ResponseWriter, r *http.Request) {
	var req VerifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	telegramUser := &model.User{
		TelegramID: req.User.ID,
		Username:   req.User.Username,
		FirstName:  req.User.FirstName,
		LastName:   req.User.LastName,
		PhotoURL:   req.User.PhotoURL,
	}

	user, err := h.authService.VerifyToken(r.Context(), req.Token, telegramUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Extract IP from RemoteAddr (format: IP:port)
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		ip = r.RemoteAddr
	}

	sessionToken, err := h.authService.CreateSession(
		r.Context(),
		user.ID,
		r.UserAgent(),
		ip,
	)
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionToken,
		Path:     "/",
		MaxAge:   604800,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	message := AuthSuccessMessage{
		Type: "auth_success",
		User: user,
	}
	messageBytes, _ := json.Marshal(message)

	if !h.hub.SendToToken(req.Token, messageBytes) {
		log.Printf("Failed to send auth success to token: %s", req.Token)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"user":    user,
	})
}

func (h *AuthHandler) WebSocket(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Token required", http.StatusBadRequest)
		return
	}

	conn, err := ws.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := ws.NewClient(h.hub, conn, token)
	h.hub.Register(client)

	go client.WritePump()
	go client.ReadPump()
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "No session", http.StatusBadRequest)
		return
	}

	session, err := h.authService.ValidateSession(r.Context(), cookie.Value)
	if err != nil {
		http.Error(w, "Invalid session", http.StatusBadRequest)
		return
	}

	if err := h.authService.Logout(r.Context(), session.ID); err != nil {
		http.Error(w, "Failed to logout", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "session_id",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	user, err := h.authService.GetUser(r.Context(), userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) DevLogin(w http.ResponseWriter, r *http.Request) {
	devUser := &model.User{
		TelegramID: 123456789,
		Username:   "dev_user",
		FirstName:  "Dev",
		LastName:   "User",
		PhotoURL:   "",
	}

	user, err := h.authService.CreateOrGetUser(r.Context(), devUser)
	if err != nil {
		http.Error(w, "Failed to create dev user", http.StatusInternalServerError)
		return
	}

	// Extract IP from RemoteAddr (format: IP:port)
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		ip = r.RemoteAddr
	}

	sessionToken, err := h.authService.CreateSession(
		r.Context(),
		user.ID,
		r.UserAgent(),
		ip,
	)
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionToken,
		Path:     "/",
		MaxAge:   604800,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
