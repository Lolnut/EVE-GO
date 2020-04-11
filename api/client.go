package api

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const defaultURI = "https://esi.evetech.net"
const defaultVersion = "/latest"

type Client struct {
	baseURI string
	httpClient *http.Client
}

func NewClient(options ...func(*Client)) *Client {
	client := &Client{
		httpClient: &http.Client{},
	}

	for _, option := range options {
		option(client)
	}

	if client.baseURI == "" {
		client.baseURI = fmt.Sprintf("%s%s", defaultURI, defaultVersion)
	}
	return client
}

func (c *Client) getBytes(ctx context.Context, r *http.Request) ([]byte, error) {
	resp, err := c.httpClient.Do(r.WithContext(ctx))
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("%d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}
	return ioutil.ReadAll(resp.Body)
}

func (c *Client) prepareAddress(endpoint string, params map[string]string) (string, error) {
	address, err := url.Parse(c.baseURI + endpoint)
	if err != nil {
		return "", err
	}

	q := address.Query()
	for param, value := range params {
		q.Add(param, value)
	}
	address.RawQuery = q.Encode()
	return address.String(), nil
}

func (c *Client) Get(ctx context.Context, endpoint string, params map[string]string) ([]byte, error) {

	address, err := c.prepareAddress(endpoint, params)
	if err != nil {
		return []byte{}, err
	}

	req, err := http.NewRequest("GET", address, nil)
	if err != nil {
		return []byte{}, err
	}

	data, err := c.getBytes(ctx, req)
	if err != nil {
		return []byte{}, err
	}
	return data, nil
}

