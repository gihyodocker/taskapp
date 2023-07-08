package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"golang.org/x/exp/slog"

	"github.com/gihyodocker/taskapp/pkg/id"
	"github.com/gihyodocker/taskapp/pkg/model"
	"github.com/gihyodocker/taskapp/pkg/payload"
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

func (h *Task) Update(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "id")
	if taskID == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	task, err := h.taskRepo.FindByID(r.Context(), taskID)
	if err == sql.ErrNoRows {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if err != nil {
		slog.Error("failed to get task", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	var input payload.Task
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		slog.Error("failed to decode json", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	task.Title = input.Title
	task.Content = input.Content
	task.Status = input.Status

	if err := h.taskRepo.Upsert(r.Context(), task); err != nil {
		slog.Error("failed to update task", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	slog.Info("updated task", slog.String("ID", taskID))

	w.WriteHeader(http.StatusNoContent)
}

func (h *Task) Delete(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "id")

	rows, err := h.taskRepo.DeleteByID(r.Context(), taskID)
	if err != nil {
		slog.Error("failed to delete task", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if rows == 0 {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	slog.Info("deleted task", slog.String("ID", taskID))

	w.WriteHeader(http.StatusNoContent)
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

func (h *Task) Create(w http.ResponseWriter, r *http.Request) {
	var input payload.Task
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		slog.Error("failed to decode json", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	m := model.Task{
		ID:      id.MakeULID().String(),
		Title:   input.Title,
		Content: input.Content,
		Status:  input.Status,
		Created: time.Now(),
		Updated: time.Now(),
	}

	if err := h.taskRepo.Upsert(r.Context(), &m); err != nil {
		slog.Error("failed to create task", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	slog.Info("created task", slog.String("ID", m.ID))
	w.WriteHeader(http.StatusCreated)
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
