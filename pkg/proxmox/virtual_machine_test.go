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

	vm, err := client.GetVMs("pve")

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

	vm, err := client.GetVM("pve", 100)

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

	request := VirtualMachineRequest{
		VMID:         102,
		Cdrom:        "local:iso/ubuntu-22.04.4-live-server-amd64.iso",
		SCSI1:        "local-lvm:8",
		Net1:         "model=virtio,bridge=vmbr0,firewall=1",
		SCSIHardware: "virtio-scsi-pci",
	}

	vm, err := client.CreateVM("pve", &request, false)

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

	request := VirtualMachineRequest{
		VMID:         102,
		Cdrom:        "local:iso/ubuntu-22.04.4-live-server-amd64.iso",
		SCSI1:        "local-lvm:8",
		Net1:         "model=virtio,bridge=vmbr0,firewall=1",
		SCSIHardware: "virtio-scsi-pci",
	}

	vm, err := client.CreateVM("pve", &request, true)

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

	err = client.DeleteVM("pve", 102)
	if err != nil {
		t.Fatal(err)
	}
}
