package proxmox

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
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

func (client *Client) Login() (*Ticket, error) {
	var ticket = Ticket{}

	authPayload := url.Values{}
	authPayload.Add("username", client.Username)
	authPayload.Add("password", client.Password)

	client.HTTPClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	response, err := client.HTTPClient.Post(
		client.Host+ApiPath+AuthenticationTicketPath,
		"application/x-www-form-urlencoded",
		strings.NewReader(authPayload.Encode()),
	)
	if err != nil {
		return &ticket, err
	}

	if response.StatusCode != http.StatusOK {
		return &ticket, errors.New(response.Status)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return &ticket, err
	}

	err = json.Unmarshal(body, &ticket)
	if err != nil {
		return &ticket, err
	}

	err = response.Body.Close()
	if err != nil {
		return &ticket, err
	}

	return &ticket, nil
}
