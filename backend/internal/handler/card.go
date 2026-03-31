package handler

import (
	"encoding/json"
	"net/http"

	"kanban/internal/model"
	"kanban/internal/repository"
)

type CardHandler struct {
	cardRepo    *repository.CardRepository
	columnRepo  *repository.ColumnRepository
	boardRepo   *repository.BoardRepository
	projectRepo *repository.ProjectRepository
}

func NewCardHandler(
	cardRepo *repository.CardRepository,
	columnRepo *repository.ColumnRepository,
	boardRepo *repository.BoardRepository,
	projectRepo *repository.ProjectRepository,
) *CardHandler {
	return &CardHandler{
		cardRepo:    cardRepo,
		columnRepo:  columnRepo,
		boardRepo:   boardRepo,
		projectRepo: projectRepo,
	}
}

func (h *CardHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	var req struct {
		ColumnID string `json:"column_id"`
		Content  string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	column, err := h.columnRepo.FindByID(r.Context(), req.ColumnID)
	if err != nil {
		http.Error(w, "Column not found", http.StatusNotFound)
		return
	}

	board, err := h.boardRepo.FindByID(r.Context(), column.BoardID)
	if err != nil {
		http.Error(w, "Board not found", http.StatusNotFound)
		return
	}

	project, err := h.projectRepo.FindByID(r.Context(), board.ProjectID)
	if err != nil || project.UserID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	maxPos, _ := h.cardRepo.GetMaxPosition(r.Context(), req.ColumnID)

	card := &model.Card{
		ColumnID: req.ColumnID,
		Content:  req.Content,
		Position: maxPos + 1,
	}

	if err := h.cardRepo.Create(r.Context(), card); err != nil {
		http.Error(w, "Failed to create card", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(card)
}

func (h *CardHandler) Update(w http.ResponseWriter, r *http.Request) {
	cardID := r.PathValue("id")
	userID := r.Context().Value("user_id").(string)

	card, err := h.cardRepo.FindByID(r.Context(), cardID)
	if err != nil {
		http.Error(w, "Card not found", http.StatusNotFound)
		return
	}

	column, err := h.columnRepo.FindByID(r.Context(), card.ColumnID)
	if err != nil {
		http.Error(w, "Column not found", http.StatusNotFound)
		return
	}

	board, err := h.boardRepo.FindByID(r.Context(), column.BoardID)
	if err != nil {
		http.Error(w, "Board not found", http.StatusNotFound)
		return
	}

	project, err := h.projectRepo.FindByID(r.Context(), board.ProjectID)
	if err != nil || project.UserID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var req struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	card.Content = req.Content

	if err := h.cardRepo.Update(r.Context(), card); err != nil {
		http.Error(w, "Failed to update card", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(card)
}

func (h *CardHandler) Delete(w http.ResponseWriter, r *http.Request) {
	cardID := r.PathValue("id")
	userID := r.Context().Value("user_id").(string)

	card, err := h.cardRepo.FindByID(r.Context(), cardID)
	if err != nil {
		http.Error(w, "Card not found", http.StatusNotFound)
		return
	}

	column, err := h.columnRepo.FindByID(r.Context(), card.ColumnID)
	if err != nil {
		http.Error(w, "Column not found", http.StatusNotFound)
		return
	}

	board, err := h.boardRepo.FindByID(r.Context(), column.BoardID)
	if err != nil {
		http.Error(w, "Board not found", http.StatusNotFound)
		return
	}

	project, err := h.projectRepo.FindByID(r.Context(), board.ProjectID)
	if err != nil || project.UserID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if err := h.cardRepo.Delete(r.Context(), cardID); err != nil {
		http.Error(w, "Failed to delete card", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *CardHandler) Move(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	var req struct {
		CardID         string `json:"card_id"`
		TargetColumnID string `json:"target_column_id"`
		Position       int    `json:"position"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	card, err := h.cardRepo.FindByID(r.Context(), req.CardID)
	if err != nil {
		http.Error(w, "Card not found", http.StatusNotFound)
		return
	}

	targetColumn, err := h.columnRepo.FindByID(r.Context(), req.TargetColumnID)
	if err != nil {
		http.Error(w, "Target column not found", http.StatusNotFound)
		return
	}

	sourceColumn, err := h.columnRepo.FindByID(r.Context(), card.ColumnID)
	if err != nil {
		http.Error(w, "Source column not found", http.StatusNotFound)
		return
	}

	sourceBoard, err := h.boardRepo.FindByID(r.Context(), sourceColumn.BoardID)
	if err != nil {
		http.Error(w, "Board not found", http.StatusNotFound)
		return
	}

	targetBoard, err := h.boardRepo.FindByID(r.Context(), targetColumn.BoardID)
	if err != nil {
		http.Error(w, "Target board not found", http.StatusNotFound)
		return
	}

	sourceProject, err := h.projectRepo.FindByID(r.Context(), sourceBoard.ProjectID)
	if err != nil || sourceProject.UserID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	targetProject, err := h.projectRepo.FindByID(r.Context(), targetBoard.ProjectID)
	if err != nil || targetProject.UserID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if err := h.cardRepo.Move(r.Context(), req.CardID, req.TargetColumnID, req.Position); err != nil {
		http.Error(w, "Failed to move card", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *CardHandler) Reorder(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	var req struct {
		ColumnID string   `json:"column_id"`
		CardIDs  []string `json:"card_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	column, err := h.columnRepo.FindByID(r.Context(), req.ColumnID)
	if err != nil {
		http.Error(w, "Column not found", http.StatusNotFound)
		return
	}

	board, err := h.boardRepo.FindByID(r.Context(), column.BoardID)
	if err != nil {
		http.Error(w, "Board not found", http.StatusNotFound)
		return
	}

	project, err := h.projectRepo.FindByID(r.Context(), board.ProjectID)
	if err != nil || project.UserID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if err := h.cardRepo.Reorder(r.Context(), req.ColumnID, req.CardIDs); err != nil {
		http.Error(w, "Failed to reorder cards", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
