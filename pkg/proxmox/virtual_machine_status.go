package proxmox

import (
	"encoding/json"
	"fmt"
	"strconv"
)

func (client *Client) GetVMStatus(node string, id int64) (VirtualMachineStatus, error) {
	path := client.Host + ApiPath + NodesPath + "/" + node + VirtualMachinePath + "/" + strconv.FormatInt(id, 10) + "/status/current"
	body, err := client.MakeRESTRequest("GET", path, nil)
	if err != nil {
		return VirtualMachineStatus{}, fmt.Errorf("GetVMStatus-make-request: %w", err)
	}

	vmStatus := VirtualMachineStatusResponse{}
	err = json.Unmarshal(body, &vmStatus)
	if err != nil {
		return VirtualMachineStatus{}, fmt.Errorf("GetVMStatus-unmarshal-response: %w", err)
	}

	return vmStatus.Data, nil
}

func (client *Client) StartVm(node string, id int64) error {
	path := client.Host + ApiPath + NodesPath + "/" + node + VirtualMachinePath + "/" + strconv.FormatInt(id, 10) + "/status/start"
	body, err := client.MakeRESTRequest("POST", path, nil)
	if err != nil {
		return fmt.Errorf("StartVM-make-request: %w", err)
	}

	job := AsynchronousTaskResponse{}
	err = json.Unmarshal(body, &job)
	if err != nil {
		return fmt.Errorf("StartVM-unmarshal-response: %w", err)
	}
	err = client.AwaitAsynchronousTask(node, job.ID)

	return nil
}

func (client *Client) StopVM(node string, id int64) error {
	path := client.Host + ApiPath + NodesPath + "/" + node + VirtualMachinePath + "/" + strconv.FormatInt(id, 10) + "/status/stop"
	body, err := client.MakeRESTRequest("POST", path, nil)
	if err != nil {
		return fmt.Errorf("StopVM-make-request: %w", err)
	}

	job := AsynchronousTaskResponse{}
	err = json.Unmarshal(body, &job)
	if err != nil {
		return fmt.Errorf("StopVM-unmarshal-response: %w", err)
	}
	err = client.AwaitAsynchronousTask(node, job.ID)
	if err != nil {
		return fmt.Errorf("StopVM-await-task: %w", err)
	}

	return nil
}
