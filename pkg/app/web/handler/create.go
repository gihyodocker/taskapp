package handler

import (
	"html/template"
	"net/http"

	"golang.org/x/exp/slog"

	"github.com/gihyodocker/taskapp/pkg/app/web/client"
	"github.com/gihyodocker/taskapp/pkg/payload"
)

type Create struct {
	taskCli client.TaskClient
}

func NewCreate(taskCli client.TaskClient) *Create {
	return &Create{
		taskCli: taskCli,
	}
}

func (p *Create) Input(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFS(templateFS, "template/create.html"))
	if err := tmpl.Execute(w, nil); err != nil {
		slog.Error("failed to execute template", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (p *Create) Complete(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	title := r.Form.Get("title")
	content := r.Form.Get("content")
	status := r.Form.Get("status")

	input := payload.Task{
		Title:   title,
		Content: content,
		Status:  status,
	}

	if err := p.taskCli.Create(input); err != nil {
		slog.Error("failed to create task", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("location", "/")
	w.WriteHeader(http.StatusMovedPermanently)
}
