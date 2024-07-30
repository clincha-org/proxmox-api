package proxmox

type VirtualMachineRequest struct {
	ID           int64  `json:"vmid"`
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

type VirtualMachineConfigResponse struct {
	Data VirtualMachineConfig `json:"data"`
}
type VirtualMachineConfig struct {
	Meta    string `json:"meta"`
	Boot    string `json:"boot"`
	Sockets int64  `json:"sockets"`
	Cpu     string `json:"cpu"`
	Ide2    string `json:"ide2"`
	Cores   int64  `json:"cores"`
	Numa    int64  `json:"numa"`
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
