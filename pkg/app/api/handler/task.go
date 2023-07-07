package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"golang.org/x/exp/slog"

	"github.com/gihyodocker/taskapp/pkg/repository"
)

type Task struct {
	taskRepo repository.Task
}

func NewTask(taskRepo repository.Task) *Task {
	return &Task{
		taskRepo: taskRepo,
	}
}

func (h *Task) List(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.taskRepo.FindAll(r.Context())
	if err != nil {
		slog.Error("failed to get tasks", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		slog.Error("failed to marshal json", err)
	}
}

func (h *Task) Get(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "id")

	t, err := h.taskRepo.FindByID(r.Context(), taskID)
	if err == sql.ErrNoRows {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		slog.Error("failed to get task", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(t); err != nil {
		slog.Error("failed to marshal json", err)
	}
}
