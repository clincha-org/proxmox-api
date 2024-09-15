package proxmox

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
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

func NewClient(Host string, Username string, Password string, LogLevel slog.Level) (*Client, error) {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: LogLevel}))
	slog.SetDefault(logger)

	client := Client{
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		Host: DefaultHostURL,
	}

	if Host == "" {
		return &client, fmt.Errorf("NewClient-Host: %w", errors.New("host is required"))
	}
	client.Host = Host

	if Username == "" || Password == "" {
		return &client, fmt.Errorf("NewClient-Username-Password: %w", errors.New("username and password are required"))
	} else {
		client.Username = Username
		client.Password = Password
	}

	ticket, err := client.Login()
	if err != nil {
		return &client, fmt.Errorf("NewClient-Login: %w", err)
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
		return &ticket, fmt.Errorf("Login-do-request: %w", err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return &ticket, fmt.Errorf("Login-read-response: %w", err)
	}

	err = response.Body.Close()
	if err != nil {
		return &ticket, fmt.Errorf("Login-close-response: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		return &ticket, fmt.Errorf("Login-status-error: %s %s", response.Status, body)
	}

	err = json.Unmarshal(body, &ticket)
	if err != nil {
		return &ticket, fmt.Errorf("Login-unmarshal-response: %w", err)
	}

	return &ticket, nil
}

func (client *Client) MakeRESTRequest(method string, path string, body *bytes.Buffer) ([]byte, error) {
	request := &http.Request{}
	err := error(nil)

	// We need to check if body is nil before creating the request
	// because passing a nil body of type bytes.Buffer to http.NewRequest will cause a panic
	// https://go.dev/doc/faq#nil_error
	if body == nil {
		request, err = http.NewRequest(method, path, nil)
		if err != nil {
			return nil, fmt.Errorf("MakeRESTRequest-new-request: %w", err)
		}
	} else {
		request, err = http.NewRequest(method, path, body)
		if err != nil {
			return nil, fmt.Errorf("MakeRESTRequest-new-request-with-body: %w", err)
		}
		request.Header.Set("Content-Type", "application/json")
	}

	request.AddCookie(&http.Cookie{Name: "PVEAuthCookie", Value: client.Ticket.Data.Ticket})
	request.Header.Set("CSRFPreventionToken", client.Ticket.Data.CSRFPreventionToken)

	response, err := client.HTTPClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("MakeRESTRequest-do-request: %w", err)
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("MakeRESTRequest-read-response: %w", err)
	}

	err = response.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("MakeRESTRequest-close-response: %w", err)
	}

	slog.Debug("api-response", "method", "MakeRESTRequest", "rest_method", method, "path", path, "status", response.Status, "response", string(responseBody))

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GetVM-status-error: %s %s", response.Status, body)
	}

	return responseBody, nil
}
