package handler

import (
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
	"golang.org/x/exp/slog"

	"github.com/gihyodocker/taskapp/pkg/app/web/client"
)

type Delete struct {
	taskCli client.TaskClient
}

func NewDelete(taskCli client.TaskClient) *Delete {
	return &Delete{
		taskCli: taskCli,
	}
}

func (p *Delete) Confirm(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "id")
	if taskID == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	task, err := p.taskCli.Get(taskID)
	if err != nil {
		slog.Error("failed to get task", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	param := toTaskParam(task)

	tmpl := template.Must(template.ParseFS(templateFS, "template/delete.html"))
	if err := tmpl.Execute(w, param); err != nil {
		slog.Error("failed to execute template", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (p *Delete) Complete(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "id")
	if taskID == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if err := p.taskCli.Delete(taskID); err != nil {
		slog.Error("failed to delete task", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("location", "/")
	w.WriteHeader(http.StatusMovedPermanently)
}
