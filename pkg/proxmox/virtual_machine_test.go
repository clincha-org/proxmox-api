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

	cdrom := "local:iso/ubuntu-24.04-live-server-amd64.iso"
	scsi1 := "local-lvm:8"
	net1 := "model=virtio,bridge=vmbr0,firewall=1"
	scsiHardware := "virtio-scsi-pci"
	cores := int64(1)
	memory := int64(2048)

	request := VirtualMachineRequest{
		ID:           102,
		Cdrom:        &cdrom,
		SCSI1:        &scsi1,
		Net1:         &net1,
		SCSIHardware: &scsiHardware,
		Cores:        &cores,
		Memory:       &memory,
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

	cdrom := "local:iso/ubuntu-24.04-live-server-amd64.iso"
	scsi1 := "local-lvm:8"
	net1 := "model=virtio,bridge=vmbr0,firewall=1"
	scsiHardware := "virtio-scsi-pci"
	cores := int64(1)
	memory := int64(2048)

	request := VirtualMachineRequest{
		ID:           102,
		Cdrom:        &cdrom,
		SCSI1:        &scsi1,
		Net1:         &net1,
		SCSIHardware: &scsiHardware,
		Cores:        &cores,
		Memory:       &memory,
	}

	_, err = client.CreateVM("pve", &request, true)

	if err != nil {
		t.Fatal(err)
	}
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

func TestUpdateVM(t *testing.T) {
	DebugLogs()
	client, err := NewClient(DefaultHostURL, TestUsername, TestPassword)
	if err != nil {
		t.Fatal(err)
	}

	cdrom := "local:iso/ubuntu-24.04-live-server-amd64.iso"
	scsi1 := "local-lvm:8"
	net1 := "model=virtio,bridge=vmbr0,firewall=1"
	scsiHardware := "virtio-scsi-pci"
	cores := int64(1)
	memory := int64(2048)

	request := VirtualMachineRequest{
		ID:           102,
		Cdrom:        &cdrom,
		SCSI1:        &scsi1,
		Net1:         &net1,
		SCSIHardware: &scsiHardware,
		Cores:        &cores,
		Memory:       &memory,
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

	cores = int64(1)
	memory = int64(1024)

	request.Cores = &cores
	request.Memory = &memory
	request.Net1 = nil
	request.SCSI1 = nil

	vm, err := client.UpdateVM("pve", &request)
	if err != nil {
		t.Fatal(err)
	}

	if vm.Cores != 1 {
		t.Errorf("Expected 1 cores, got %d", vm.Cores)
	}

	if vm.Memory != 1024 {
		t.Errorf("Expected 1024 memory, got %d", vm.Memory)
	}
}
