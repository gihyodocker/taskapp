package handler

import (
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
	"golang.org/x/exp/slog"

	"github.com/gihyodocker/taskapp/pkg/app/web/client"
	"github.com/gihyodocker/taskapp/pkg/payload"
)

type Update struct {
	taskCli client.TaskClient
}

func NewUpdate(taskCli client.TaskClient) *Update {
	return &Update{
		taskCli: taskCli,
	}
}

func (p *Update) Input(w http.ResponseWriter, r *http.Request) {
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

	tmpl := template.Must(template.ParseFS(templateFS, "template/update.html"))
	if err := tmpl.Execute(w, param); err != nil {
		slog.Error("failed to execute template", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (p *Update) Complete(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "id")
	if taskID == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	r.ParseForm()
	title := r.Form.Get("title")
	content := r.Form.Get("content")
	status := r.Form.Get("status")

	input := payload.Task{
		Title:   title,
		Content: content,
		Status:  status,
	}

	if err := p.taskCli.Update(taskID, input); err != nil {
		slog.Error("failed to update task", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("location", "/")
	w.WriteHeader(http.StatusMovedPermanently)
}
