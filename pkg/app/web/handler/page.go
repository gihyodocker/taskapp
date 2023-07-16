package handler

import (
	"embed"

	"github.com/gihyodocker/taskapp/pkg/model"
)

//go:embed template
var templateFS embed.FS

type taskParam struct {
	ID      string
	Title   string
	Content string
	Status  string
}

func toTaskParam(m *model.Task) taskParam {
	return taskParam{
		ID:      m.ID,
		Title:   m.Title,
		Content: m.Content,
		Status:  m.Status,
	}
}
