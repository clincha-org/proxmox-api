package proxmox

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type TicketData struct {
	Data Ticket
}

type Ticket struct {
	Ticket              string
	CSRFPreventionToken string
	Cap                 Capabilities
}

type Capabilities struct {
	Storage                map[string]int8
	DataCenter             map[string]int8
	SoftwareDefinedNetwork map[string]int8
	VirtualMachines        map[string]int8
	Nodes                  map[string]int8
	Access                 map[string]int8
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

	var accessTicketResponse TicketData
	err = json.Unmarshal(body, &accessTicketResponse)
	if err != nil {
		return &ticket, err
	}

	ticket = accessTicketResponse.Data

	err = response.Body.Close()
	if err != nil {
		return &ticket, err
	}

	return &ticket, nil
}
