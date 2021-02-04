package acapy

import (
	"errors"
	"fmt"
	"strings"
)

type Schema struct {
	Ver            string   `json:"ver"`
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Version        string   `json:"version"`
	AttributeNames []string `json:"attrNames"`
	SeqNo          int      `json:"seqNo"`
}

type schemaRequest struct {
	Version    string   `json:"schema_version"`
	Name       string   `json:"schema_name"`
	Attributes []string `json:"attributes"`
}

type schemaResponse struct {
	SchemaID string `json:"schema_id"`
	Schema   Schema `json:"schema"`
}

var ErrInvalidSchemaID = errors.New("invalid schema ID")

// SchemaIDToParts takes a schemaID, for example 6qnvgJtqwK44D8LFYnV5Yf:2:registration.dflow:1.0.0
// and returns the schema's issuer DID, the `ver`, the schema name and the schema version.
func SchemaIDToParts(schemaID string) (string, string, string, string, error) {
	parts := strings.Split(schemaID, ":")
	if len(parts) != 4 {
		return "", "", "", "", ErrInvalidSchemaID
	}
	return parts[0], parts[1], parts[2], parts[3], nil
}

func (c *Client) RegisterSchema(name string, version string, attributes []string) (Schema, error) {
	var request = schemaRequest{
		Name:       name,
		Version:    version,
		Attributes: attributes,
	}
	var response schemaResponse
	err := c.post("/schemas", nil, request, &response)
	if err != nil {
		return Schema{}, err
	}
	return response.Schema, err
}

type QuerySchemasParams struct {
	SchemaID        string `json:"schema_id"`
	SchemaIssuerDID string `json:"schema_issuer_did"`
	SchemaName      string `json:"schema_name"`
	SchemaVersion   string `json:"schema_version"`
}

func (c *Client) QuerySchemas(params QuerySchemasParams) ([]string, error) {
	type result struct {
		SchemaIDs []string `json:"schema_ids"`
	}
	var r result
	queryParams := map[string]string{
		"schema_id":         params.SchemaID,
		"schema_issuer_did": params.SchemaIssuerDID,
		"schema_name":       params.SchemaName,
		"schema_version":    params.SchemaVersion,
	}
	err := c.get("/schemas/created", queryParams, &r)
	if err != nil {
		return nil, err
	}
	return r.SchemaIDs, nil
}

func (c *Client) GetSchema(schemaID string) (Schema, error) {
	var schemaResponse schemaResponse
	err := c.get(fmt.Sprintf("/schemas/%s", schemaID), nil, &schemaResponse)
	if err != nil {
		return Schema{}, err
	}
	return schemaResponse.Schema, nil
}
