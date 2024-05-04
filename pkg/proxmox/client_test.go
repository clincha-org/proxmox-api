package proxmox

import (
	"testing"
)

func TestLogin(t *testing.T) {
	host := "https://localhost:8006"
	username := "root@pam"
	password := "vagrant"

	client, err := NewClient(&host, &username, &password)
	if err != nil {
		t.Error(err)
	}

	if client == nil {
		t.Errorf("Client was not initialised")
	}

	if client.Ticket.Data.Ticket == "" {
		t.Errorf("Expected ticket, got empty string")
	}

	if client.Ticket.Data.CSRFPreventionToken == "" {
		t.Errorf("Expected CSRFPreventionToken, got empty string")
	}

}

func TestIncorrectUsername(t *testing.T) {
	host := "https://localhost:8006"
	username := "thisuserdoesnotexist@pam"
	password := "thisuserdoesnotexist"

	_, err := NewClient(&host, &username, &password)
	if err == nil {
		t.Error("Expected authentication failure")
	}
}
