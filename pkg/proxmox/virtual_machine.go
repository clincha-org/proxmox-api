package proxmox

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/clincha-org/proxmox-api/pkg/ide"
	"io"
	"log/slog"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

const VirtualMachinePath = "/qemu"

func (client *Client) GetVM(node string, id int64) (VirtualMachine, error) {
	request, err := http.NewRequest(
		"GET",
		client.Host+ApiPath+NodesPath+"/"+node+VirtualMachinePath+"/"+strconv.FormatInt(id, 10)+"/config",
		nil,
	)
	if err != nil {
		return VirtualMachine{}, fmt.Errorf("GetVM-build-request: %w", err)
	}

	request.AddCookie(&http.Cookie{Name: "PVEAuthCookie", Value: client.Ticket.Data.Ticket})
	request.Header.Set("CSRFPreventionToken", client.Ticket.Data.CSRFPreventionToken)

	response, err := client.HTTPClient.Do(request)
	if err != nil {
		return VirtualMachine{}, fmt.Errorf("GetVM-do-request: %w", err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return VirtualMachine{}, fmt.Errorf("GetVM-read-response: %w", err)
	}

	err = response.Body.Close()
	if err != nil {
		return VirtualMachine{}, fmt.Errorf("GetVM-close-response: %w", err)
	}

	slog.Debug("api-response", "method", "GetVM", "node", node, "status", response.Status, "response", string(body))

	if response.StatusCode != http.StatusOK {
		return VirtualMachine{}, fmt.Errorf("GetVM-status-error: %s %s", response.Status, body)
	}

	vmModel := VirtualMachineConfigResponse{}

	// The API returns numbers with and without quotes, so we quote all numbers to make it easier to unmarshal
	re := regexp.MustCompile(`(":\s*)([\d.]+)(\s*[,}])`)
	body = re.ReplaceAll(body, []byte(`$1"$2"$3`))

	slog.Debug("api-response-quoted", "method", "GetVM", "node", node, "response", string(body))

	err = json.Unmarshal(body, &vmModel)
	if err != nil {
		return VirtualMachine{}, fmt.Errorf("GetVM-unmarshal-response: %w", err)
	}

	vm := VirtualMachine{
		ID:           id,
		Net1:         &vmModel.Data.Net1,
		SCSIHardware: &vmModel.Data.Scsihw,
		Cores:        vmModel.Data.Cores,
		Memory:       vmModel.Data.Memory,
	}

	if vmModel.Data.IDE0 != nil || vmModel.Data.IDE1 != nil || vmModel.Data.IDE2 != nil || vmModel.Data.IDE3 != nil {
		var IdeDevices []ide.InternalDataStorage
		for index, IDEDeviceString := range []*string{vmModel.Data.IDE0, vmModel.Data.IDE1, vmModel.Data.IDE2, vmModel.Data.IDE3} {
			if IDEDeviceString == nil {
				continue
			}

			device := ide.InternalDataStorage{}
			err := ide.Unmarshal(int64(index), *IDEDeviceString, &device)
			if err != nil {
				return VirtualMachine{}, err
			}
			IdeDevices = append(IdeDevices, device)
		}
		vm.IDEDevices = &IdeDevices
	}
	return vm, nil
}

func (client *Client) GetVMs(node string) ([]VirtualMachineListItem, error) {
	request, err := http.NewRequest(
		"GET",
		client.Host+ApiPath+NodesPath+"/"+node+VirtualMachinePath,
		nil,
	)
	if err != nil {
		return []VirtualMachineListItem{}, fmt.Errorf("GetVMs-build-request: %w", err)
	}

	request.AddCookie(&http.Cookie{Name: "PVEAuthCookie", Value: client.Ticket.Data.Ticket})
	request.Header.Set("CSRFPreventionToken", client.Ticket.Data.CSRFPreventionToken)

	response, err := client.HTTPClient.Do(request)
	if err != nil {
		return []VirtualMachineListItem{}, fmt.Errorf("GetVMs-do-request: %w", err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return []VirtualMachineListItem{}, fmt.Errorf("GetVMs-read-response: %w", err)
	}

	err = response.Body.Close()
	if err != nil {
		return []VirtualMachineListItem{}, fmt.Errorf("GetVMs-close-response: %w", err)
	}

	slog.Debug("api-response", "method", "GetVMs", "node", node, "status", response.Status, "response", string(body))

	if response.StatusCode != http.StatusOK {
		return []VirtualMachineListItem{}, fmt.Errorf("GetVMs-status-error: %s %s", response.Status, body)
	}

	vmModel := VirtualMachinesResponse{}
	err = json.Unmarshal(body, &vmModel)
	if err != nil {
		return []VirtualMachineListItem{}, fmt.Errorf("GetVMs-unmarshal-response: %w", err)
	}

	return vmModel.Data, nil
}

func (client *Client) CreateVM(node string, vm *VirtualMachine, start bool) (VirtualMachine, error) {
	vmRequest := VirtualMachineRequest{
		ID:           vm.ID,
		SCSI1:        vm.SCSI1,
		Net1:         vm.Net1,
		SCSIHardware: vm.SCSIHardware,
		Cores:        vm.Cores,
		Memory:       vm.Memory,
	}

	if vm.IDEDevices != nil {
		slog.Debug("ide-devices", "method", "CreateVM", "devices", vm.IDEDevices)
		if len(*vm.IDEDevices) > 4 {
			return VirtualMachine{}, fmt.Errorf("CreateVM-invalid-number-of-ide-devices: %d. Proxmox only allows 4 IDE devices", len(*vm.IDEDevices))
		}

		for _, ideDevice := range *vm.IDEDevices {
			slog.Debug("ide-device", "method", "CreateVM", "device", ideDevice)
			if ideDevice.ID > 3 || ideDevice.ID < 0 {
				return VirtualMachine{}, fmt.Errorf("CreateVM-invalid-ide-device: %d", ideDevice.ID)
			}

			marshal, err := ide.Marshal(&ideDevice)
			if err != nil {
				return VirtualMachine{}, err
			}

			switch ideDevice.ID {
			case 0:
				vmRequest.IDE0 = &marshal
			case 1:
				vmRequest.IDE1 = &marshal
			case 2:
				vmRequest.IDE2 = &marshal
			case 3:
				vmRequest.IDE3 = &marshal
			}
		}
	}

	slog.Debug("api-request", "method", "CreateVM", "node", node, "request", vmRequest)
	requestBody, err := json.Marshal(vmRequest)
	if err != nil {
		return VirtualMachine{}, fmt.Errorf("CreateVM-marshal-request: %w", err)
	}

	request, err := http.NewRequest(
		"POST",
		client.Host+ApiPath+NodesPath+"/"+node+VirtualMachinePath,
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return VirtualMachine{}, fmt.Errorf("CreateVM-build-request: %w", err)
	}

	request.AddCookie(&http.Cookie{Name: "PVEAuthCookie", Value: client.Ticket.Data.Ticket})
	request.Header.Set("CSRFPreventionToken", client.Ticket.Data.CSRFPreventionToken)
	request.Header.Set("Content-Type", "application/json")

	response, err := client.HTTPClient.Do(request)
	if err != nil {
		return VirtualMachine{}, fmt.Errorf("CreateVM-do-request: %w", err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return VirtualMachine{}, fmt.Errorf("CreateVM-read-response: %w", err)
	}

	err = response.Body.Close()
	if err != nil {
		return VirtualMachine{}, fmt.Errorf("CreateVM-close-response: %w", err)
	}

	slog.Debug("api-response", "method", "CreateVM", "node", node, "status", response.Status, "response", string(body))

	if response.StatusCode != http.StatusOK {
		return VirtualMachine{}, fmt.Errorf("CreateVM-status-error: %s %s", response.Status, body)
	}

	// Make sure the VM has finished configuring
	job := AsynchronousTaskResponse{}
	err = json.Unmarshal(body, &job)
	if err != nil {
		return VirtualMachine{}, fmt.Errorf("CreateVM-unmarshal-response: %w", err)
	}
	err = client.AwaitAsynchronousTask(node, job.ID)
	if err != nil {
		return VirtualMachine{}, fmt.Errorf("CreateVM-await-task: %w", err)
	}

	if start {
		err = client.StartVm(node, vm.ID)
		if err != nil {
			return VirtualMachine{}, fmt.Errorf("CreateVM-start-vm: %w", err)
		}
	}

	return client.GetVM(node, vm.ID)
}

func (client *Client) UpdateVM(node string, vm *VirtualMachine) (VirtualMachine, error) {
	vmRequest := VirtualMachineRequest{
		ID:           vm.ID,
		SCSI1:        vm.SCSI1,
		Net1:         vm.Net1,
		SCSIHardware: vm.SCSIHardware,
		Cores:        vm.Cores,
		Memory:       vm.Memory,
	}

	currentState, err := client.GetVM(node, vm.ID)
	if err != nil {
		return VirtualMachine{}, fmt.Errorf("UpdateVM-get-current-state: %w", err)
	}

	// If there are no IDE devices in the current state, then there is nothing to delete
	if currentState.IDEDevices != nil {
		// If there are no IDE devices in the new state, then we need to delete all IDE devices
		if vm.IDEDevices == nil {
			var deleteDevices []string
			for _, ideDevice := range *currentState.IDEDevices {
				deleteDevices = append(deleteDevices, "ide"+strconv.FormatInt(ideDevice.ID, 10))
			}
			ideDeviceStrings := strings.Join(deleteDevices, ",")
			vmRequest.Delete = &ideDeviceStrings
		} else {
			// Delete any IDE devices that are not in the new state
			var deleteDevices []string
			for _, ideDevice := range *currentState.IDEDevices {
				found := false
				for _, newIdeDevice := range *vm.IDEDevices {
					if ideDevice.ID == newIdeDevice.ID {
						found = true
						break
					}
				}
				if !found {
					deleteDevices = append(deleteDevices, "ide"+strconv.FormatInt(ideDevice.ID, 10))
				}
			}

			if len(deleteDevices) > 0 {
				ideDeviceStrings := strings.Join(deleteDevices, ",")
				vmRequest.Delete = &ideDeviceStrings
			}
		}
	}

	if vm.IDEDevices != nil {
		for _, ideDevice := range *vm.IDEDevices {

			slog.Debug("ide-device", "method", "UpdateVM", "device", ideDevice)

			if ideDevice.ID > 3 || ideDevice.ID < 0 {
				return VirtualMachine{}, fmt.Errorf("CreateVM-invalid-ide-device: %d", ideDevice.ID)
			}

			marshal, err := ide.Marshal(&ideDevice)
			if err != nil {
				return VirtualMachine{}, err
			}

			switch ideDevice.ID {
			case 0:
				vmRequest.IDE0 = &marshal
			case 1:
				vmRequest.IDE1 = &marshal
			case 2:
				vmRequest.IDE2 = &marshal
			case 3:
				vmRequest.IDE3 = &marshal
			}
		}

		if len(*vm.IDEDevices) > 4 {
			return VirtualMachine{}, fmt.Errorf("CreateVM-invalid-number-of-ide-devices: %d. Proxmox only allows 4 IDE devices", len(*vm.IDEDevices))
		}
	}

	requestBody, err := json.Marshal(vmRequest)
	if err != nil {
		return VirtualMachine{}, fmt.Errorf("UpdateVM-marshal-request: %w", err)
	}

	slog.Debug("api-request", "method", "UpdateVM", "node", node, "request", string(requestBody))

	request, err := http.NewRequest(
		"POST",
		client.Host+ApiPath+NodesPath+"/"+node+VirtualMachinePath+"/"+strconv.FormatInt(vmRequest.ID, 10)+"/config",
		bytes.NewBuffer(requestBody),
	)
	if err != nil {
		return VirtualMachine{}, fmt.Errorf("UpdateVM-build-request: %w", err)
	}

	request.AddCookie(&http.Cookie{Name: "PVEAuthCookie", Value: client.Ticket.Data.Ticket})
	request.Header.Set("CSRFPreventionToken", client.Ticket.Data.CSRFPreventionToken)
	request.Header.Set("Content-Type", "application/json")

	response, err := client.HTTPClient.Do(request)
	if err != nil {
		return VirtualMachine{}, fmt.Errorf("UpdateVM-do-request: %w", err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return VirtualMachine{}, fmt.Errorf("UpdateVM-read-response: %w", err)
	}

	err = response.Body.Close()
	if err != nil {
		return VirtualMachine{}, fmt.Errorf("UpdateVM-close-response: %w", err)
	}

	slog.Debug("api-response", "method", "UpdateVM", "node", node, "status", response.Status, "response", string(body))

	if response.StatusCode != http.StatusOK {
		return VirtualMachine{}, fmt.Errorf("UpdateVM-status-error: %s %s", response.Status, body)
	}

	// Make sure the VM has finished configuring
	job := AsynchronousTaskResponse{}
	err = json.Unmarshal(body, &job)
	if err != nil {
		return VirtualMachine{}, fmt.Errorf("UpdateVM-unmarshal-response: %w", err)
	}
	err = client.AwaitAsynchronousTask(node, job.ID)
	if err != nil {
		return VirtualMachine{}, fmt.Errorf("UpdateVM-await-task: %w", err)
	}

	return client.GetVM(node, vmRequest.ID)
}

func (client *Client) DeleteVM(node string, id int64) error {

	// Check if the VM is still running
	vmStatus, err := client.GetVMStatus(node, id)
	if err != nil {
		return fmt.Errorf("DeleteVM-get-vm-status: %w", err)
	}

	if vmStatus.Status != "stopped" {
		// Stop the VM
		err = client.StopVM(node, id)
		if err != nil {
			return fmt.Errorf("DeleteVM-stop-vm: %w", err)
		}
	}

	// Once the VM is stopped, delete it
	request, err := http.NewRequest(
		"DELETE",
		client.Host+ApiPath+NodesPath+"/"+node+VirtualMachinePath+"/"+strconv.FormatInt(id, 10),
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

	slog.Debug("api-response", "method", "DeleteVM", "node", node, "status", response.Status, "response", string(body))

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("DeleteVM-status-error: %s %s", response.Status, body)
	}

	// Make sure the VM has finished configuring
	job := AsynchronousTaskResponse{}
	err = json.Unmarshal(body, &job)
	if err != nil {
		return fmt.Errorf("DeleteVM-unmarshal-response: %w", err)
	}
	err = client.AwaitAsynchronousTask(node, job.ID)
	if err != nil {
		return fmt.Errorf("DeleteVM-await-task: %w", err)
	}

	return nil
}
