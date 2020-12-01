package acapy

import (
	"fmt"
)

func (c *Client) CreateCredentialDefinition(tag string, supportRevocation bool, revocationRegistrySize int, schemaID string) (string, error) {
	var request = struct {
		Tag                    string `json:"tag"`
		SupportRevocation      bool   `json:"support_revocation,omitempty"`
		RevocationRegistrySize int    `json:"revocation_registry_size,omitempty"`
		SchemaID               string `json:"schema_id"`
	}{
		Tag:                    tag,
		SupportRevocation:      supportRevocation,
		RevocationRegistrySize: revocationRegistrySize,
		SchemaID:               schemaID,
	}
	var response = struct {
		CredentialDefinitionID string `json:"credential_definition_id"`
	}{}
	err := c.post("/credential-definitions", nil, request, &response)
	if err != nil {
		return "", err
	}
	return response.CredentialDefinitionID, nil
}

type QueryCredentialDefinitionsParams struct {
	CredentialDefinitionID string `json:"cred_def_id"`
	IssuerDID              string `json:"issuer_did"`
	SchemaID               string `json:"schema_id"`
	SchemaIssuerDID        string `json:"schema_issuer_did"`
	SchemaName             string `json:"schema_name"`
	SchemaVersion          string `json:"schema_version"`
}

func (c *Client) QueryCredentialDefinitions(params QueryCredentialDefinitionsParams) ([]string, error) {
	var result = struct {
		CredentialDefinitionIDs []string `json:"credential_definition_ids"`
	}{}
	queryParams := map[string]string{
		"cred_def_id":       params.CredentialDefinitionID,
		"issuer_did":        params.IssuerDID,
		"schema_id":         params.SchemaID,
		"schema_issuer_did": params.SchemaIssuerDID,
		"schema_name":       params.SchemaName,
		"schema_version":    params.SchemaVersion,
	}
	err := c.get("/credential-definitions/created", queryParams, &result)
	if err != nil {
		return nil, err
	}
	return result.CredentialDefinitionIDs, nil
}

type R struct {
	Score        string `json:"score"`
	MasterSecret string `json:"master_secret"`
}
type Primary struct {
	N     string `json:"n"`
	S     string `json:"s"`
	R     R      `json:"r"`
	Rctxt string `json:"rctxt"`
	Z     string `json:"z"`
}
type Revocation struct {
	G      string `json:"g"`
	GDash  string `json:"g_dash"`
	H      string `json:"h"`
	H0     string `json:"h0"`
	H1     string `json:"h1"`
	H2     string `json:"h2"`
	Htilde string `json:"htilde"`
	HCap   string `json:"h_cap"`
	U      string `json:"u"`
	Pk     string `json:"pk"`
	Y      string `json:"y"`
}
type CredentialDefinitionValue struct {
	Primary    Primary    `json:"primary"`
	Revocation Revocation `json:"revocation"`
}
type CredentialDefinition struct {
	Version  string                    `json:"ver"`
	ID       string                    `json:"id"`
	SchemaID string                    `json:"schemaId"`
	Type     string                    `json:"type"`
	Tag      string                    `json:"tag"`
	Value    CredentialDefinitionValue `json:"value"`
}

func (c *Client) GetCredentialDefinition(credentialDefinitionID string) (CredentialDefinition, error) {
	var credentialDefinition CredentialDefinition
	err := c.get(fmt.Sprintf("/credential-definitions/%s", credentialDefinitionID), nil, &credentialDefinition)
	if err != nil {
		return CredentialDefinition{}, err
	}
	return credentialDefinition, nil
}
