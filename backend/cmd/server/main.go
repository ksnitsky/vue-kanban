package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"kanban/internal/config"
	"kanban/internal/handler"
	"kanban/internal/middleware"
	"kanban/internal/repository"
	"kanban/internal/service"
	"kanban/internal/ws"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	cfg := config.Load()

	db, err := repository.New(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Connected to database")

	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	authTokenRepo := repository.NewAuthTokenRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	boardRepo := repository.NewBoardRepository(db)
	columnRepo := repository.NewColumnRepository(db)
	cardRepo := repository.NewCardRepository(db)

	authService := service.NewAuthService(
		authTokenRepo,
		userRepo,
		sessionRepo,
		time.Duration(cfg.SessionTTL)*time.Second,
	)

	hub := ws.NewHub()
	go hub.Run()

	authHandler := handler.NewAuthHandler(authService, hub)
	projectHandler := handler.NewProjectHandler(projectRepo)
	boardHandler := handler.NewBoardHandler(boardRepo, columnRepo, cardRepo, projectRepo)
	columnHandler := handler.NewColumnHandler(columnRepo, boardRepo, projectRepo)
	cardHandler := handler.NewCardHandler(cardRepo, columnRepo, boardRepo, projectRepo)
	telegramHandler := handler.NewTelegramHandler(authService, cfg.TelegramToken)

	mux := http.NewServeMux()

	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	})

	mux.HandleFunc("/api/auth/token", authHandler.GetToken)
	mux.HandleFunc("/api/auth/verify", authHandler.Verify)
	mux.HandleFunc("/api/auth/ws", authHandler.WebSocket)
	mux.HandleFunc("/api/auth/logout", middleware.Auth(authService)(http.HandlerFunc(authHandler.Logout)).ServeHTTP)
	mux.HandleFunc("/api/auth/me", middleware.Auth(authService)(http.HandlerFunc(authHandler.Me)).ServeHTTP)

	if cfg.Env == "development" {
		log.Println("Development mode: dev-login endpoint enabled")
		mux.HandleFunc("/api/auth/dev-login", authHandler.DevLogin)
	}

	mux.HandleFunc("/api/projects", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			middleware.Auth(authService)(http.HandlerFunc(projectHandler.List)).ServeHTTP(w, r)
		case http.MethodPost:
			middleware.Auth(authService)(http.HandlerFunc(projectHandler.Create)).ServeHTTP(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/projects/{id}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			middleware.Auth(authService)(http.HandlerFunc(projectHandler.Get)).ServeHTTP(w, r)
		case http.MethodPut:
			middleware.Auth(authService)(http.HandlerFunc(projectHandler.Update)).ServeHTTP(w, r)
		case http.MethodDelete:
			middleware.Auth(authService)(http.HandlerFunc(projectHandler.Delete)).ServeHTTP(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/projects/{project_id}/boards", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			middleware.Auth(authService)(http.HandlerFunc(boardHandler.List)).ServeHTTP(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/boards", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			middleware.Auth(authService)(http.HandlerFunc(boardHandler.Create)).ServeHTTP(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/boards/{id}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			middleware.Auth(authService)(http.HandlerFunc(boardHandler.Get)).ServeHTTP(w, r)
		case http.MethodPut:
			middleware.Auth(authService)(http.HandlerFunc(boardHandler.Update)).ServeHTTP(w, r)
		case http.MethodDelete:
			middleware.Auth(authService)(http.HandlerFunc(boardHandler.Delete)).ServeHTTP(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/columns", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			middleware.Auth(authService)(http.HandlerFunc(columnHandler.Create)).ServeHTTP(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/columns/{id}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			middleware.Auth(authService)(http.HandlerFunc(columnHandler.Update)).ServeHTTP(w, r)
		case http.MethodDelete:
			middleware.Auth(authService)(http.HandlerFunc(columnHandler.Delete)).ServeHTTP(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/columns/reorder", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			middleware.Auth(authService)(http.HandlerFunc(columnHandler.Reorder)).ServeHTTP(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/cards", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			middleware.Auth(authService)(http.HandlerFunc(cardHandler.Create)).ServeHTTP(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/cards/{id}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			middleware.Auth(authService)(http.HandlerFunc(cardHandler.Update)).ServeHTTP(w, r)
		case http.MethodDelete:
			middleware.Auth(authService)(http.HandlerFunc(cardHandler.Delete)).ServeHTTP(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/cards/move", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			middleware.Auth(authService)(http.HandlerFunc(cardHandler.Move)).ServeHTTP(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/cards/reorder", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			middleware.Auth(authService)(http.HandlerFunc(cardHandler.Reorder)).ServeHTTP(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/telegram/webhook", telegramHandler.Webhook)

	handler := middleware.CORS(mux)

	addr := ":" + cfg.Port
	log.Printf("Server starting on %s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
