package acapy

import (
	"fmt"
	"strconv"
)

type Credential struct {
	Referent               string            `json:"referent"` // Also know as CredentialID
	CredentialDefinitionID string            `json:"cred_def_id"`
	CredentialRevokeID     string            `json:"cred_rev_id"`
	SchemaID               string            `json:"schema_id"`
	RevokeRegistryID       string            `json:"rev_reg_id"`
	Attributes             map[string]string `json:"attrs"`
}

func (c *Client) GetCredentials(max int, index int, wql string) ([]Credential, error) {
	var results struct {
		Credentials []Credential `json:"results"`
	}
	queryParams := map[string]string{
		"max":   strconv.Itoa(max),
		"index": strconv.Itoa(index),
		"wql":   wql,
	}
	err := c.get("/credentials", queryParams, &results)
	if err != nil {
		return nil, err
	}
	return results.Credentials, nil
}

func (c *Client) GetCredential(credentialID string) (Credential, error) {
	var credential Credential
	err := c.get(fmt.Sprintf("/credential/%s", credentialID), nil, &credential)
	if err != nil {
		return Credential{}, err
	}
	return credential, nil
}

func (c *Client) IsCredentialRevoked(credentialID string) (bool, error) {
	var result = struct{
		Revoked bool `json:"revoked"`
	}{}
	err := c.get(fmt.Sprintf("/credential/revoked/%s", credentialID), nil, &result)
	if err != nil {
		return false, err
	}
	return result.Revoked, nil
}

func (c *Client) CredentialMimeTypes(credentialID string) (map[string]string, error) {
	var result map[string]string
	err := c.get(fmt.Sprintf("/credential/mime-types/%s", credentialID), nil, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) RemoveCredential(credentialID string) error {
	return c.delete(fmt.Sprintf("/credential/%s", credentialID))
}
