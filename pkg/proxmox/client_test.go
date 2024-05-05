package proxmox

import (
	"testing"
)

const TestUsername = "root@pam"
const TestPassword = "vagrant"

func TestLogin(t *testing.T) {
	client, err := NewClient(DefaultHostURL, TestUsername, TestPassword)
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
	_, err := NewClient(DefaultHostURL, TestUsername, "wrong")
	if err == nil {
		t.Error("Expected authentication failure")
	}
}
