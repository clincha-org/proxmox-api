package proxmox

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

const NetworkPath = "/network"

type NetworkModel struct {
	Data []Network `json:"data"`
}

type Network struct {
	Gateway     string   `json:"gateway"`
	Type        string   `json:"type"`
	Autostart   int      `json:"autostart"`
	Families    []string `json:"families"`
	Method6     string   `json:"method6"`
	Iface       string   `json:"iface"`
	BridgeFd    string   `json:"bridge_fd"`
	Netmask     string   `json:"netmask"`
	Priority    int      `json:"priority"`
	Active      int      `json:"active"`
	Method      string   `json:"method"`
	BridgeStp   string   `json:"bridge_stp"`
	Address     string   `json:"address"`
	Cidr        string   `json:"cidr"`
	BridgePorts string   `json:"bridge_ports"`
}

func (client *Client) GetNetworks(node *Node) ([]Network, error) {
	var network []Network
	request, err := http.NewRequest(
		"GET",
		client.Host+ApiPath+NodesPath+"/"+node.Node+NetworkPath,
		nil,
	)
	if err != nil {
		return network, err
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
