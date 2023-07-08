package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gihyodocker/taskapp/pkg/model"
	"github.com/gihyodocker/taskapp/pkg/payload"
)

type TaskClient interface {
	Update(id string, input payload.Task) error
	Delete(id string) error
	Get(id string) (*model.Task, error)
	Create(input payload.Task) error
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

func (c taskClient) Update(id string, input payload.Task) error {
	url := fmt.Sprintf("%s/api/tasks/%s", c.apiAddress, id)
	body, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("failed to marshal json: %w", err)
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpCli.Do(req)
	if err != nil {
		return fmt.Errorf("failed to request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to update task: %s", resp.Status)
	}
	return nil
}

func (c taskClient) Delete(id string) error {
	url := fmt.Sprintf("%s/api/tasks/%s", c.apiAddress, id)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpCli.Do(req)
	if err != nil {
		return fmt.Errorf("failed to request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete task: %s", resp.Status)
	}
	return nil
}

func (c taskClient) Get(id string) (*model.Task, error) {
	url := fmt.Sprintf("%s/api/tasks/%s", c.apiAddress, id)
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

func (c taskClient) Create(input payload.Task) error {
	url := fmt.Sprintf("%s/api/tasks", c.apiAddress)
	body, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("failed to marshal json: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpCli.Do(req)
	if err != nil {
		return fmt.Errorf("failed to request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to create task: %s", resp.Status)
	}
	return nil
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
