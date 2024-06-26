package proxmox

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"reflect"
	"slices"
	"strings"
	"time"
)

const NetworkPath = "/network"

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

	networkModel := NetworksResponse{}
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

	slog.Debug(fmt.Sprintf("Response from GetNetwork endpoint was: %v", string(body)))

	networkModel := NetworkResponse{}
	err = json.Unmarshal(body, &networkModel)
	if err != nil {
		return network, fmt.Errorf("unable to unmarshall JSON, error was: %v", err)
	}

	// Check to make sure there aren't any fields that are missing in the struct. Warn if there are missing fields.
	// More information in: https://github.com/clincha-org/proxmox-api/issues/5
	var networkJSON map[string]map[string]interface{}

	err = json.Unmarshal(body, &networkJSON)
	if err != nil {
		slog.Error(fmt.Sprintf("Unable to unmarshall body for missing field check, error was %v", err))
	}

	networkModelStruct := reflect.ValueOf(&network).Elem()
	var NetworkStructFields []string
	for i := 0; i < networkModelStruct.NumField(); i++ {
		JSONFieldName := networkModelStruct.Type().Field(i).Tag.Get("json") // Get the JSON representation of the field
		CommaIndex := strings.Index(JSONFieldName, ",")                     // Get the index of the comma (,omitempty)
		JSONFieldName = JSONFieldName[:CommaIndex]                          // Remove everything after the comma
		NetworkStructFields = append(NetworkStructFields, JSONFieldName)
	}

	for name, _ := range networkJSON["data"] {
		if !slices.Contains(NetworkStructFields, name) {
			slog.Warn(fmt.Sprintf("field %q returned by Proxmox API but field does not exist in NetworkModel struct. Please report this to the developers. See: https://github.com/clincha-org/proxmox-api/issues/5", name))
		}
	}

	network = networkModel.Data
	network.Interface = networkName

	return network, nil
}

func (client *Client) CreateNetwork(node *Node, networkRequest *NetworkRequest) (Network, error) {
	var url = client.Host + ApiPath + NodesPath + "/" + node.Node + NetworkPath + "/"

	jsonData, err := json.Marshal(&networkRequest)
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

	err = client.ReloadNetwork(node)
	if err != nil {
		return Network{}, fmt.Errorf("network reload failed, error was: %v", err)
	}

	return client.GetNetwork(node, networkRequest.Interface)
}

func (client *Client) UpdateNetwork(node *Node, networkRequest *NetworkRequest) (Network, error) {
	var url = client.Host + ApiPath + NodesPath + "/" + node.Node + NetworkPath + "/" + networkRequest.Interface

	jsonData, err := json.Marshal(&networkRequest)
	if err != nil {
		return Network{}, fmt.Errorf("unable to marshall JSON, error was: %v", err)
	}

	request, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return Network{}, fmt.Errorf("unable to create update network HTTP PUT request. Error was %v", err)
	}

	request.AddCookie(&http.Cookie{Name: "PVEAuthCookie", Value: client.Ticket.Data.Ticket})
	request.Header.Set("CSRFPreventionToken", client.Ticket.Data.CSRFPreventionToken)
	request.Header.Set("Content-Type", "application/json")

	response, err := client.HTTPClient.Do(request)
	if err != nil {
		return Network{}, fmt.Errorf("error recieved when making request to update network. Error was %v", err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return Network{}, fmt.Errorf("error reading response body. Error was %v", err)
	}

	if response.StatusCode != http.StatusOK {
		return Network{}, fmt.Errorf("network update failed. Status returned %v Body of response was %v", response.Status, string(body))
	}

	err = client.ReloadNetwork(node)
	if err != nil {
		return Network{}, fmt.Errorf("network reload failed, error was: %v", err)
	}

	return client.GetNetwork(node, networkRequest.Interface)
}

func (client *Client) DeleteNetwork(node *Node, network string) error {
	var url = client.Host + ApiPath + NodesPath + "/" + node.Node + NetworkPath + "/" + network

	request, err := http.NewRequest("DELETE", url, nil)
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

	err = client.ReloadNetwork(node)
	if err != nil {
		return fmt.Errorf("network reload failed, error was: %v", err)
	}

	return nil
}

func (client *Client) ReloadNetwork(node *Node) error {
	var url = client.Host + ApiPath + NodesPath + "/" + node.Node + NetworkPath + "/"

	// Now we need to call the PUT endpoint of the node's network to reload the network configuration
	request, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return fmt.Errorf("unable to create new network HTTP PUT request, error was %v", err)
	}

	request.AddCookie(&http.Cookie{Name: "PVEAuthCookie", Value: client.Ticket.Data.Ticket})
	request.Header.Set("CSRFPreventionToken", client.Ticket.Data.CSRFPreventionToken)

	response, err := client.HTTPClient.Do(request)
	if err != nil {
		return fmt.Errorf("error recieved when making request to reload network, error was %v", err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("error reading response body. Error was %v", err)
	}

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("network reload failed: %v Body of response was %v", response.Status, string(body))
	}

	//Allow the network daemon time to reload the configuration
	time.Sleep(1 * time.Second)

	return nil
}
