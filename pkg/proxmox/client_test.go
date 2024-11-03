package proxmox

import (
	"log/slog"
	"testing"
)

const TestUsername = "root@pam"
const TestPassword = "vagrant"

func TestLogin(t *testing.T) {
	client, err := NewClient(DefaultHostURL, TestUsername, TestPassword, slog.LevelDebug)
	if err != nil {
		t.Error(err)
	}

	if client == nil {
		t.Fatalf("Client was not initialised")
	}

	if client.Ticket.Data.Ticket == "" {
		t.Fatalf("Expected ticket, got empty string")
	}

	if client.Ticket.Data.CSRFPreventionToken == "" {
		t.Fatalf("Expected CSRFPreventionToken, got empty string")
	}

}

func TestIncorrectUsername(t *testing.T) {
	_, err := NewClient(DefaultHostURL, TestUsername, "wrong", slog.LevelDebug)
	if err == nil {
		t.Error("Expected authentication failure")
	}
}

func TestGetVersion(t *testing.T) {
	client, err := NewClient(DefaultHostURL, TestUsername, TestPassword, slog.LevelDebug)
	if err != nil {
		t.Fatal(err)
	}

	version, err := client.GetVersion()
	if err != nil {
		t.Fatal(err)
	}

	if version.Version == "" {
		t.Errorf("Expected version, got empty string")
	}

	if version.MajorVersion == "" {
		t.Errorf("Expected major version, got empty string")
	}

	if version.Release == "" {
		t.Errorf("Expected release, got empty string")
	}

	if version.RepoID == "" {
		t.Errorf("Expected repo ID, got empty string")
	}

}
