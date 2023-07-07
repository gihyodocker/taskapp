package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gihyodocker/taskapp/pkg/model"
)

type TaskClient interface {
	List() ([]*model.Task, error)
}

type taskClient struct {
	apiAddress string
	httpCli    *http.Client
}

func NewTask(apiAddress string) TaskClient {
	return &taskClient{
		apiAddress: apiAddress,
		httpCli:    &http.Client{},
	}
}

func (c taskClient) List() ([]*model.Task, error) {
	url := fmt.Sprintf("%s/api/tasks", c.apiAddress)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpCli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to request: %w", err)
	}
	defer resp.Body.Close()

	var tasks []*model.Task
	if err := json.NewDecoder(resp.Body).Decode(&tasks); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return tasks, nil
}

func (c taskClient) Get(id string) (*model.Task, error) {
	url := fmt.Sprintf("%s/api/%s", c.apiAddress, id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpCli.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to request: %w", err)
	}
	defer resp.Body.Close()

	var task model.Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &task, nil
}
