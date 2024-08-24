package ide

import (
	"testing"
)

func TestIdeMarshal(t *testing.T) {
	proxmox_ide_response := "local:iso/ubuntu-24.04-live-server-amd64.iso,media=cdrom,size=2690412K"

	cdrom := &InternalDataStorage{}
	err := Unmarshal(proxmox_ide_response, cdrom)
	if err != nil {
		t.Fatal(err)
	}
}
