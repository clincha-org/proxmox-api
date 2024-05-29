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
		VlanID:    0,
	}

	network, err := client.UpdateNetwork(node, &request)
	if err != nil {
		t.Fatal(err)
	}

	if network.VlanID != 2 {
		t.Fatalf("Expected VLAN ID to be 2, got %v instead", network.VlanID)
	}
}

func TestNetworkNetmaskUpdate(t *testing.T) {
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
		Address:   "10.0.3.13",
		Netmask:   "255.255.255.0",
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
		Address:   "10.0.3.13",
		Netmask:   "255.255.255.252",
	}

	network, err := client.UpdateNetwork(node, &request)
	if err != nil {
		t.Fatal(err)
	}

	if network.Netmask != 30 {
		t.Fatalf("Expected netmask to be 30, got %v instead", network.Netmask)
	}
}

func TestNetworkCIDRUpdate(t *testing.T) {
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
		CIDR:      "10.0.3.13/24",
	}

	network, err := client.CreateNetwork(node, &request)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = client.DeleteNetwork(node, "vmbr22")
	})

	if network.Address != "10.0.3.13" {
		t.Fatalf("Expected network address to be 10.0.3.13, got %v instead", network.Address)
	}

	if network.Netmask != 24 {
		t.Fatalf("Expected netmask to be 24, got %v instead", network.Netmask)
	}

	request = NetworkRequest{
		Interface: "vmbr22",
		Type:      "bridge",
		Address:   "10.0.3.13",
		Netmask:   "255.255.255.252",
	}

	network, err = client.UpdateNetwork(node, &request)
	if err != nil {
		t.Fatal(err)
	}

	if network.CIDR != "10.0.3.13/30" {
		t.Fatalf("Expected CIDR to be 10.0.3.13/30, got %v instead", network.CIDR)
	}

	if network.Netmask != 30 {
		t.Fatalf("Expected netmask to be 30, got %v instead", network.Netmask)
	}
}

func TestNetworkAutostart(t *testing.T) {
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
		AutoStart: true,
	}

	network, err := client.CreateNetwork(node, &request)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = client.DeleteNetwork(node, "vmbr22")
	})

	if network.Autostart != true {
		t.Fatalf("Expected network autostart to be true, got %v instead", network.Autostart)
	}

	request = NetworkRequest{
		Interface: "vmbr22",
		Type:      "bridge",
		AutoStart: false,
	}

	network, err = client.UpdateNetwork(node, &request)
	if err != nil {
		t.Fatal(err)
	}

	if network.Autostart != false {
		t.Fatalf("Expected network autostart to be false, got %v instead", network.Autostart)
	}
}

func TestNetworkBridgePorts(t *testing.T) {
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
		Interface: "enp0s4",
		Type:      "eth",
	}

	network, err := client.CreateNetwork(node, &request)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = client.DeleteNetwork(node, "enp0s4")
	})

	request = NetworkRequest{
		Interface:   "vmbr22",
		Type:        "bridge",
		CIDR:        "10.0.3.13/24",
		BridgePorts: "enp0s4",
	}

	network, err = client.CreateNetwork(node, &request)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = client.DeleteNetwork(node, "vmbr22")
	})

	if network.BridgePorts != "enp0s4" {
		t.Fatalf("Expected bridge port to be enp0s4, got %v instead", network.BridgePorts)
	}
}
