package proxmox

import (
	"slices"
	"testing"
)

const PrimaryVirtualBridgeName = "vmbr0"

func TestGetNetworks(t *testing.T) {
	client, err := NewClient(DefaultHostURL, TestUsername, TestPassword)
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
	client, err := NewClient(DefaultHostURL, TestUsername, TestPassword)
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

func TestCreateDeleteNetworkBridge(t *testing.T) {
	client, err := NewClient(DefaultHostURL, TestUsername, TestPassword)
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

	err = client.DeleteNetwork(&nodes[0], network.Interface)
	if err != nil {
		t.Error(err)
	}
}

func TestCreateUpdateDeleteNetworkBridge(t *testing.T) {
	client, err := NewClient(DefaultHostURL, TestUsername, TestPassword)
	if err != nil {
		t.Error(err)
	}

	nodes, err := client.GetNodes()
	if err != nil {
		t.Error(err)
	}

	testInterface := "vmbr47"

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

	network.Autostart = 1

	network.BridgeStp = ""
	network.BridgeFd = ""
	network.Method6 = ""
	network.Families = nil
	network.Priority = 0
	network.Method = ""

	network, err = client.UpdateNetwork(&nodes[0], &network)
	if err != nil {
		t.Error(err)
	}

	if network.Autostart != 1 {
		t.Errorf("expected autostart to be 1 but got %v", network.Autostart)
	}

	err = client.DeleteNetwork(&nodes[0], network.Interface)
	if err != nil {
		t.Error(err)
	}
}
