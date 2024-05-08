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

	networkRequest := NetworkRequest{
		Interface: testInterface,
		Type:      "bridge",
	}

	network, err := client.CreateNetwork(&nodes[0], &networkRequest)
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
	testCIDR := "10.0.3.15/14"
	//testComments := "testComment"

	networkRequest := NetworkRequest{
		Type:      "bridge",
		Interface: testInterface,
	}

	network, err := client.CreateNetwork(&nodes[0], &networkRequest)
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

	networkRequest = NetworkRequest{
		Interface: network.Interface,
		Type:      network.Type,
		AutoStart: 1,
		CIDR:      testCIDR,
	}

	network, err = client.UpdateNetwork(&nodes[0], &networkRequest)
	if err != nil {
		t.Error(err)
	}

	if network.Autostart != 1 {
		t.Errorf("expected autostart to be 1 but got %v", network.Autostart)
	}

	if network.Cidr != testCIDR {
		t.Errorf("expected CIDR to be %v but got %v", testCIDR, network.Cidr)
	}

	networkRequest = NetworkRequest{
		Interface: network.Interface,
		Type:      network.Type,
		Comments:  "Hello",
		MTU:       8000,
	}

	network, err = client.UpdateNetwork(&nodes[0], &networkRequest)
	if err != nil {
		t.Error(err)
	}

	if network.Autostart != 1 {
		t.Errorf("expected autostart to be 1 but got %v", network.Autostart)
	}

	if network.Cidr != testCIDR {
		t.Errorf("expected CIDR to be %v but got %v", testCIDR, network.Cidr)
	}

	//if network.comments != testComments {
	//	t.Errorf("expected comments to be %v but got %v", testComments, network.Cidr)
	//}

	err = client.DeleteNetwork(&nodes[0], testInterface)
	if err != nil {
		t.Error(err)
	}
}
