package handler

import (
	"encoding/json"
	"net/http"

	"kanban/internal/model"
	"kanban/internal/repository"
)

type BoardHandler struct {
	boardRepo   *repository.BoardRepository
	columnRepo  *repository.ColumnRepository
	cardRepo    *repository.CardRepository
	projectRepo *repository.ProjectRepository
}

func NewBoardHandler(
	boardRepo *repository.BoardRepository,
	columnRepo *repository.ColumnRepository,
	cardRepo *repository.CardRepository,
	projectRepo *repository.ProjectRepository,
) *BoardHandler {
	return &BoardHandler{
		boardRepo:   boardRepo,
		columnRepo:  columnRepo,
		cardRepo:    cardRepo,
		projectRepo: projectRepo,
	}
}

type BoardWithColumns struct {
	*model.Board
	Columns []ColumnWithCards `json:"columns"`
}

type ColumnWithCards struct {
	*model.Column
	Cards []*model.Card `json:"cards"`
}

func (h *BoardHandler) List(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("project_id")
	userID := r.Context().Value("user_id").(string)

	project, err := h.projectRepo.FindByID(r.Context(), projectID)
	if err != nil || project.UserID != userID {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	boards, err := h.boardRepo.FindByProjectID(r.Context(), projectID)
	if err != nil {
		http.Error(w, "Failed to fetch boards", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(boards)
}

func (h *BoardHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	var req struct {
		ProjectID string `json:"project_id"`
		Name      string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	project, err := h.projectRepo.FindByID(r.Context(), req.ProjectID)
	if err != nil || project.UserID != userID {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	board := &model.Board{
		ProjectID: req.ProjectID,
		Name:      req.Name,
	}

	if err := h.boardRepo.Create(r.Context(), board); err != nil {
		http.Error(w, "Failed to create board", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(board)
}

func (h *BoardHandler) Get(w http.ResponseWriter, r *http.Request) {
	boardID := r.PathValue("id")
	userID := r.Context().Value("user_id").(string)

	board, err := h.boardRepo.FindByID(r.Context(), boardID)
	if err != nil {
		http.Error(w, "Board not found", http.StatusNotFound)
		return
	}

	project, err := h.projectRepo.FindByID(r.Context(), board.ProjectID)
	if err != nil || project.UserID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	columns, err := h.columnRepo.FindByBoardID(r.Context(), boardID)
	if err != nil {
		http.Error(w, "Failed to fetch columns", http.StatusInternalServerError)
		return
	}

	columnsWithCards := make([]ColumnWithCards, len(columns))
	for i, col := range columns {
		cards, err := h.cardRepo.FindByColumnID(r.Context(), col.ID)
		if err != nil {
			http.Error(w, "Failed to fetch cards", http.StatusInternalServerError)
			return
		}
		columnsWithCards[i] = ColumnWithCards{
			Column: col,
			Cards:  cards,
		}
	}

	result := BoardWithColumns{
		Board:   board,
		Columns: columnsWithCards,
	}

	json.NewEncoder(w).Encode(result)
}

func (h *BoardHandler) Update(w http.ResponseWriter, r *http.Request) {
	boardID := r.PathValue("id")
	userID := r.Context().Value("user_id").(string)

	board, err := h.boardRepo.FindByID(r.Context(), boardID)
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
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	board.Name = req.Name

	if err := h.boardRepo.Update(r.Context(), board); err != nil {
		http.Error(w, "Failed to update board", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(board)
}

func (h *BoardHandler) Delete(w http.ResponseWriter, r *http.Request) {
	boardID := r.PathValue("id")
	userID := r.Context().Value("user_id").(string)

	board, err := h.boardRepo.FindByID(r.Context(), boardID)
	if err != nil {
		http.Error(w, "Board not found", http.StatusNotFound)
		return
	}

	project, err := h.projectRepo.FindByID(r.Context(), board.ProjectID)
	if err != nil || project.UserID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if err := h.boardRepo.Delete(r.Context(), boardID); err != nil {
		http.Error(w, "Failed to delete board", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
