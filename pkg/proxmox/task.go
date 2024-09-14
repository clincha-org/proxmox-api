package proxmox

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

const taskPath = "/tasks"

func (client *Client) GetTaskStatus(node string, id string) (Task, error) {
	request, err := http.NewRequest(
		"GET",
		client.Host+ApiPath+NodesPath+"/"+node+taskPath+"/"+id+"/status",
		nil,
	)
	if err != nil {
		return Task{}, fmt.Errorf("GetTaskStatus-build-request: %w", err)
	}

	request.AddCookie(&http.Cookie{Name: "PVEAuthCookie", Value: client.Ticket.Data.Ticket})
	request.Header.Set("CSRFPreventionToken", client.Ticket.Data.CSRFPreventionToken)

	response, err := client.HTTPClient.Do(request)
	if err != nil {
		return Task{}, fmt.Errorf("GetTaskStatus-do-request: %w", err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return Task{}, fmt.Errorf("GetTaskStatus-read-response: %w", err)
	}

	err = response.Body.Close()
	if err != nil {
		return Task{}, fmt.Errorf("GetTaskStatus-close-response: %w", err)
	}

	slog.Debug("api-response", "method", "GetTaskStatus", "status", response.Status, "response", string(body))

	if response.StatusCode != http.StatusOK {
		return Task{}, fmt.Errorf("GetTaskStatus-status-error: %s %s", response.Status, body)
	}

	taskModel := TaskResponse{}
	err = json.Unmarshal(body, &taskModel)
	if err != nil {
		return Task{}, fmt.Errorf("GetTaskStatus-unmarshal-response: %w", err)
	}

	return taskModel.Data, nil

}

func (client *Client) AwaitAsynchronousTask(node string, taskID string) error {
	task, err := client.GetTaskStatus(node, taskID)
	if err != nil {
		return fmt.Errorf("AwaitAsynchronousTask-get-task-status: '%s' %w", taskID, err)
	}

	// Poll the task status until the task is completed
	for ok := true; ok; ok = task.Status != "stopped" {
		task, err = client.GetTaskStatus(node, taskID)
		if err != nil {
			return fmt.Errorf("AwaitAsynchronousTask-get-job-loop: '%s' %w", taskID, err)
		}
		// Sleep for 1 second before polling again
		time.Sleep(1 * time.Second)
	}

	return nil
}
