package proxmox

// AsynchronousTaskResponse Proxmox can send us an ID of a task that is being executed asynchronously. We can use this to query the status of the task.
type AsynchronousTaskResponse struct {
	ID string `json:"data"`
}

// TaskResponse We can use the AsynchronousTaskResponse.ID to query the status of the task.
type TaskResponse struct {
	Data Task `json:"data"`
}

type Task struct {
	VirtualMachineID string  `json:"id"`
	Node             string  `json:"node"`
	PID              int64   `json:"pid"`
	StartTime        int64   `json:"starttime"`
	Status           string  `json:"status"`
	Type             string  `json:"type"`
	UPID             string  `json:"upid"`
	User             string  `json:"user"`
	ExitStatus       *string `json:"exitstatus"`
}
