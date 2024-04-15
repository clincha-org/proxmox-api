package main

import (
	"github.com/clincha-org/proxmox-api/pkg/proxmox"
	"testing"
)

func TestLogin(t *testing.T) {
	host := "https://localhost:8006"
	username := "root@pam"
	password := "vagrant"

	client, err := proxmox.NewClient(&host, &username, &password)
	if err != nil {
		t.Error(err)
	}

	if client == nil {
		t.Errorf("Client was not initialised")
	}

	if client.Ticket.Ticket == "" {
		t.Errorf("Expected ticket, got empty string")
	}

	if client.Ticket.CSRFPreventionToken == "" {
		t.Errorf("Expected CSRFPreventionToken, got empty string")
	}

	for k, v := range client.Ticket.Cap.Storage {
		if v != 1 {
			t.Errorf("Expected root user to have %q privilage but got %d", k, v)
		}
	}
	for k, v := range client.Ticket.Cap.DataCenter {
		if v != 1 {
			t.Errorf("Expected root user to have %q privilage but got %d", k, v)
		}
	}
	for k, v := range client.Ticket.Cap.SoftwareDefinedNetwork {
		if v != 1 {
			t.Errorf("Expected root user to have %q privilage but got %d", k, v)
		}
	}
	for k, v := range client.Ticket.Cap.VirtualMachines {
		if v != 1 {
			t.Errorf("Expected root user to have %q privilage but got %d", k, v)
		}
	}
	for k, v := range client.Ticket.Cap.Nodes {
		if v != 1 {
			t.Errorf("Expected root user to have %q privilage but got %d", k, v)
		}
	}
	for k, v := range client.Ticket.Cap.Access {
		if v != 1 {
			t.Errorf("Expected root user to have %q privilage but got %d", k, v)
		}
	}
}

func TestIncorrectUsername(t *testing.T) {
	host := "https://localhost:8006"
	username := "thisuserdoesnotexist@pam"
	password := "thisuserdoesnotexist"

	_, err := proxmox.NewClient(&host, &username, &password)
	if err == nil {
		t.Error("Expected authentication failure")
	}
}
