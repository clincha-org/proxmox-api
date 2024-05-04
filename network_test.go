package main

import (
	"github.com/clincha-org/proxmox-api/pkg/proxmox"
	"slices"
	"testing"
)

func TestGetNetwork(t *testing.T) {
	host := "https://localhost:8006"
	username := "root@pam"
	password := "vagrant"

	client, err := proxmox.NewClient(&host, &username, &password)
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

	PrimaryVirtualBridgeName := "vmbr0"

	PrimaryVirtualBridgeIndex := slices.IndexFunc(networks, func(network proxmox.Network) bool {
		return network.Interface == PrimaryVirtualBridgeName
	})

	if PrimaryVirtualBridgeIndex == -1 {
		t.Errorf("Unable to find primary virtual bridge interface called %v", PrimaryVirtualBridgeName)
	}

	VirtualBridge := networks[PrimaryVirtualBridgeIndex]

	if VirtualBridge.Interface != PrimaryVirtualBridgeName {
		t.Errorf("Expected first network to be called %v. Got %v instead", PrimaryVirtualBridgeName, VirtualBridge.Interface)
	}

	if VirtualBridge.Type != "bridge" {
		t.Errorf("Expected first network to be of type bridge. Got type %v instead", VirtualBridge.Type)
	}
}

func TestCreateAndDeleteNetworkBridge(t *testing.T) {
	host := "https://localhost:8006"
	username := "root@pam"
	password := "vagrant"

	client, err := proxmox.NewClient(&host, &username, &password)
	if err != nil {
		t.Error(err)
	}

	nodes, err := client.GetNodes()
	if err != nil {
		t.Error(err)
	}

	network := proxmox.Network{
		Type:      "bridge",
		Interface: "vmbr99",
	}

	network, err = client.CreateNetwork(&nodes[0], &network)
	if err != nil {
		t.Error(err)
	}

	err = client.DeleteNetwork(&nodes[0], &network)
	if err != nil {
		t.Error(err)
	}
}
