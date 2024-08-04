package proxmox

import (
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

	_, err = client.GetVMs("pve")

	if err != nil {
		t.Fatal(err)
	}
}

func TestGetVM(t *testing.T) {
	DebugLogs()
	client, err := NewClient(DefaultHostURL, TestUsername, TestPassword)
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.GetVM("pve", 102)

	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateVM(t *testing.T) {
	DebugLogs()
	client, err := NewClient(DefaultHostURL, TestUsername, TestPassword)
	if err != nil {
		t.Fatal(err)
	}

	request := VirtualMachineRequest{
		ID:           102,
		Cdrom:        "local:iso/ubuntu-22.04.4-live-server-amd64.iso",
		SCSI1:        "local-lvm:8",
		Net1:         "model=virtio,bridge=vmbr0,firewall=1",
		SCSIHardware: "virtio-scsi-pci",
		Cores:        1,
		Memory:       2048,
	}

	vm, err := client.CreateVM("pve", &request, false)
	t.Cleanup(func() {
		err := client.DeleteVM("pve", 102)
		if err != nil {
			t.Fatal(err)
		}
	})

	if err != nil {
		t.Fatal(err)
	}

	if vm.Cores != 1 {
		t.Errorf("Expected 1 core, got %d", vm.Cores)
	}

	if vm.Memory != 2048 {
		t.Errorf("Expected 2048 memory, got %d", vm.Memory)
	}

}

func TestCreateVMWithStart(t *testing.T) {
	DebugLogs()
	client, err := NewClient(DefaultHostURL, TestUsername, TestPassword)
	if err != nil {
		t.Fatal(err)
	}

	request := VirtualMachineRequest{
		ID:           102,
		Cdrom:        "local:iso/ubuntu-22.04.4-live-server-amd64.iso",
		SCSI1:        "local-lvm:8",
		Net1:         "model=virtio,bridge=vmbr0,firewall=1",
		SCSIHardware: "virtio-scsi-pci",
		Cores:        1,
		Memory:       2048,
	}

	_, err = client.CreateVM("pve", &request, true)
	t.Cleanup(func() {
		err := client.DeleteVM("pve", 102)
		if err != nil {
			t.Fatal(err)
		}
	})

	if err != nil {
		t.Fatal(err)
	}
}
