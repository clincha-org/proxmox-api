package proxmox

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

const taskPath = "/tasks"

func (client *Client) GetTask(node string, id string) (Task, error) {
	request, err := http.NewRequest(
		"GET",
		client.Host+ApiPath+NodesPath+"/"+node+taskPath+"/"+id+"/status",
		nil,
	)
	if err != nil {
		return Task{}, fmt.Errorf("GetTask-build-request: %w", err)
	}

	request.AddCookie(&http.Cookie{Name: "PVEAuthCookie", Value: client.Ticket.Data.Ticket})
	request.Header.Set("CSRFPreventionToken", client.Ticket.Data.CSRFPreventionToken)

	response, err := client.HTTPClient.Do(request)
	if err != nil {
		return Task{}, fmt.Errorf("GetTask-do-request: %w", err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return Task{}, fmt.Errorf("GetTask-read-response: %w", err)
	}

	err = response.Body.Close()
	if err != nil {
		return Task{}, fmt.Errorf("GetTask-close-response: %w", err)
	}

	slog.Debug("api-response", "method", "GetTask", "status", response.Status, "response", string(body))

	if response.StatusCode != http.StatusOK {
		return Task{}, fmt.Errorf("GetTask-status-error: %s %s", response.Status, body)
	}

	taskModel := TaskResponse{}
	err = json.Unmarshal(body, &taskModel)
	if err != nil {
		return Task{}, fmt.Errorf("GetTask-unmarshal-response: %w", err)
	}

	return taskModel.Data, nil

}
