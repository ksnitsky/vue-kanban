package handler

import (
	"encoding/json"
	"net/http"

	"kanban/internal/model"
	"kanban/internal/repository"
)

type ProjectHandler struct {
	projectRepo *repository.ProjectRepository
}

func NewProjectHandler(projectRepo *repository.ProjectRepository) *ProjectHandler {
	return &ProjectHandler{projectRepo: projectRepo}
}

func (h *ProjectHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	projects, err := h.projectRepo.FindByUserID(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to fetch projects", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(projects)
}

func (h *ProjectHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	project := &model.Project{
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
	}

	if err := h.projectRepo.Create(r.Context(), project); err != nil {
		http.Error(w, "Failed to create project", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(project)
}

func (h *ProjectHandler) Get(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("id")
	userID := r.Context().Value("user_id").(string)

	project, err := h.projectRepo.FindByID(r.Context(), projectID)
	if err != nil {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	if project.UserID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	json.NewEncoder(w).Encode(project)
}

func (h *ProjectHandler) Update(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("id")
	userID := r.Context().Value("user_id").(string)

	project, err := h.projectRepo.FindByID(r.Context(), projectID)
	if err != nil {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	if project.UserID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	project.Name = req.Name
	project.Description = req.Description

	if err := h.projectRepo.Update(r.Context(), project); err != nil {
		http.Error(w, "Failed to update project", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(project)
}

func (h *ProjectHandler) Delete(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("id")
	userID := r.Context().Value("user_id").(string)

	project, err := h.projectRepo.FindByID(r.Context(), projectID)
	if err != nil {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	if project.UserID != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if err := h.projectRepo.Delete(r.Context(), projectID); err != nil {
		http.Error(w, "Failed to delete project", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
