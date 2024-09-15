package proxmox

import (
	"encoding/json"
	"fmt"
	"time"
)

const taskPath = "/tasks"

func (client *Client) GetTaskStatus(node string, id string) (Task, error) {
	path := client.Host + ApiPath + NodesPath + "/" + node + taskPath + "/" + id + "/status"
	body, err := client.MakeRESTRequest("GET", path, nil)
	if err != nil {
		return Task{}, fmt.Errorf("GetTaskStatus-make-request: %w", err)
	}

	taskModel := TaskResponse{}
	err = json.Unmarshal(body, &taskModel)
	if err != nil {
		return Task{}, fmt.Errorf("GetTaskStatus-unmarshal-response: %w", err)
	}

	return taskModel.Data, nil

}

func (client *Client) AwaitAsynchronousTask(node string, taskID string) error {
	task, err := client.GetTaskStatus(node, taskID)
	if err != nil {
		return fmt.Errorf("AwaitAsynchronousTask-get-task-status: '%s' %w", taskID, err)
	}

	// Poll the task status until the task is completed
	for ok := true; ok; ok = task.Status != "stopped" {
		task, err = client.GetTaskStatus(node, task.UPID)
		if err != nil {
			return fmt.Errorf("AwaitAsynchronousTask-get-job-loop: '%s' %w", taskID, err)
		}
		// Sleep for 1 second before polling again
		time.Sleep(1 * time.Second)
	}

	if task.ExitStatus == nil {
		return fmt.Errorf("AwaitAsynchronousTask-exit-status-nil: '%s'", taskID)
	}

	if *task.ExitStatus != "OK" {
		return fmt.Errorf("AwaitAsynchronousTask-exit-status-not-ok: '%s' %s", taskID, *task.ExitStatus)
	}

	return nil
}
