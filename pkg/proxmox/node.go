package proxmox

import (
	"encoding/json"
	"fmt"
)

const NodesPath string = "nodes"

func (client *Client) GetNodes() ([]Node, error) {
	path := client.Host + ApiPath + NodesPath
	body, err := client.MakeRESTRequest("GET", path, nil)
	if err != nil {
		return []Node{}, fmt.Errorf("GetNodes-make-request: %w", err)
	}

	nodeModel := NodeResponse{}
	err = json.Unmarshal(body, &nodeModel)
	if err != nil {
		return []Node{}, fmt.Errorf("GetNodes-unmarshal-response: %w", err)
	}

	return nodeModel.Data, nil
}
