package proxmox

import (
	"strings"
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

	TestAddress := "10.0.3.13"
	TestNetmask := "255.255.255.0"
	request := NetworkRequest{
		Interface: "vmbr22",
		Type:      "bridge",
		Address:   &TestAddress,
		Netmask:   &TestNetmask,
	}

	network, err := client.CreateNetwork(node, &request)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = client.DeleteNetwork(node, "vmbr22")
	})

	if *network.Netmask != TestNetmask {
		t.Fatalf("Expected netmask to be %v, got %v instead", TestNetmask, network.Netmask)
	}

	TestNetmask = "255.255.255.252"
	request = NetworkRequest{
		Interface: "vmbr22",
		Type:      "bridge",
		Address:   &TestAddress,
		Netmask:   &TestNetmask,
	}

	network, err = client.UpdateNetwork(node, &request)
	if err != nil {
		t.Fatal(err)
	}

	if *network.Netmask != TestNetmask {
		t.Fatalf("Expected netmask to be %v, got %v instead", TestNetmask, network.Netmask)
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

	TestCIDR := "10.0.3.13/24"
	request := NetworkRequest{
		Interface: "vmbr22",
		Type:      "bridge",
		CIDR:      &TestCIDR,
	}

	network, err := client.CreateNetwork(node, &request)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = client.DeleteNetwork(node, "vmbr22")
	})

	if *network.Address != strings.Split(TestCIDR, "/")[0] {
		t.Fatalf("Expected network address to be %v, got %v instead", strings.Split(TestCIDR, "/")[0], network.Address)
	}

	if *network.Netmask != "255.255.255.0" {
		t.Fatalf("Expected netmask to be 255.255.255.0, got %v instead", network.Netmask)
	}

	TestAddress := "10.0.3.13"
	TestNetmask := "255.255.255.252"
	request = NetworkRequest{
		Interface: "vmbr22",
		Type:      "bridge",
		Address:   &TestAddress,
		Netmask:   &TestNetmask,
	}

	network, err = client.UpdateNetwork(node, &request)
	if err != nil {
		t.Fatal(err)
	}

	if *network.CIDR != "10.0.3.13/30" {
		t.Fatalf("Expected CIDR to be 10.0.3.13/30, got %v instead", network.CIDR)
	}

	if *network.Netmask != "255.255.255.252" {
		t.Fatalf("Expected netmask to be 255.255.255.252, got %v instead", network.Netmask)
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

	autostart := true
	request := NetworkRequest{
		Interface: "vmbr22",
		Type:      "bridge",
		AutoStart: &autostart,
	}

	network, err := client.CreateNetwork(node, &request)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = client.DeleteNetwork(node, "vmbr22")
	})

	if network.Autostart != 1 {
		t.Fatalf("Expected network autostart to be 1, got %v instead", network.Autostart)
	}

	autostart = false
	request = NetworkRequest{
		Interface: "vmbr22",
		Type:      "bridge",
		AutoStart: &autostart,
	}

	network, err = client.UpdateNetwork(node, &request)
	if err != nil {
		t.Fatal(err)
	}

	if network.Autostart != 0 {
		t.Fatalf("Expected network autostart to be 0, got %v instead", network.Autostart)
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

	TestCIDR := "10.0.3.13/24"
	TestBridgePorts := "enp0s4"
	request = NetworkRequest{
		Interface:   "vmbr22",
		Type:        "bridge",
		CIDR:        &TestCIDR,
		BridgePorts: &TestBridgePorts,
	}

	network, err = client.CreateNetwork(node, &request)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = client.DeleteNetwork(node, "vmbr22")
	})

	if *network.BridgePorts != TestBridgePorts {
		t.Fatalf("Expected bridge port to be %v, got %v instead", TestBridgePorts, network.BridgePorts)
	}
}

func TestNetworkOmittedFields(t *testing.T) {
	client, err := NewClient(DefaultHostURL, TestUsername, TestPassword)
	if err != nil {
		t.Fatal(err)
	}

	nodes, err := client.GetNodes()
	if err != nil {
		t.Fatal(err)
	}

	node := &nodes[0]

	TestAutostart := true
	TestComments := "test comments"
	request := NetworkRequest{
		Interface: "vmbr22",
		Type:      "bridge",
		AutoStart: &TestAutostart,
		Comments:  TestComments,
	}

	network, err := client.CreateNetwork(node, &request)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = client.DeleteNetwork(node, "vmbr22")
	})
	if network.Autostart != 1 {
		t.Fatalf("Expected network autostart to be 1, got %v instead", network.Autostart)
	}
	if network.Comments != TestComments {
		t.Fatalf("Expected network comments to be %v, got %v instead", TestComments, network.Comments)
	}

	request = NetworkRequest{
		Interface: "vmbr22",
		Type:      "bridge",
	}

	network, err = client.UpdateNetwork(node, &request)
	if err != nil {
		t.Fatal(err)
	}
	if network.Autostart != 1 {
		t.Fatalf("Expected network autostart to be 1, got %v instead", network.Autostart)
	}
	if network.Comments != "" {
		t.Fatalf("Expected network comments to be %v, got %v instead", TestComments, network.Comments)
	}

	TestAutostart = false
	TestComments = ""
	request = NetworkRequest{
		Interface: "vmbr22",
		Type:      "bridge",
		AutoStart: &TestAutostart,
		Comments:  TestComments,
	}

	network, err = client.UpdateNetwork(node, &request)
	if err != nil {
		t.Fatal(err)
	}
	if network.Autostart != 0 {
		t.Fatalf("Expected network autostart to be 0, got %v instead", network.Autostart)
	}
	if network.Comments != "" {
		t.Fatalf("Expected network comments to be empty, got %v instead", network.Comments)
	}
}

func TestSubnetMaskReturnedInSameFormat(t *testing.T) {
	client, err := NewClient(DefaultHostURL, TestUsername, TestPassword)
	if err != nil {
		t.Fatal(err)
	}

	nodes, err := client.GetNodes()
	if err != nil {
		t.Fatal(err)
	}

	node := &nodes[0]

	TestAddress := "10.1.2.3"
	TestNetmask := "255.255.255.0"
	request := NetworkRequest{
		Interface: "vmbr22",
		Type:      "bridge",
		Address:   &TestAddress,
		Netmask:   &TestNetmask,
	}

	network, err := client.CreateNetwork(node, &request)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = client.DeleteNetwork(node, "vmbr22")
	})

	if *network.Netmask != TestNetmask {
		t.Fatalf("Expected network mask to be %v, got %v instead", TestNetmask, network.Autostart)
	}
}
