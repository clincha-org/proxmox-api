package proxmox

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
)

func (client *Client) GetVMStatus(node string, vmid int) (VirtualMachineStatus, error) {
	request, err := http.NewRequest(
		"GET",
		client.Host+ApiPath+NodesPath+"/"+node+VirtualMachinePath+"/"+strconv.Itoa(vmid)+"/status/current",
		nil,
	)
	if err != nil {
		return VirtualMachineStatus{}, fmt.Errorf("GetVMStatus-build-request: %w", err)
	}

	request.AddCookie(&http.Cookie{Name: "PVEAuthCookie", Value: client.Ticket.Data.Ticket})
	request.Header.Set("CSRFPreventionToken", client.Ticket.Data.CSRFPreventionToken)

	response, err := client.HTTPClient.Do(request)
	if err != nil {
		return VirtualMachineStatus{}, fmt.Errorf("GetVMStatus-do-request: %w", err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return VirtualMachineStatus{}, fmt.Errorf("GetVMStatus-read-response: %w", err)
	}

	err = response.Body.Close()
	if err != nil {
		return VirtualMachineStatus{}, fmt.Errorf("GetVMStatus-close-response: %w", err)
	}

	slog.Debug("api-response", "method", "GetVMStatus", "node", node, "status", response.Status, "response", string(body))

	if response.StatusCode != http.StatusOK {
		return VirtualMachineStatus{}, fmt.Errorf("GetVMStatus-status-error: %s %s", response.Status, body)
	}

	vmStatus := VirtualMachineStatusResponse{}
	err = json.Unmarshal(body, &vmStatus)
	if err != nil {
		return VirtualMachineStatus{}, fmt.Errorf("GetVMStatus-unmarshal-response: %w", err)
	}

	return vmStatus.Data, nil
}

func (client *Client) StopVM(node string, vmid int) error {
	request, err := http.NewRequest(
		"POST",
		client.Host+ApiPath+NodesPath+"/"+node+VirtualMachinePath+"/"+strconv.Itoa(vmid)+"/status/stop",
		nil,
	)
	if err != nil {
		return fmt.Errorf("StopVM-build-request: %w", err)
	}

	request.AddCookie(&http.Cookie{Name: "PVEAuthCookie", Value: client.Ticket.Data.Ticket})
	request.Header.Set("CSRFPreventionToken", client.Ticket.Data.CSRFPreventionToken)

	response, err := client.HTTPClient.Do(request)
	if err != nil {
		return fmt.Errorf("StopVM-do-request: %w", err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("StopVM-read-response: %w", err)
	}

	err = response.Body.Close()
	if err != nil {
		return fmt.Errorf("StopVM-close-response: %w", err)
	}

	slog.Debug("api-response", "method", "StopVM", "node", node, "status", response.Status, "response", string(body))

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("StopVM-status-error: %s %s", response.Status, body)
	}

	return nil
}

func (client *Client) StartVm(node string, vmid int) error {
	request, err := http.NewRequest(
		"POST",
		client.Host+ApiPath+NodesPath+"/"+node+VirtualMachinePath+"/"+strconv.Itoa(vmid)+"/status/start",
		nil,
	)
	if err != nil {
		return fmt.Errorf("StartVM-build-request: %w", err)
	}

	request.AddCookie(&http.Cookie{Name: "PVEAuthCookie", Value: client.Ticket.Data.Ticket})
	request.Header.Set("CSRFPreventionToken", client.Ticket.Data.CSRFPreventionToken)

	response, err := client.HTTPClient.Do(request)
	if err != nil {
		return fmt.Errorf("StartVM-do-request: %w", err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("StartVM-read-response: %w", err)
	}

	err = response.Body.Close()
	if err != nil {
		return fmt.Errorf("StartVM-close-response: %w", err)
	}

	slog.Debug("api-response", "method", "StartVM", "node", node, "status", response.Status, "response", string(body))

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("StartVM-status-error: %s %s", response.Status, body)
	}

	return nil
}
