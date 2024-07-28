package proxmox

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const NodesPath string = "nodes"

func (client *Client) GetNodes() ([]Node, error) {
	var node []Node
	request, err := http.NewRequest(
		"GET",
		client.Host+ApiPath+NodesPath,
		nil,
	)
	if err != nil {
		return node, fmt.Errorf("GetNodes-build-request: %w", err)
	}

	request.AddCookie(&http.Cookie{Name: "PVEAuthCookie", Value: client.Ticket.Data.Ticket})
	request.Header.Set("CSRFPreventionToken", client.Ticket.Data.CSRFPreventionToken)

	response, err := client.HTTPClient.Do(request)
	if err != nil {
		return node, fmt.Errorf("GetNodes-do-request: %w", err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return node, fmt.Errorf("GetNodes-read-response: %w", err)
	}

	err = response.Body.Close()
	if err != nil {
		return node, fmt.Errorf("GetNodes-close-response: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		return node, fmt.Errorf("GetNodes-status-error: %s %s", response.Status, body)
	}

	nodeModel := NodeResponse{}
	err = json.Unmarshal(body, &nodeModel)
	if err != nil {
		return node, fmt.Errorf("GetNodes-unmarshal-response: %w", err)
	}

	node = nodeModel.Data

	return node, nil
}
