package ide

import (
	"testing"
)

const UbuntuTestIso = "ubuntu-24.04.1-live-server-amd64.iso"

func TestIdeMarshal(t *testing.T) {
	proxmox_ide_response := "local:iso/" + UbuntuTestIso + ",media=cdrom,size=2690412K"

	cdrom := &InternalDataStorage{}
	err := Unmarshal(2, proxmox_ide_response, cdrom)
	if err != nil {
		t.Fatal(err)
	}
}
