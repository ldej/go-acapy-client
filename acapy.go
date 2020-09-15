package acapy

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type Client struct {
	LedgerURL string
	AcapyURL  string

	HTTPClient http.Client
}

func NewClient(ledgerURL string, acapyURL string) *Client {
	return &Client{
		LedgerURL:  ledgerURL,
		AcapyURL:   acapyURL,
		HTTPClient: http.Client{},
	}
}

func (c *Client) post(url string, queryParam map[string]string, body interface{}, response interface{}) error {
	return c.request(http.MethodPost, url, queryParam, body, response)
}

func (c *Client) get(url string, queryParams map[string]string, response interface{}) error {
	return c.request(http.MethodGet, url, queryParams, nil, response)
}

func (c *Client) request(method string, url string, queryParams map[string]string, body interface{}, response interface{}) error {
	var input io.Reader
	var err error

	if body != nil {
		jsonInput, err := json.Marshal(body)
		if err != nil {
			return err
		}
		input = bytes.NewReader(jsonInput)
	}

	r, err := http.NewRequest(method, url, input)
	if err != nil {
		return err
	}
	r.Header.Add("Content-Type", "application/json")

	q := r.URL.Query()
	for k, v := range queryParams {
		if k != "" && v != "" {
			q.Add(k, v)
		}
	}
	r.URL.RawQuery = q.Encode()

	result, err := c.HTTPClient.Do(r)
	if err != nil {
		return err
	}

	err = json.NewDecoder(result.Body).Decode(response)
	if err != nil {
		return err
	}
	return nil
}
