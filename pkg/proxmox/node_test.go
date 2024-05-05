package proxmox

import (
	"testing"
)

func TestGetNodes(t *testing.T) {
	client, err := NewClient(DefaultHostURL, TestUsername, TestPassword)
	if err != nil {
		t.Error(err)
	}

	nodes, err := client.GetNodes()
	if err != nil {
		t.Error(err)
	}

	if len(nodes) == 0 {
		t.Error("Expecting at least one nodes, got none")
	}

	if nodes[0].Status != "online" {
		t.Errorf("Expected nodes to be online, got %q", nodes[0].Status)
	}

	if nodes[0].Type != "node" {
		t.Errorf("Expected type of nodes to be nodes, got %q", nodes[0].Type)
	}

	if nodes[0].Uptime <= 0 {
		t.Errorf("Expected nodes uptime to be greater than 0, got %q", nodes[0].Uptime)
	}
}
