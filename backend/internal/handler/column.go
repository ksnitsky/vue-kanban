package handler

import (
	"encoding/json"
	"net/http"

	"kanban/internal/model"
	"kanban/internal/repository"
)

type ColumnHandler struct {
	columnRepo  *repository.ColumnRepository
	boardRepo   *repository.BoardRepository
	projectRepo *repository.ProjectRepository
}

func NewColumnHandler(
	columnRepo *repository.ColumnRepository,
	boardRepo *repository.BoardRepository,
	projectRepo *repository.ProjectRepository,
) *ColumnHandler {
	return &ColumnHandler{
		columnRepo:  columnRepo,
		boardRepo:   boardRepo,
		projectRepo: projectRepo,
	}
}

func (h *ColumnHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	var req struct {
		BoardID string `json:"board_id"`
		Title   string `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	board, err := h.boardRepo.FindByID(r.Context(), req.BoardID)
	if err != nil {
		http.Error(w, "Board not found", http.StatusNotFound)
		return
	}

	project, err := h.projectRepo.FindByID(r.Context(), board.ProjectID)
	if err != nil || project.UserID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	maxPos, _ := h.columnRepo.GetMaxPosition(r.Context(), req.BoardID)

	column := &model.Column{
		BoardID:  req.BoardID,
		Title:    req.Title,
		Position: maxPos + 1,
	}

	if err := h.columnRepo.Create(r.Context(), column); err != nil {
		http.Error(w, "Failed to create column", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(column)
}

func (h *ColumnHandler) Update(w http.ResponseWriter, r *http.Request) {
	columnID := r.PathValue("id")
	userID := r.Context().Value("user_id").(string)

	column, err := h.columnRepo.FindByID(r.Context(), columnID)
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
		Title string `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	column.Title = req.Title

	if err := h.columnRepo.Update(r.Context(), column); err != nil {
		http.Error(w, "Failed to update column", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(column)
}

func (h *ColumnHandler) Delete(w http.ResponseWriter, r *http.Request) {
	columnID := r.PathValue("id")
	userID := r.Context().Value("user_id").(string)

	column, err := h.columnRepo.FindByID(r.Context(), columnID)
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

	if err := h.columnRepo.Delete(r.Context(), columnID); err != nil {
		http.Error(w, "Failed to delete column", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ColumnHandler) Reorder(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	var req struct {
		BoardID   string   `json:"board_id"`
		ColumnIDs []string `json:"column_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	board, err := h.boardRepo.FindByID(r.Context(), req.BoardID)
	if err != nil {
		http.Error(w, "Board not found", http.StatusNotFound)
		return
	}

	project, err := h.projectRepo.FindByID(r.Context(), board.ProjectID)
	if err != nil || project.UserID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if err := h.columnRepo.Reorder(r.Context(), req.ColumnIDs); err != nil {
		http.Error(w, "Failed to reorder columns", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
