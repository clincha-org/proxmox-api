package proxmox

import (
	"github.com/clincha-org/proxmox-api/pkg/ide"
	"log/slog"
	"testing"
)

func TestClone(t *testing.T) {
	client, err := NewClient(DefaultHostURL, TestUsername, TestPassword, slog.LevelDebug)
	if err != nil {
		t.Fatal(err)
	}

	isoPath := "iso/" + UbuntuTestIso
	size := "4"
	diskStorage := "local-lvm"
	ide1 := ide.InternalDataStorage{
		ID:      1,
		Storage: &diskStorage,
		Size:    &size,
	}
	cdromStorage := "local"
	cdrom := ide.InternalDataStorage{
		ID:      2,
		Storage: &cdromStorage,
		Path:    &isoPath,
	}
	scsi1 := "local-lvm:8"
	net0 := "model=virtio,bridge=vmbr0,firewall=1"
	scsiHardware := "virtio-scsi-pci"
	vm := VirtualMachine{
		ID:           103,
		IDEDevices:   &[]ide.InternalDataStorage{cdrom, ide1},
		SCSI1:        &scsi1,
		Net0:         &net0,
		SCSIHardware: &scsiHardware,
		Cores:        1,
		Memory:       2048,
	}

	vm, err = client.CloneVM("pve", &vm, 100, false)
	t.Cleanup(func() {
		err := client.DeleteVM("pve", 103)
		if err != nil {
			t.Fatal(err)
		}
	})

	vm, err = client.UpdateVM("pve", &vm)
	if err != nil {
		t.Fatal(err)
	}

	if vm.Parent == nil {
		t.Errorf("Expected parent, got nil")
	}

	if *vm.Parent != 100 {
		t.Errorf("Expected parent 100, got %d", *vm.Parent)
	}
}
