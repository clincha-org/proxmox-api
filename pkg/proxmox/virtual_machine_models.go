package proxmox

import (
	"github.com/clincha-org/proxmox-api/pkg/ide"
)

type VirtualMachine struct {
	ID           int64                      `json:"vmid"`
	IDEDevices   *[]ide.InternalDataStorage `json:"-"`
	SCSI1        *string                    `json:"scsi1"`
	Net1         *string                    `json:"net1"`
	SCSIHardware *string                    `json:"scsihw"`
	Cores        int64                      `json:"cores"`
	Memory       int64                      `json:"memory"`
}

type VirtualMachineRequest struct {
	ID           int64   `json:"vmid"`
	IDE0         *string `json:"ide0,omitempty"`
	IDE1         *string `json:"ide1,omitempty"`
	IDE2         *string `json:"ide2,omitempty"`
	IDE3         *string `json:"ide3,omitempty"`
	SCSI1        *string `json:"scsi1,omitempty"`
	Net1         *string `json:"net1,omitempty"`
	SCSIHardware *string `json:"scsihw,omitempty"`
	Cores        int64   `json:"cores,omitempty"`
	Memory       int64   `json:"memory,omitempty"`
	Delete       *string `json:"delete,omitempty"`
}

type VirtualMachineCloneRequest struct {
	NewID     int64 `json:"newid"`
	FullClone *bool `json:"full,omitempty"`
}

type VirtualMachinesResponse struct {
	Data []VirtualMachineListItem `json:"data"`
}
type VirtualMachineListItem struct {
	Status    string  `json:"status"`
	Cpu       float32 `json:"cpu"`
	Diskwrite int64   `json:"diskwrite"`
	ID        int64   `json:"vmid"`
	Uptime    int64   `json:"uptime"`
	Netout    int64   `json:"netout"`
	Name      string  `json:"name"`
	Maxmem    int64   `json:"maxmem"`
	Mem       int64   `json:"mem"`
	Diskread  int64   `json:"diskread"`
	Disk      int64   `json:"disk"`
	Netin     int64   `json:"netin"`
	Maxdisk   int64   `json:"maxdisk"`
	Cpus      int64   `json:"cpus"`
}

type VirtualMachineStatusResponse struct {
	Data VirtualMachineStatus `json:"data"`
}
type VirtualMachineStatus struct {
	Diskread       int64   `json:"diskread"`
	Maxmem         int64   `json:"maxmem"`
	Mem            int64   `json:"mem"`
	Disk           int64   `json:"disk"`
	Netin          int64   `json:"netin"`
	Cpus           float32 `json:"cpus"`
	Maxdisk        int64   `json:"maxdisk"`
	Balloon        int64   `json:"balloon"`
	RunningMachine string  `json:"running-machine"`
	RunningQemu    string  `json:"running-qemu"`
	Status         string  `json:"status"`
	Diskwrite      int64   `json:"diskwrite"`
	Cpu            float32 `json:"cpu"`
	Name           string  `json:"name"`
	Qmpstatus      string  `json:"qmpstatus"`
	Pid            int64   `json:"pid"`
	ID             int64   `json:"vmid"`
	Netout         int64   `json:"netout"`
	Uptime         int64   `json:"uptime"`
}

type VirtualMachineConfigResponse struct {
	Data VirtualMachineConfig `json:"data"`
}
type VirtualMachineConfig struct {
	Meta    string  `json:"meta"`
	Boot    string  `json:"boot"`
	Sockets int64   `json:"sockets,string"`
	Cpu     string  `json:"cpu"`
	IDE0    *string `json:"ide0,omitempty"`
	IDE1    *string `json:"ide1,omitempty"`
	IDE2    *string `json:"ide2,omitempty"`
	IDE3    *string `json:"ide3,omitempty"`
	Cores   int64   `json:"cores,string"`
	Numa    int64   `json:"numa,string"`
	Smbios1 string  `json:"smbios1"`
	Vmgenid string  `json:"vmgenid"`
	Net1    string  `json:"net1"`
	Ostype  string  `json:"ostype"`
	Scsi0   string  `json:"scsi0"`
	Digest  string  `json:"digest"`
	Scsihw  string  `json:"scsihw"`
	Memory  int64   `json:"memory,string"`
}
