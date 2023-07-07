package page

import (
	"embed"
	"html/template"
	"net/http"

	"golang.org/x/exp/slog"

	"github.com/gihyodocker/taskapp/pkg/app/web/client"
)

type Index struct {
	taskCli client.TaskClient
}

func NewIndex(taskCli client.TaskClient) *Index {
	return &Index{
		taskCli: taskCli,
	}
}

//go:embed template
var templateFS embed.FS

func (p *Index) Get(w http.ResponseWriter, r *http.Request) {
	tasks, err := p.taskCli.List()
	if err != nil {
		slog.Error("failed to get tasks", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFS(templateFS, "template/index.html"))
	if err := tmpl.Execute(w, tasks); err != nil {
		slog.Error("failed to execute template", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
