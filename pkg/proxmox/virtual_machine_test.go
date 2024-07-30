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
		VMID:         102,
		Cdrom:        "local:iso/ubuntu-22.04.4-live-server-amd64.iso",
		SCSI1:        "local-lvm:8",
		Net1:         "model=virtio,bridge=vmbr0,firewall=1",
		SCSIHardware: "virtio-scsi-pci",
	}

	vm, err := client.CreateVM(node, &request, false)

	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("VM: %+v", vm)
}

func TestCreateVMWithStart(t *testing.T) {
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
		VMID:         102,
		Cdrom:        "local:iso/ubuntu-22.04.4-live-server-amd64.iso",
		SCSI1:        "local-lvm:8",
		Net1:         "model=virtio,bridge=vmbr0,firewall=1",
		SCSIHardware: "virtio-scsi-pci",
	}

	vm, err := client.CreateVM(node, &request, true)

	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("VM: %+v", vm)
}

func TestDeleteVM(t *testing.T) {
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

	err = client.DeleteVM(node, 102)
	if err != nil {
		t.Fatal(err)
	}
}
