package main

import (
	"github.com/clincha-org/proxmox-api/pkg/proxmox"
	"testing"
)

func TestGetNodes(t *testing.T) {
	host := "https://localhost:8006"
	username := "root@pam"
	password := "vagrant"

	client, err := proxmox.NewClient(&host, &username, &password)
	if err != nil {
		t.Error(err)
	}

	node, err := client.GetNodes()
	if err != nil {
		t.Error(err)
	}

	if len(node.Data) == 0 {
		t.Error("Expecting at least one node, got none")
	}

	if node.Data[0].Status != "online" {
		t.Errorf("Expected node to be online, got %q", node.Data[0].Status)
	}

	if node.Data[0].Type != "node" {
		t.Errorf("Expected type of node to be node, got %q", node.Data[0].Type)
	}

	if node.Data[0].Uptime <= 0 {
		t.Errorf("Expected node uptime to be greater than 0, got %q", node.Data[0].Uptime)
	}

}
