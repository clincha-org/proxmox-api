package proxmox

type VirtualMachineRequest struct {
	ID           int    `json:"vmid"`
	Cdrom        string `json:"cdrom"`
	SCSI1        string `json:"scsi1"`
	Net1         string `json:"net1"`
	SCSIHardware string `json:"scsihw"`
}

type VirtualMachinesResponse struct {
	Data []VirtualMachine `json:"data"`
}

type VirtualMachine struct {
	Status    string  `json:"status"`
	Cpu       float32 `json:"cpu"`
	Diskwrite int     `json:"diskwrite"`
	Vmid      int     `json:"vmid"`
	Uptime    int     `json:"uptime"`
	Netout    int     `json:"netout"`
	Name      string  `json:"name"`
	Maxmem    int64   `json:"maxmem"`
	Mem       int     `json:"mem"`
	Diskread  int     `json:"diskread"`
	Disk      int     `json:"disk"`
	Netin     int     `json:"netin"`
	Maxdisk   int64   `json:"maxdisk"`
	Cpus      int     `json:"cpus"`
}

type VirtualMachineConfigResponse struct {
	Data VirtualMachineConfig `json:"data"`
}
type VirtualMachineConfig struct {
	Meta    string `json:"meta"`
	Boot    string `json:"boot"`
	Sockets int    `json:"sockets"`
	Cpu     string `json:"cpu"`
	Ide2    string `json:"ide2"`
	Cores   int    `json:"cores"`
	Numa    int    `json:"numa"`
	Smbios1 string `json:"smbios1"`
	Vmgenid string `json:"vmgenid"`
	Net0    string `json:"net0"`
	Ostype  string `json:"ostype"`
	Scsi0   string `json:"scsi0"`
	Digest  string `json:"digest"`
	Scsihw  string `json:"scsihw"`
	Memory  string `json:"memory"`
}

type VirtualMachineStatusResponse struct {
	Data VirtualMachineStatus `json:"data"`
}

type VirtualMachineStatus struct {
	Diskread       int     `json:"diskread"`
	Maxmem         int     `json:"maxmem"`
	Mem            int     `json:"mem"`
	Disk           int     `json:"disk"`
	Netin          int     `json:"netin"`
	Cpus           float32 `json:"cpus"`
	Maxdisk        int64   `json:"maxdisk"`
	Balloon        int     `json:"balloon"`
	RunningMachine string  `json:"running-machine"`
	RunningQemu    string  `json:"running-qemu"`
	Status         string  `json:"status"`
	Diskwrite      int     `json:"diskwrite"`
	Cpu            float32 `json:"cpu"`
	Name           string  `json:"name"`
	Qmpstatus      string  `json:"qmpstatus"`
	Pid            int     `json:"pid"`
	Vmid           int     `json:"vmid"`
	Netout         int     `json:"netout"`
	Uptime         int     `json:"uptime"`
}
