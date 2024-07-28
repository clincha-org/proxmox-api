package proxmox

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const NetworkPath = "/network"

func (client *Client) GetNetworks(node *Node) ([]Network, error) {
	request, err := http.NewRequest(
		"GET",
		client.Host+ApiPath+NodesPath+"/"+node.Node+NetworkPath,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("GetNetworks-build-request: %w", err)
	}

	request.AddCookie(&http.Cookie{Name: "PVEAuthCookie", Value: client.Ticket.Data.Ticket})
	request.Header.Set("CSRFPreventionToken", client.Ticket.Data.CSRFPreventionToken)

	response, err := client.HTTPClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("GetNetworks-do-request: %w", err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("GetNetworks-read-response: %w", err)
	}

	err = response.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("GetNetworks-close-response: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GetNetworks-status-error: %s %s", response.Status, body)
	}

	networkModel := NetworksResponse{}
	err = json.Unmarshal(body, &networkModel)
	if err != nil {
		return nil, fmt.Errorf("GetNetworks-unmarshal-response: %w", err)
	}

	return networkModel.Data, nil
}

func (client *Client) GetNetwork(node *Node, networkName string) (Network, error) {
	var network Network

	request, err := http.NewRequest(
		"GET",
		client.Host+ApiPath+NodesPath+"/"+node.Node+NetworkPath+"/"+networkName,
		nil,
	)
	if err != nil {
		return network, fmt.Errorf("GetNetwork-build-request: %w", err)
	}

	request.AddCookie(&http.Cookie{Name: "PVEAuthCookie", Value: client.Ticket.Data.Ticket})
	request.Header.Set("CSRFPreventionToken", client.Ticket.Data.CSRFPreventionToken)

	response, err := client.HTTPClient.Do(request)
	if err != nil {
		return network, fmt.Errorf("GetNetwork-do-request: %w", err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return network, fmt.Errorf("GetNetwork-read-response: %w", err)
	}

	err = response.Body.Close()
	if err != nil {
		return network, fmt.Errorf("GetNetwork-close-response: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		return network, fmt.Errorf("GetNetwork-status-error: %s %s", response.Status, body)
	}

	networkModel := NetworkResponse{}
	err = json.Unmarshal(body, &networkModel)
	if err != nil {
		return network, fmt.Errorf("GetNetwork-unmarshal-response: %w", err)
	}

	network = networkModel.Data
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

	request, err := http.NewRequest(
		"POST",
		client.Host+ApiPath+NodesPath+"/"+node.Node+NetworkPath+"/",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return Network{}, fmt.Errorf("CreateNetwork-build-request: %w", err)
	}

	request.AddCookie(&http.Cookie{Name: "PVEAuthCookie", Value: client.Ticket.Data.Ticket})
	request.Header.Set("CSRFPreventionToken", client.Ticket.Data.CSRFPreventionToken)
	request.Header.Set("Content-Type", "application/json")

	response, err := client.HTTPClient.Do(request)
	if err != nil {
		return Network{}, fmt.Errorf("CreateNetwork-do-request: %w", err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return Network{}, fmt.Errorf("CreateNetwork-read-response: %w", err)
	}

	err = response.Body.Close()
	if err != nil {
		return Network{}, fmt.Errorf("CreateNetwork-close-response: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		return Network{}, fmt.Errorf("CreateNetwork-status-error: %s %s", response.Status, body)
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

	request, err := http.NewRequest(
		"PUT",
		client.Host+ApiPath+NodesPath+"/"+node.Node+NetworkPath+"/"+networkRequest.Interface,
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return Network{}, fmt.Errorf("UpdateNetwork-build-request: %w", err)
	}

	request.AddCookie(&http.Cookie{Name: "PVEAuthCookie", Value: client.Ticket.Data.Ticket})
	request.Header.Set("CSRFPreventionToken", client.Ticket.Data.CSRFPreventionToken)
	request.Header.Set("Content-Type", "application/json")

	response, err := client.HTTPClient.Do(request)
	if err != nil {
		return Network{}, fmt.Errorf("UpdateNetwork-do-request: %w", err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return Network{}, fmt.Errorf("UpdateNetwork-read-response: %w", err)
	}

	err = response.Body.Close()
	if err != nil {
		return Network{}, fmt.Errorf("UpdateNetwork-close-response: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		return Network{}, fmt.Errorf("UpdateNetwork-status-error: %s %s", response.Status, body)
	}

	err = client.ReloadNetwork(node)
	if err != nil {
		return Network{}, fmt.Errorf("UpdateNetwork-reload-network Node - %s: %w", node.Node, err)
	}

	return client.GetNetwork(node, networkRequest.Interface)
}

func (client *Client) DeleteNetwork(node *Node, network string) error {
	request, err := http.NewRequest(
		"DELETE",
		client.Host+ApiPath+NodesPath+"/"+node.Node+NetworkPath+"/"+network,
		nil,
	)
	if err != nil {
		return fmt.Errorf("DeleteNetwork-build-request: %w", err)
	}

	request.AddCookie(&http.Cookie{Name: "PVEAuthCookie", Value: client.Ticket.Data.Ticket})
	request.Header.Set("CSRFPreventionToken", client.Ticket.Data.CSRFPreventionToken)

	response, err := client.HTTPClient.Do(request)
	if err != nil {
		return fmt.Errorf("DeleteNetwork-do-request: %w", err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("DeleteNetwork-read-response: %w", err)
	}

	err = response.Body.Close()
	if err != nil {
		return fmt.Errorf("DeleteNetwork-close-response: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("DeleteNetwork-status-error: %s %s", response.Status, body)
	}

	err = client.ReloadNetwork(node)
	if err != nil {
		return fmt.Errorf("DeleteNetwork-reload-network Node - %s: %w", node.Node, err)
	}

	return nil
}

func (client *Client) ReloadNetwork(node *Node) error {
	request, err := http.NewRequest(
		"PUT",
		client.Host+ApiPath+NodesPath+"/"+node.Node+NetworkPath+"/",
		nil,
	)
	if err != nil {
		return fmt.Errorf("ReloadNetwork-build-request: %w", err)
	}

	request.AddCookie(&http.Cookie{Name: "PVEAuthCookie", Value: client.Ticket.Data.Ticket})
	request.Header.Set("CSRFPreventionToken", client.Ticket.Data.CSRFPreventionToken)

	response, err := client.HTTPClient.Do(request)
	if err != nil {
		return fmt.Errorf("ReloadNetwork-do-request: %w", err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("ReloadNetwork-read-response: %w", err)
	}

	err = response.Body.Close()
	if err != nil {
		return fmt.Errorf("ReloadNetwork-close-response: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("ReloadNetwork-status-error: %s %s", response.Status, body)
	}

	// Give the daemon some time to reload the network
	// https://github.com/clincha-org/proxmox-api/issues/1
	time.Sleep(1 * time.Second)

	return nil
}
