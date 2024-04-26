package main

import (
	"fmt"
	"github.com/clincha-org/proxmox-api/pkg/proxmox"
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

	fmt.Println(networks[0])
}
