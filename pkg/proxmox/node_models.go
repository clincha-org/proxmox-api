package proxmox

type NodeResponse struct {
	Data []Node `json:"data"`
}

type Node struct {
	Type           string  `json:"type"`
	Maxcpu         int     `json:"maxcpu"`
	Cpu            float64 `json:"cpu"`
	Status         string  `json:"status"`
	Maxmem         int     `json:"maxmem"`
	SslFingerprint string  `json:"ssl_fingerprint"`
	Mem            int     `json:"mem"`
	Id             string  `json:"id"`
	Node           string  `json:"node"`
	Disk           int64   `json:"disk"`
	Uptime         int     `json:"uptime"`
	Maxdisk        int64   `json:"maxdisk"`
	Level          string  `json:"level"`
}
