package acapy

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Client struct {
	acapyURL                   string
	apiKey                     string
	tracing                    bool
	preserveExchangeRecords    bool
	autoRespondCredentialOffer bool
	HTTPClient                 http.Client
}

func NewClient(acapyURL string) *Client {
	return &Client{
		acapyURL:   strings.TrimRight(acapyURL, "/"),
		HTTPClient: http.Client{},
	}
}

func (c *Client) SetAPIKey(apiKey string) *Client {
	c.apiKey = apiKey
	return c
}

func (c *Client) EnableTracing() *Client {
	c.tracing = true
	return c
}

func (c *Client) DisableTracing() *Client {
	c.tracing = false
	return c
}

func (c *Client) PreserveExchangeRecords() *Client {
	c.preserveExchangeRecords = true
	return c
}

func (c *Client) AutoRespondCredentialOffer() *Client {
	c.autoRespondCredentialOffer = true
	return c
}

func (c *Client) post(path string, queryParam map[string]string, body interface{}, response interface{}) error {
	return c.request(http.MethodPost, c.acapyURL+path, queryParam, body, response)
}

func (c *Client) get(path string, queryParams map[string]string, response interface{}) error {
	return c.request(http.MethodGet, c.acapyURL+path, queryParams, nil, response)
}

func (c *Client) patch(path string, queryParams map[string]string, body interface{}, response interface{}) error {
	return c.request(http.MethodPatch, c.acapyURL+path, queryParams, body, response)
}

func (c *Client) put(path string) error {
	return c.request(http.MethodPut, c.acapyURL+path, nil, nil, nil)
}

func (c *Client) delete(url string) error {
	return c.request(http.MethodDelete, url, nil, nil, nil)
}

func (c *Client) request(method string, url string, queryParams map[string]string, body interface{}, responseObject interface{}) error {
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
	if c.apiKey != "" {
		r.Header.Add("X-API-KEY", c.apiKey)
	}
	r.Header.Add("Content-Type", "application/json")

	q := r.URL.Query()
	for k, v := range queryParams {
		if k != "" && v != "" {
			q.Add(k, v)
		}
	}
	r.URL.RawQuery = q.Encode()

	response, err := c.HTTPClient.Do(r)
	if err != nil || response.StatusCode >= 300 {
		if response != nil {
			log.Printf("Request failed: %s", response.Status)
			if body, err := ioutil.ReadAll(response.Body); err != nil {
				log.Printf("Response body: %s", body)
			}
		}
		return err
	}

	if responseObject != nil {
		err = json.NewDecoder(response.Body).Decode(responseObject)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) getFile(url string) ([]byte, error) {
	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	response, err := c.HTTPClient.Do(r)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	return ioutil.ReadAll(response.Body)
}
