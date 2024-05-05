package proxmox

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const NetworkPath = "/network"

type NetworkModel struct {
	Data []Network `json:"data"`
}

type SingleNetworkModel struct {
	Data Network `json:"data"`
}

type Network struct {
	Gateway     string   `json:"gateway,omitempty"`
	Type        string   `json:"type,omitempty"`
	Autostart   int      `json:"autostart,omitempty"`
	Families    []string `json:"families,omitempty"`
	Method6     string   `json:"method6,omitempty"`
	Interface   string   `json:"iface,omitempty"`
	BridgeFd    string   `json:"bridge_fd,omitempty"`
	Netmask     string   `json:"netmask,omitempty"`
	Priority    int      `json:"priority,omitempty"`
	Active      int      `json:"active,omitempty"`
	Method      string   `json:"method,omitempty"`
	BridgeStp   string   `json:"bridge_stp,omitempty"`
	Address     string   `json:"address,omitempty"`
	Cidr        string   `json:"cidr,omitempty"`
	BridgePorts string   `json:"bridge_ports,omitempty"`
}

func (client *Client) GetNetworks(node *Node) ([]Network, error) {
	url := client.Host + ApiPath + NodesPath + "/" + node.Node + NetworkPath
	var network []Network

	request, err := http.NewRequest(
		"GET",
		url,
		nil,
	)
	if err != nil {
		return network, err
	}

	request.AddCookie(&http.Cookie{Name: "PVEAuthCookie", Value: client.Ticket.Data.Ticket})
	request.Header.Set("CSRFPreventionToken", client.Ticket.Data.CSRFPreventionToken)

	response, err := client.HTTPClient.Do(request)
	if err != nil {
		return network, err
	}

	if response.StatusCode != http.StatusOK {
		return network, errors.New(response.Status)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return network, err
	}

	networkModel := NetworkModel{}
	err = json.Unmarshal(body, &networkModel)
	if err != nil {
		return network, err
	}

	network = networkModel.Data

	return network, nil
}

func (client *Client) GetNetwork(node *Node, networkName string) (Network, error) {
	url := client.Host + ApiPath + NodesPath + "/" + node.Node + NetworkPath + "/" + networkName
	var network Network

	request, err := http.NewRequest(
		"GET",
		url,
		nil,
	)
	if err != nil {
		return network, err
	}

	request.AddCookie(&http.Cookie{Name: "PVEAuthCookie", Value: client.Ticket.Data.Ticket})
	request.Header.Set("CSRFPreventionToken", client.Ticket.Data.CSRFPreventionToken)

	response, err := client.HTTPClient.Do(request)
	if err != nil {
		return network, fmt.Errorf("unable to create new HTTP request. Error was %v", err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return network, fmt.Errorf("error reading response body. Error was %v", err)
	}

	if response.StatusCode != http.StatusOK {
		return network, fmt.Errorf("network creation failed. Status returned %v Body of response was %v", response.Status, string(body))
	}

	networkModel := SingleNetworkModel{}
	err = json.Unmarshal(body, &networkModel)
	if err != nil {
		return network, fmt.Errorf("unable to marshall JSON. Error was: %v", err)
	}

	network = networkModel.Data
	network.Interface = networkName

	return network, nil
}

func (client *Client) CreateNetwork(node *Node, network *Network) (Network, error) {
	var url = client.Host + ApiPath + NodesPath + "/" + node.Node + NetworkPath + "/"

	jsonData, err := json.Marshal(&network)
	if err != nil {
		return Network{}, fmt.Errorf("unable to marshall JSON. Error was: %v", err)
	}

	request, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return Network{}, fmt.Errorf("unable to create new network HTTP POST request. Error was %v", err)
	}

	request.AddCookie(&http.Cookie{Name: "PVEAuthCookie", Value: client.Ticket.Data.Ticket})
	request.Header.Set("CSRFPreventionToken", client.Ticket.Data.CSRFPreventionToken)
	request.Header.Set("Content-Type", "application/json")

	response, err := client.HTTPClient.Do(request)
	if err != nil {
		return Network{}, fmt.Errorf("error recieved when making request to create network. Error was %v", err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return Network{}, fmt.Errorf("error reading response body. Error was %v", err)
	}

	if response.StatusCode != http.StatusOK {
		return Network{}, fmt.Errorf("network creation failed. Status returned %v Body of response was %v", response.Status, string(body))
	}

	// Now we need to call the PUT endpoint of the node's network to reload the network configuration
	request, err = http.NewRequest("PUT", url, nil)
	if err != nil {
		return Network{}, fmt.Errorf("unable to create new network HTTP PUT request, error was %v", err)
	}

	request.AddCookie(&http.Cookie{Name: "PVEAuthCookie", Value: client.Ticket.Data.Ticket})
	request.Header.Set("CSRFPreventionToken", client.Ticket.Data.CSRFPreventionToken)

	response, err = client.HTTPClient.Do(request)
	if err != nil {
		return Network{}, fmt.Errorf("error recieved when making request to reload network, error was %v", err)
	}

	body, err = io.ReadAll(response.Body)
	if err != nil {
		return Network{}, fmt.Errorf("error reading response body. Error was %v", err)
	}

	if response.StatusCode != http.StatusOK {
		return Network{}, fmt.Errorf("network reload failed: %v Body of response was %v", response.Status, string(body))
	}

	return client.GetNetwork(node, network.Interface)
}

func (client *Client) DeleteNetwork(node *Node, network *Network) error {
	var url = client.Host + ApiPath + NodesPath + "/" + node.Node + NetworkPath + "/" + network.Interface

	request, err := http.NewRequest(
		"DELETE",
		url,
		nil,
	)
	if err != nil {
		return fmt.Errorf("unable to create request, error was %v", err)
	}

	request.AddCookie(&http.Cookie{Name: "PVEAuthCookie", Value: client.Ticket.Data.Ticket})
	request.Header.Set("CSRFPreventionToken", client.Ticket.Data.CSRFPreventionToken)

	response, err := client.HTTPClient.Do(request)
	if err != nil {
		return fmt.Errorf("request failed, error was %v", err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body, error was %v", err)
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("error recieved in http response when deleting netowrk, error was %v. Response body was %v", err, body)
	}

	return nil
}
