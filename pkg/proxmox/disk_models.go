package proxmox

type proxmoxDisk struct {
}

type scsiDisk struct {
	proxmoxDisk
	File string `json:"file"`
	Size int    `json:"size"`
	SSD  int    `json:"ssd"`
}
