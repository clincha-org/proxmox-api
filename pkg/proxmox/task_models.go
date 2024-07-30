package proxmox

type TaskResponse struct {
	Data Task `json:"data"`
}

type Task struct {
	ID        string `json:"id"`
	Node      string `json:"node"`
	PID       int64  `json:"pid"`
	StartTime int64  `json:"starttime"`
	Status    string `json:"status"`
	Type      string `json:"type"`
	UPID      string `json:"upid"`
	User      string `json:"user"`
}

type JobResponse struct {
	ID string `json:"data"`
}
