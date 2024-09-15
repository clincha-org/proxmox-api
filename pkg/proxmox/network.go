package proxmox

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

const NetworkPath = "/network"

func (client *Client) GetNetworks(node *Node) ([]Network, error) {
	path := client.Host + ApiPath + NodesPath + "/" + node.Node + NetworkPath
	body, err := client.MakeRESTRequest("GET", path, nil)
	if err != nil {
		return []Network{}, fmt.Errorf("GetNetworks-make-request: %w", err)
	}

	networkModel := NetworksResponse{}
	err = json.Unmarshal(body, &networkModel)
	if err != nil {
		return nil, fmt.Errorf("GetNetworks-unmarshal-response: %w", err)
	}

	return networkModel.Data, nil
}

func (client *Client) GetNetwork(node *Node, networkName string) (Network, error) {
	path := client.Host + ApiPath + NodesPath + "/" + node.Node + NetworkPath + "/" + networkName
	body, err := client.MakeRESTRequest("GET", path, nil)
	if err != nil {
		return Network{}, fmt.Errorf("GetNetwork-make-request: %w", err)
	}

	networkModel := NetworkResponse{}
	err = json.Unmarshal(body, &networkModel)
	if err != nil {
		return Network{}, fmt.Errorf("GetNetwork-unmarshal-response: %w", err)
	}

	network := networkModel.Data
	network.Interface = networkName

	// Convert the network CIDR into subnet mask format
	if network.Netmask != nil {
		network.Netmask, err = ConvertCIDRToNetmask(network.Netmask)
		if err != nil {
			return network, fmt.Errorf("GetNetwork-convert-cidr-to-netmask - CIDR - %s: %w", *network.Netmask, err)
		}
	}

	// Remove the newline at the end of the comment
	if network.Comments != nil && *network.Comments != "" {
		trimmedString := strings.Trim(*network.Comments, "\n")
		network.Comments = &trimmedString
	}

	return network, nil
}

func (client *Client) CreateNetwork(node *Node, networkRequest *NetworkRequest) (Network, error) {
	jsonData, err := json.Marshal(&networkRequest)
	if err != nil {
		return Network{}, fmt.Errorf("CreateNetwork-marshal-request: %w", err)
	}

	path := client.Host + ApiPath + NodesPath + "/" + node.Node + NetworkPath + "/"
	_, err = client.MakeRESTRequest("POST", path, bytes.NewBuffer(jsonData))
	if err != nil {
		return Network{}, fmt.Errorf("CreateNetwork-make-request: %w", err)
	}

	err = client.ReloadNetwork(node)
	if err != nil {
		return Network{}, fmt.Errorf("CreateNetwork-reload-network Node - %s: %w", node.Node, err)
	}

	return client.GetNetwork(node, networkRequest.Interface)
}

func (client *Client) UpdateNetwork(node *Node, networkRequest *NetworkRequest) (Network, error) {
	jsonData, err := json.Marshal(&networkRequest)
	if err != nil {
		return Network{}, fmt.Errorf("UpdateNetwork-marshal-request: %w", err)
	}

	path := client.Host + ApiPath + NodesPath + "/" + node.Node + NetworkPath + "/" + networkRequest.Interface
	_, err = client.MakeRESTRequest("PUT", path, bytes.NewBuffer(jsonData))
	if err != nil {
		return Network{}, fmt.Errorf("UpdateNetwork-make-request: %w", err)
	}

	err = client.ReloadNetwork(node)
	if err != nil {
		return Network{}, fmt.Errorf("UpdateNetwork-reload-network Node - %s: %w", node.Node, err)
	}

	return client.GetNetwork(node, networkRequest.Interface)
}

func (client *Client) DeleteNetwork(node *Node, network string) error {
	path := client.Host + ApiPath + NodesPath + "/" + node.Node + NetworkPath + "/" + network
	_, err := client.MakeRESTRequest("DELETE", path, nil)
	if err != nil {
		return fmt.Errorf("DeleteNetwork-make-request: %w", err)
	}

	err = client.ReloadNetwork(node)
	if err != nil {
		return fmt.Errorf("DeleteNetwork-reload-network Node - %s: %w", node.Node, err)
	}

	return nil
}

func (client *Client) ReloadNetwork(node *Node) error {
	path := client.Host + ApiPath + NodesPath + "/" + node.Node + NetworkPath + "/"
	body, err := client.MakeRESTRequest("PUT", path, nil)
	if err != nil {
		return fmt.Errorf("ReloadNetwork-make-request: %w", err)
	}

	job := AsynchronousTaskResponse{}
	err = json.Unmarshal(body, &job)
	if err != nil {
		return fmt.Errorf("ReloadNetwork-unmarshal-response: %w", err)
	}
	err = client.AwaitAsynchronousTask(node.Node, job.ID)
	if err != nil {
		return fmt.Errorf("ReloadNetwork-await-task: %w", err)
	}

	return nil
}
