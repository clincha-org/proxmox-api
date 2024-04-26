package proxmox

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type NodeModel struct {
	Data []Node `json:"data"`
}

type Node struct {
	Type           string  `json:"type"`
	Maxcpu         int     `json:"maxcpu"`
	Cpu            float64 `json:"cpu"`
	Status         string  `json:"status"`
	Maxmem         int     `json:"maxmem"`
	SslFingerprint string  `json:"ssl_fingerprint"`
	Mem            int     `json:"mem"`
	Id             string  `json:"id"`
	Node           string  `json:"node"`
	Disk           int64   `json:"disk"`
	Uptime         int     `json:"uptime"`
	Maxdisk        int64   `json:"maxdisk"`
	Level          string  `json:"level"`
}

const NodesPath string = "nodes"

func (client *Client) GetNodes() ([]Node, error) {
	var node []Node
	request, err := http.NewRequest(
		"GET",
		client.Host+ApiPath+NodesPath,
		nil,
	)
	if err != nil {
		return node, err
	}

	cookie := &http.Cookie{
		Name:  "PVEAuthCookie",
		Value: client.Ticket.Data.Ticket,
	}
	request.AddCookie(cookie)

	request.Header.Set(
		"CSRFPreventionToken",
		client.Ticket.Data.CSRFPreventionToken,
	)

	response, err := client.HTTPClient.Do(request)
	if err != nil {
		return node, err
	}

	if response.StatusCode != http.StatusOK {
		return node, errors.New(response.Status)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return node, err
	}

	nodeModel := NodeModel{}
	err = json.Unmarshal(body, &nodeModel)
	if err != nil {
		return node, err
	}

	node = nodeModel.Data

	err = response.Body.Close()
	if err != nil {
		return node, nil
	}

	return node, nil
}
