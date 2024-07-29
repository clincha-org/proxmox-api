package proxmox

import (
	"fmt"
	"log/slog"
	"os"
	"testing"
)

func DebugLogs() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)
}

func TestGetVMs(t *testing.T) {
	client, err := NewClient(DefaultHostURL, TestUsername, TestPassword)
	if err != nil {
		t.Fatal(err)
	}

	nodes, err := client.GetNodes()
	if err != nil {
		t.Fatal(err)
	}

	node := &nodes[0]

	vm, err := client.GetVMs(node)

	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("VM: %+v", vm)
}

func TestGetVM(t *testing.T) {
	client, err := NewClient(DefaultHostURL, TestUsername, TestPassword)
	if err != nil {
		t.Fatal(err)
	}

	nodes, err := client.GetNodes()
	if err != nil {
		t.Fatal(err)
	}

	node := &nodes[0]

	vm, err := client.GetVM(node, 100)

	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("VM: %+v", vm.Ide2)
}

func TestCreateVM(t *testing.T) {
	DebugLogs()
	client, err := NewClient(DefaultHostURL, TestUsername, TestPassword)
	if err != nil {
		t.Fatal(err)
	}

	nodes, err := client.GetNodes()
	if err != nil {
		t.Fatal(err)
	}

	node := &nodes[0]

	request := VirtualMachineRequest{
		VMID: 101,
	}

	vm, err := client.CreateVM(node, &request)

	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("VM: %+v", vm)
}
