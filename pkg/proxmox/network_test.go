package proxmox

import (
	"slices"
	"testing"
)

const PrimaryVirtualBridgeName = "vmbr0"

func TestGetNetworks(t *testing.T) {
	client, err := NewClient(DefaultHostURL, TestUsername, TestPassword)
	if err != nil {
		t.Fatal(err)
	}

	nodes, err := client.GetNodes()
	if err != nil {
		t.Fatal(err)
	}

	networks, err := client.GetNetworks(&nodes[0])
	if err != nil {
		t.Fatal(err)
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
		t.Fatal(err)
	}

	nodes, err := client.GetNodes()
	if err != nil {
		t.Fatal(err)
	}

	network, err := client.GetNetwork(&nodes[0], PrimaryVirtualBridgeName)
	if err != nil {
		t.Fatal(err)
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
		t.Fatal(err)
	}

	nodes, err := client.GetNodes()
	if err != nil {
		t.Fatal(err)
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
		t.Fatal(err)
	}

	nodes, err := client.GetNodes()
	if err != nil {
		t.Fatal(err)
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

	if network.CIDR != testCIDR {
		t.Errorf("expected CIDR to be %v but got %v", testCIDR, network.CIDR)
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

	if network.CIDR != testCIDR {
		t.Errorf("expected CIDR to be %v but got %v", testCIDR, network.CIDR)
	}

	//if network.comments != testComments {
	//	t.Errorf("expected comments to be %v but got %v", testComments, network.Cidr)
	//}

	err = client.DeleteNetwork(&nodes[0], testInterface)
	if err != nil {
		t.Error(err)
	}
}

func TestBondedNetworkConfiguration(t *testing.T) {
	client, err := NewClient(DefaultHostURL, TestUsername, TestPassword)
	if err != nil {
		t.Fatal(err)
	}

	nodes, err := client.GetNodes()
	if err != nil {
		t.Fatal(err)
	}

	request := NetworkRequest{
		Interface: "eno1",
		Type:      "eth",
		AutoStart: 1,
	}

	_, err = client.CreateNetwork(&nodes[0], &request)
	if err != nil {
		t.Fatal(err)
	}

	request.Interface = "eno2"
	_, err = client.CreateNetwork(&nodes[0], &request)
	if err != nil {
		t.Fatal(err)
	}

	Interface := "bond1"
	Type := "bond"
	OVSBonds := "eno1,eno2"
	BondMode := "active-backup"

	request = NetworkRequest{
		Interface: Interface,
		Type:      Type,
		OVSBonds:  OVSBonds,
		BondMode:  BondMode,
	}
	network, err := client.CreateNetwork(&nodes[0], &request)
	if err != nil {
		t.Error(err)
	}

	if network.Interface != Interface {
		t.Errorf("expected %v for Interface but got %v", Interface, network.Interface)
	}
	if network.Type != Type {
		t.Errorf("expected %v for Type but got %v", Type, network.Type)
	}
	if network.OVSBonds != OVSBonds {
		t.Errorf("expected %v for OVSBonds but got %v", OVSBonds, network.OVSBonds)
	}
	if network.BondMode != BondMode {
		t.Errorf("expected %v for BondMode but got %v", BondMode, network.BondMode)
	}

	err = client.DeleteNetwork(&nodes[0], "bond1")
	if err != nil {
		t.Fatal(err)
	}
	err = client.DeleteNetwork(&nodes[0], "eno1")
	if err != nil {
		t.Fatal(err)
	}
	err = client.DeleteNetwork(&nodes[0], "eno2")
	if err != nil {
		t.Fatal(err)
	}

}

func TestVLANNetworkConfiguration(t *testing.T) {
	client, err := NewClient(DefaultHostURL, TestUsername, TestPassword)
	if err != nil {
		t.Fatal(err)
	}

	nodes, err := client.GetNodes()
	if err != nil {
		t.Fatal(err)
	}

	request := NetworkRequest{
		Interface: "vmbr74",
		Type:      "bridge",
		VlanID:    3,
	}

	network, err := client.CreateNetwork(&nodes[0], &request)
	if err != nil {
		t.Fatal(err)
	}

	if network.VlanID != 3 {
		t.Errorf("expected %v for VlanID but got %v", 3, network.VlanID)
	}

	err = client.DeleteNetwork(&nodes[0], "vmbr74")
	if err != nil {
		t.Fatal(err)
	}
}
