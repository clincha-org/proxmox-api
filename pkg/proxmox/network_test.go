package proxmox

import (
	"fmt"
	"testing"
)

func TestGetDefaultInterface(t *testing.T) {
	client, err := NewClient(DefaultHostURL, TestUsername, TestPassword)
	if err != nil {
		t.Fatal(err)
	}

	nodes, err := client.GetNodes()
	if err != nil {
		t.Fatal(err)
	}

	node := &nodes[0]

	network, err := client.GetNetwork(node, "vmbr0")
	if err != nil {
		t.Fatal(err)
	}

	if network.Interface != "vmbr0" {
		t.Errorf("Incorrect interface returned. Expected vmbr0, got %v", network.Interface)
	}

	fmt.Printf("%+v", network)
}

// https://github.com/clincha-org/proxmox-api/issues/11
func TestNetworkUpdateWithGolangZeroValues(t *testing.T) {
	client, err := NewClient(DefaultHostURL, TestUsername, TestPassword)
	if err != nil {
		t.Fatal(err)
	}

	nodes, err := client.GetNodes()
	if err != nil {
		t.Fatal(err)
	}

	node := &nodes[0]

	request := NetworkRequest{
		Interface: "vmbr22",
		Type:      "bridge",
		VlanID:    2,
	}

	_, err = client.CreateNetwork(node, &request)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = client.DeleteNetwork(node, "vmbr22")
	})

	request = NetworkRequest{
		Interface: "vmbr22",
		Type:      "bridge",
	}

	network, err := client.UpdateNetwork(node, &request)
	if err != nil {
		t.Fatal(err)
	}

	if network.VlanID != 0 {
		t.Fatalf("Expected VLAN ID to be 1, got %v instead", network.VlanID)
	}
}
