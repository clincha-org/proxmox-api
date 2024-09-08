package proxmox

import (
	"github.com/clincha-org/proxmox-api/pkg/ide"
	"log/slog"
	"testing"
)

const UbuntuTestIso = "ubuntu-24.04.1-live-server-amd64.iso"

func TestGetVMs(t *testing.T) {
	client, err := NewClient(DefaultHostURL, TestUsername, TestPassword, slog.LevelDebug)
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.GetVMs("pve")

	if err != nil {
		t.Fatal(err)
	}
}

func TestGetVM(t *testing.T) {
	client, err := NewClient(DefaultHostURL, TestUsername, TestPassword, slog.LevelDebug)
	if err != nil {
		t.Fatal(err)
	}

	isoPath := "iso/" + UbuntuTestIso
	cdrom := ide.InternalDataStorage{
		ID:      2,
		Storage: "local",
		Path:    &isoPath,
	}
	scsi1 := "local-lvm:8"
	net1 := "model=virtio,bridge=vmbr0,firewall=1"
	scsiHardware := "virtio-scsi-pci"
	cores := int64(1)
	memory := int64(2048)

	request := VirtualMachine{
		ID:           102,
		IDEDevices:   &[]ide.InternalDataStorage{cdrom},
		SCSI1:        &scsi1,
		Net1:         &net1,
		SCSIHardware: &scsiHardware,
		Cores:        cores,
		Memory:       memory,
	}

	_, err = client.CreateVM("pve", &request, false)
	t.Cleanup(func() {
		err := client.DeleteVM("pve", 102)
		if err != nil {
			t.Fatal(err)
		}
	})

	if err != nil {
		t.Fatal(err)
	}

	vm, err := client.GetVM("pve", 102)

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

func TestCreateVM(t *testing.T) {
	client, err := NewClient(DefaultHostURL, TestUsername, TestPassword, slog.LevelDebug)
	if err != nil {
		t.Fatal(err)
	}

	isoPath := "iso/" + UbuntuTestIso
	cdrom := ide.InternalDataStorage{
		ID:      2,
		Storage: "local",
		Path:    &isoPath,
	}
	scsi1 := "local-lvm:8"
	net1 := "model=virtio,bridge=vmbr0,firewall=1"
	scsiHardware := "virtio-scsi-pci"
	cores := int64(1)
	memory := int64(2048)

	request := VirtualMachine{
		ID:           102,
		IDEDevices:   &[]ide.InternalDataStorage{cdrom},
		SCSI1:        &scsi1,
		Net1:         &net1,
		SCSIHardware: &scsiHardware,
		Cores:        cores,
		Memory:       memory,
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
	client, err := NewClient(DefaultHostURL, TestUsername, TestPassword, slog.LevelDebug)
	if err != nil {
		t.Fatal(err)
	}

	isoPath := "iso/" + UbuntuTestIso
	cdrom := ide.InternalDataStorage{
		ID:      2,
		Storage: "local",
		Path:    &isoPath,
	}
	scsi1 := "local-lvm:8"
	net1 := "model=virtio,bridge=vmbr0,firewall=1"
	scsiHardware := "virtio-scsi-pci"
	cores := int64(1)
	memory := int64(2048)

	request := VirtualMachine{
		ID:           102,
		IDEDevices:   &[]ide.InternalDataStorage{cdrom},
		SCSI1:        &scsi1,
		Net1:         &net1,
		SCSIHardware: &scsiHardware,
		Cores:        cores,
		Memory:       memory,
	}

	vm, err := client.CreateVM("pve", &request, true)
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

}

func TestUpdateVM(t *testing.T) {
	client, err := NewClient(DefaultHostURL, TestUsername, TestPassword, slog.LevelDebug)
	if err != nil {
		t.Fatal(err)
	}

	isoPath := "iso/" + UbuntuTestIso
	cdrom := ide.InternalDataStorage{
		ID:      2,
		Storage: "local",
		Path:    &isoPath,
	}
	newDiskSize := "4"
	ide1 := ide.InternalDataStorage{
		ID:      1,
		Storage: "local-lvm",
		Size:    &newDiskSize,
	}
	scsi1 := "local-lvm:8"
	net1 := "model=virtio,bridge=vmbr0,firewall=1"
	scsiHardware := "virtio-scsi-pci"
	vm := VirtualMachine{
		ID:           102,
		IDEDevices:   &[]ide.InternalDataStorage{cdrom, ide1},
		SCSI1:        &scsi1,
		Net1:         &net1,
		SCSIHardware: &scsiHardware,
		Cores:        1,
		Memory:       2048,
	}

	vm, err = client.CreateVM("pve", &vm, true)
	t.Cleanup(func() {
		err := client.DeleteVM("pve", 102)
		if err != nil {
			t.Fatal(err)
		}
	})
	if err != nil {
		t.Fatal(err)
	}

	if vm.IDEDevices == nil {
		t.Errorf("Expected ide devices, got %v", vm.IDEDevices)
	}

	if len(*vm.IDEDevices) != 2 {
		t.Errorf("Expected 2 ide devices, got %d", len(*vm.IDEDevices))
	}

	vm.Memory = 1024
	vm.Net1 = nil
	vm.SCSI1 = nil

	vm, err = client.UpdateVM("pve", &vm)
	if err != nil {
		t.Fatal(err)
	}

	if vm.Cores != 1 {
		t.Errorf("Expected 1 cores, got %d", vm.Cores)
	}

	if vm.Memory != 1024 {
		t.Errorf("Expected 1024 memory, got %d", vm.Memory)
	}

	if len(*vm.IDEDevices) != 2 {
		t.Errorf("Expected 2 ide devices, got %d", len(*vm.IDEDevices))
	}
}

func TestIDERemoval(t *testing.T) {
	client, err := NewClient(DefaultHostURL, TestUsername, TestPassword, slog.LevelDebug)
	if err != nil {
		t.Fatal(err)
	}

	isoPath := "iso/" + UbuntuTestIso
	size := "4"
	ide1 := ide.InternalDataStorage{
		ID:      1,
		Storage: "local-lvm",
		Size:    &size,
	}
	cdrom := ide.InternalDataStorage{
		ID:      2,
		Storage: "local",
		Path:    &isoPath,
	}
	scsi1 := "local-lvm:8"
	net1 := "model=virtio,bridge=vmbr0,firewall=1"
	scsiHardware := "virtio-scsi-pci"
	vm := VirtualMachine{
		ID:           102,
		IDEDevices:   &[]ide.InternalDataStorage{cdrom, ide1},
		SCSI1:        &scsi1,
		Net1:         &net1,
		SCSIHardware: &scsiHardware,
		Cores:        1,
		Memory:       2048,
	}

	vm, err = client.CreateVM("pve", &vm, true)
	t.Cleanup(func() {
		err := client.DeleteVM("pve", 102)
		if err != nil {
			t.Fatal(err)
		}
	})

	if err != nil {
		t.Fatal(err)
	}

	if len(*vm.IDEDevices) != 2 {
		t.Errorf("Expected 2 ide devices, got %d", len(*vm.IDEDevices))
	}

	vm.IDEDevices = &[]ide.InternalDataStorage{cdrom}

	vm, err = client.UpdateVM("pve", &vm)
	if err != nil {
		t.Fatal(err)
	}

	if len(*vm.IDEDevices) != 1 {
		t.Errorf("Expected 1 ide devices, got %d", len(*vm.IDEDevices))
	}
}

func TestCreateVMWithNoIDEDevices(t *testing.T) {
	client, err := NewClient(DefaultHostURL, TestUsername, TestPassword, slog.LevelDebug)
	if err != nil {
		t.Fatal(err)
	}

	scsi1 := "local-lvm:8"
	net1 := "model=virtio,bridge=vmbr0,firewall=1"
	scsiHardware := "virtio-scsi-pci"
	vm := VirtualMachine{
		ID:           102,
		SCSI1:        &scsi1,
		Net1:         &net1,
		SCSIHardware: &scsiHardware,
		Cores:        1,
		Memory:       2048,
	}

	vm, err = client.CreateVM("pve", &vm, true)
	t.Cleanup(func() {
		err := client.DeleteVM("pve", 102)
		if err != nil {
			t.Fatal(err)
		}
	})

	if err != nil {
		t.Fatal(err)
	}

	if vm.IDEDevices != nil {
		t.Errorf("Expected nil ide devices, got %v", vm.IDEDevices)
	}
}
