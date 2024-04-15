package proxmox

import (
	"errors"
	"log/slog"
	"net/http"
	"time"
)

const DefaultHostURL string = "https://localhost:8006"
const ApiPath string = "/api2/json/"
const AuthenticationTicketPath string = "access/ticket"

type Client struct {
	Host       string
	Username   string
	Password   string
	HTTPClient *http.Client
	Ticket     *Ticket
}

func NewClient(Host *string, Username *string, Password *string) (*Client, error) {
	client := Client{
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		Host: DefaultHostURL,
	}

	if Host != nil {
		client.Host = *Host
	} else {
		slog.Warn("Using default host as none was specified")
	}

	if Username == nil || Password == nil {
		return &client, errors.New("username and password are required")
	} else {
		client.Username = *Username
		client.Password = *Password
	}

	ticket, err := client.Login()
	if err != nil {
		return &client, err
	}

	client.Ticket = ticket

	return &client, nil
}
