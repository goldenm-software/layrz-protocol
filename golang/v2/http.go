package layrzprotocol

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type HttpScheme string

const (
	HTTP  HttpScheme = "http"
	HTTPS HttpScheme = "https"
)

type HttpComm struct {
	Scheme      HttpScheme
	Host        string
	Ident       string
	Passwd      string
	initialized bool
}

// New creates a new intance of LayrzProtocol using HTTP communication
func (p *HttpComm) New(scheme HttpScheme, host, ident, password string) {
	p.Scheme = scheme
	p.Host = host
	p.Ident = ident
	p.Passwd = password

	p.initialized = true
}

// Send sends a packet to the server
//
// Returns an error if the packet is invalid
// And, may return an Packet stored on an interface{}
func (p *HttpComm) Send(packet interface{}) (*interface{}, error) {
	if !p.initialized {
		return nil, errors.New("HttpComm not initialized")
	}

	data, err := parsePacketToString(packet)
	if err != nil {
		return nil, err
	}

	headers := http.Header{}
	headers.Add("Authorization", fmt.Sprintf("LayrzAuth %s;%s", p.Ident, p.Passwd))

	request := http.Request{
		Method: "POST",
		URL: &url.URL{
			Scheme: string(p.Scheme),
			Host:   p.Host,
			Path:   "/v2/message",
		},
		Header: headers,
		Body:   io.NopCloser(bytes.NewReader([]byte(*data))),
	}

	client := http.Client{}
	response, err := client.Do(&request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	body, err := io.ReadAll(io.Reader(response.Body))
	if err != nil {
		return nil, err
	}

	resp := string(body)
	return handleServerOutput(resp)
}

// Get new commands from the server
func (p *HttpComm) GetCommands() (*interface{}, error) {
	if !p.initialized {
		return nil, errors.New("HttpComm not initialized")
	}

	headers := http.Header{}
	headers.Add("Authorization", fmt.Sprintf("LayrzAuth %s;%s", p.Ident, p.Passwd))

	request := http.Request{
		Method: "GET",
		URL: &url.URL{
			Scheme: string(p.Scheme),
			Host:   p.Host,
			Path:   "/v2/commands",
		},
		Header: headers,
	}

	client := http.Client{}
	response, err := client.Do(&request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	body, err := io.ReadAll(io.Reader(response.Body))
	if err != nil {
		return nil, err
	}

	resp := string(body)
	return handleServerOutput(resp)
}
