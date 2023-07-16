package handler

import (
	"html/template"
	"net/http"

	"golang.org/x/exp/slog"

	"github.com/gihyodocker/taskapp/pkg/app/web/client"
	"github.com/gihyodocker/taskapp/pkg/model"
)

type Index struct {
	taskCli client.TaskClient
}

func NewIndex(taskCli client.TaskClient) *Index {
	return &Index{
		taskCli: taskCli,
	}
}

type indexParam struct {
	Backlog  []*model.Task
	Progress []*model.Task
	Done     []*model.Task
}

func (p *Index) Index(w http.ResponseWriter, r *http.Request) {
	tasks, err := p.taskCli.List()
	if err != nil {
		slog.Error("failed to get tasks", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	param := indexParam{
		Backlog:  make([]*model.Task, 0),
		Progress: make([]*model.Task, 0),
		Done:     make([]*model.Task, 0),
	}

	for _, t := range tasks {
		switch t.Status {
		case model.TaskStatusBACKLOG:
			param.Backlog = append(param.Backlog, t)
		case model.TaskStatusPROGRESS:
			param.Progress = append(param.Progress, t)
		case model.TaskStatusDONE:
			param.Done = append(param.Done, t)
		default:
			slog.Error("unknown status: %s", t.Status)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

	tmpl := template.Must(template.ParseFS(templateFS, "template/index.html"))
	if err := tmpl.Execute(w, param); err != nil {
		slog.Error("failed to execute template", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
