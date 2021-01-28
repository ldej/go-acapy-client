package acapy

import (
	"encoding/json"
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

// wql: https://ldej.nl/post/becoming-a-hyperledger-aries-developer-part-5-issue-credentials#wql-some-query-language
// https://github.com/hyperledger/aries-cloudagent-python/blob/master/aries_cloudagent/storage/basic.py#L135
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

// TODO from/to query params
func (c *Client) IsCredentialRevoked(credentialID string) (bool, error) {
	var result = struct {
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

func (c *Client) FindMatchingCredentials(request PresentationRequest) (map[string]PresentationProofAttribute, error) {

	requestedAttributes := map[string]PresentationProofAttribute{}

	for attrName, attr := range request.RequestedAttributes {
		restrictions, err := json.Marshal(attr.Restrictions[0])
		if err != nil {
			return nil, err
		}
		credentials, err := c.GetCredentials(10, 0, string(restrictions))

		if err != nil {
			return nil, err
		}

		if len(credentials) == 0 {
			return nil, fmt.Errorf("no credentials found for %s", attrName)
		} else if len(credentials) > 1 {
			return nil, fmt.Errorf("multiple credentials found for %s", attrName)
		}

		if containsAllAttributes(credentials[0], attr.Names) {
			requestedAttributes[attrName] = PresentationProofAttribute{
				Revealed: true,
				//Timestamp:    time.Now().Unix(),
				CredentialID: credentials[0].Referent,
			}
		}
	}
	return requestedAttributes, nil
}

func containsAllAttributes(credential Credential, attrs []string) bool {
	for _, attr := range attrs {
		if value, found := credential.Attributes[attr]; !found || value == "" {
			return false
		}
	}
	return true
}
