package proxmox

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
)

const VirtualMachinePath = "/qemu"

func (client *Client) GetVM(node *Node, vmID int) (VirtualMachineConfig, error) {
	request, err := http.NewRequest(
		"GET",
		client.Host+ApiPath+NodesPath+"/"+node.Node+VirtualMachinePath+"/"+strconv.Itoa(vmID)+"/config",
		nil,
	)
	if err != nil {
		return VirtualMachineConfig{}, fmt.Errorf("GetVM-build-request: %w", err)
	}

	request.AddCookie(&http.Cookie{Name: "PVEAuthCookie", Value: client.Ticket.Data.Ticket})
	request.Header.Set("CSRFPreventionToken", client.Ticket.Data.CSRFPreventionToken)

	response, err := client.HTTPClient.Do(request)
	if err != nil {
		return VirtualMachineConfig{}, fmt.Errorf("GetVM-do-request: %w", err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return VirtualMachineConfig{}, fmt.Errorf("GetVM-read-response: %w", err)
	}

	err = response.Body.Close()
	if err != nil {
		return VirtualMachineConfig{}, fmt.Errorf("GetVM-close-response: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		return VirtualMachineConfig{}, fmt.Errorf("GetVM-status-error: %s %s", response.Status, body)
	}

	slog.Debug("api-response", "method", "GetVM", "node", node.Node, "response", string(body))

	vmModel := VirtualMachineConfigResponse{}
	err = json.Unmarshal(body, &vmModel)
	if err != nil {
		return VirtualMachineConfig{}, fmt.Errorf("GetVM-unmarshal-response: %w", err)
	}

	return vmModel.Data, nil
}

func (client *Client) GetVMs(node *Node) ([]VirtualMachine, error) {
	request, err := http.NewRequest(
		"GET",
		client.Host+ApiPath+NodesPath+"/"+node.Node+VirtualMachinePath,
		nil,
	)
	if err != nil {
		return []VirtualMachine{}, fmt.Errorf("GetVMs-build-request: %w", err)
	}

	request.AddCookie(&http.Cookie{Name: "PVEAuthCookie", Value: client.Ticket.Data.Ticket})
	request.Header.Set("CSRFPreventionToken", client.Ticket.Data.CSRFPreventionToken)

	response, err := client.HTTPClient.Do(request)
	if err != nil {
		return []VirtualMachine{}, fmt.Errorf("GetVMs-do-request: %w", err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return []VirtualMachine{}, fmt.Errorf("GetVMs-read-response: %w", err)
	}

	err = response.Body.Close()
	if err != nil {
		return []VirtualMachine{}, fmt.Errorf("GetVMs-close-response: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		return []VirtualMachine{}, fmt.Errorf("GetVMs-status-error: %s %s", response.Status, body)
	}

	slog.Debug("api-response", "method", "GetVMs", "node", node.Node, "response", string(body))

	vmModel := VirtualMachinesResponse{}
	err = json.Unmarshal(body, &vmModel)
	if err != nil {
		return []VirtualMachine{}, fmt.Errorf("GetVMs-unmarshal-response: %w", err)
	}

	return vmModel.Data, nil
}

func (client *Client) CreateVM(node *Node, vmRequest *VirtualMachineRequest, start bool) (VirtualMachineConfig, error) {
	requestBody, err := json.Marshal(vmRequest)
	if err != nil {
		return VirtualMachineConfig{}, fmt.Errorf("CreateVM-marshal-request: %w", err)
	}

	request, err := http.NewRequest(
		"POST",
		client.Host+ApiPath+NodesPath+"/"+node.Node+VirtualMachinePath,
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return VirtualMachineConfig{}, fmt.Errorf("CreateVM-build-request: %w", err)
	}

	request.AddCookie(&http.Cookie{Name: "PVEAuthCookie", Value: client.Ticket.Data.Ticket})
	request.Header.Set("CSRFPreventionToken", client.Ticket.Data.CSRFPreventionToken)
	request.Header.Set("Content-Type", "application/json")

	response, err := client.HTTPClient.Do(request)
	if err != nil {
		return VirtualMachineConfig{}, fmt.Errorf("CreateVM-do-request: %w", err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return VirtualMachineConfig{}, fmt.Errorf("CreateVM-read-response: %w", err)
	}

	err = response.Body.Close()
	if err != nil {
		return VirtualMachineConfig{}, fmt.Errorf("CreateVM-close-response: %w", err)
	}

	slog.Debug("api-response", "method", "CreateVM", "node", node.Node, "status", response.Status, "response", string(body))

	if response.StatusCode != http.StatusOK {
		return VirtualMachineConfig{}, fmt.Errorf("CreateVM-status-error: %s %s", response.Status, body)
	}

	if start {
		err = client.StartVm(node, vmRequest.VMID)
		if err != nil {
			return VirtualMachineConfig{}, fmt.Errorf("CreateVM-start-vm: %w", err)
		}
	}

	return client.GetVM(node, vmRequest.VMID)
}

func (client *Client) DeleteVM(node *Node, vmid int) error {

	// Check if the VM is still running
	vmStatus, err := client.GetVMStatus(node, vmid)
	if err != nil {
		return fmt.Errorf("DeleteVM-get-vm-status: %w", err)
	}

	if vmStatus.Status != "stopped" {
		// Stop the VM
		err = client.StopVM(node, vmid)
		if err != nil {
			return fmt.Errorf("DeleteVM-stop-vm: %w", err)
		}

		// Poll the VM status until it is stopped
		for ok := true; ok; ok = vmStatus.Status != "stopped" {
			vmStatus, err = client.GetVMStatus(node, vmid)
			if err != nil {
				return fmt.Errorf("DeleteVM-get-vm-status-loop: %w", err)
			}
		}
	}

	// Once the VM is stopped, delete it
	request, err := http.NewRequest(
		"DELETE",
		client.Host+ApiPath+NodesPath+"/"+node.Node+VirtualMachinePath+"/"+strconv.Itoa(vmid),
		nil,
	)
	if err != nil {
		return fmt.Errorf("DeleteVM-build-request: %w", err)
	}

	request.AddCookie(&http.Cookie{Name: "PVEAuthCookie", Value: client.Ticket.Data.Ticket})
	request.Header.Set("CSRFPreventionToken", client.Ticket.Data.CSRFPreventionToken)

	response, err := client.HTTPClient.Do(request)
	if err != nil {
		return fmt.Errorf("DeleteVM-do-request: %w", err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("DeleteVM-read-response: %w", err)
	}

	err = response.Body.Close()
	if err != nil {
		return fmt.Errorf("DeleteVM-close-response: %w", err)
	}

	slog.Debug("api-response", "method", "DeleteVM", "node", node.Node, "status", response.Status, "response", string(body))

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("DeleteVM-status-error: %s %s", response.Status, body)
	}

	return nil
}
