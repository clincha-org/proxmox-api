package proxmox

import (
	"slices"
	"testing"
)

const PrimaryVirtualBridgeName = "vmbr0"

func TestGetNetworks(t *testing.T) {
	host := "https://localhost:8006"
	username := "root@pam"
	password := "vagrant"

	client, err := NewClient(&host, &username, &password)
	if err != nil {
		t.Error(err)
	}

	nodes, err := client.GetNodes()
	if err != nil {
		t.Error(err)
	}

	networks, err := client.GetNetworks(&nodes[0])
	if err != nil {
		t.Error(err)
	}

	if len(networks) <= 0 {
		t.Error("No networks returned by GetNetwork")
	}

	PrimaryVirtualBridgeIndex := slices.IndexFunc(networks, func(network Network) bool {
		return network.Interface == PrimaryVirtualBridgeName
	})

	if PrimaryVirtualBridgeIndex == -1 {
		t.Errorf("Unable to find primary virtual bridge interface called %q", PrimaryVirtualBridgeName)
	}

	VirtualBridge := networks[PrimaryVirtualBridgeIndex]

	if VirtualBridge.Interface != PrimaryVirtualBridgeName {
		t.Errorf("Expected first network to be called %q. Got %q instead", PrimaryVirtualBridgeName, VirtualBridge.Interface)
	}

	if VirtualBridge.Type != "bridge" {
		t.Errorf("Expected first network to be of type bridge. Got type %q instead", VirtualBridge.Type)
	}
}

func TestGetNetwork(t *testing.T) {
	host := "https://localhost:8006"
	username := "root@pam"
	password := "vagrant"

	client, err := NewClient(&host, &username, &password)
	if err != nil {
		t.Error(err)
	}

	nodes, err := client.GetNodes()
	if err != nil {
		t.Error(err)
	}

	network, err := client.GetNetwork(&nodes[0], PrimaryVirtualBridgeName)
	if err != nil {
		t.Error(err)
	}

	if network.Interface != PrimaryVirtualBridgeName {
		t.Errorf("Expected first network to be called %q. Got %q instead", PrimaryVirtualBridgeName, network.Interface)
	}

	if network.Type != "bridge" {
		t.Errorf("Expected first network to be of type bridge. Got type %q instead", network.Type)
	}

	if network.Method != "static" {
		t.Errorf("Expected first network method to be static. Got %q instead", network.Method)
	}
}

func TestCreateAndDeleteNetworkBridge(t *testing.T) {
	host := "https://localhost:8006"
	username := "root@pam"
	password := "vagrant"

	client, err := NewClient(&host, &username, &password)
	if err != nil {
		t.Error(err)
	}

	nodes, err := client.GetNodes()
	if err != nil {
		t.Error(err)
	}

	testInterface := "vmbr99"

	network := Network{
		Type:      "bridge",
		Interface: testInterface,
	}

	network, err = client.CreateNetwork(&nodes[0], &network)
	if err != nil {
		t.Error(err)
	}

	if network.Interface != testInterface {
		t.Errorf("Expected network to be called %q. Got %q instead", testInterface, network.Interface)
	}

	if network.Type != "bridge" {
		t.Errorf("Expected network to be of type bridge. Got type %q instead", network.Type)
	}

	if network.Method != "manual" {
		t.Errorf("Expected network method to be manual. Got %q instead", network.Method)
	}

	err = client.DeleteNetwork(&nodes[0], &network)
	if err != nil {
		t.Error(err)
	}
}
