package acapy

import (
	"fmt"
	"strconv"
)

type RevocationRegistry struct {
	UpdatedAt            string                       `json:"updated_at"`
	Type                 string                       `json:"revoc_def_type"`
	PendingPub           []interface{}                `json:"pending_pub"`
	CreatedAt            string                       `json:"created_at"`
	Tag                  string                       `json:"tag"`
	RecordID             string                       `json:"record_id"`
	CredDefID            string                       `json:"cred_def_id"`
	TailsHash            string                       `json:"tails_hash"`
	MaxCredNum           int                          `json:"max_cred_num"`
	State                string                       `json:"state"`
	IssuerDid            string                       `json:"issuer_did"`
	Definition           RevocationRegistryDefinition `json:"revoc_reg_def"`
	TailsLocalPath       string                       `json:"tails_local_path"`
	Entry                RevocationRegistryEntry      `json:"revoc_reg_entry"`
	RevocationRegistryID string                       `json:"revoc_reg_id"`
}
type AccumKey struct {
	Z string `json:"z"`
}
type PublicKeys struct {
	AccumKey AccumKey `json:"accumKey"`
}
type RevocationRegistryDefinitionValue struct {
	IssuanceType  string     `json:"issuanceType"`
	MaxCredNum    int        `json:"maxCredNum"`
	PublicKeys    PublicKeys `json:"publicKeys"`
	TailsHash     string     `json:"tailsHash"`
	TailsLocation string     `json:"tailsLocation"`
}
type RevocationRegistryDefinition struct {
	Ver          string                            `json:"ver"`
	ID           string                            `json:"id"`
	RevocDefType string                            `json:"revocDefType"`
	Tag          string                            `json:"tag"`
	CredDefID    string                            `json:"credDefId"`
	Value        RevocationRegistryDefinitionValue `json:"value"`
}
type Value struct {
	Accum string `json:"accum"`
}
type RevocationRegistryEntry struct {
	Ver   string                            `json:"ver"`
	Value RevocationRegistryDefinitionValue `json:"value"`
}

func (c *Client) CreateRevocationRegistry(credentialDefinitionID string, maxCredNum int) (RevocationRegistry, error) {
	var result = struct {
		RevocationRegistry RevocationRegistry `json:"result"`
	}{}
	var body = struct {
		CredentialDefinitionID string `json:"credential_definition_id"`
		MaxCredNum             int    `json:"max_cred_num"`
	}{
		CredentialDefinitionID: credentialDefinitionID,
		MaxCredNum:             maxCredNum,
	}
	err := c.post("/revocation/create-registry", nil, body, &result)
	if err != nil {
		return RevocationRegistry{}, err
	}
	return result.RevocationRegistry, nil
}

func (c *Client) QueryRevocationRegistries(credentialDefinitionID string, state string) ([]string, error) {
	var queryParams = map[string]string{
		"cred_def_id": credentialDefinitionID,
		"state":       state,
	}
	var result = struct {
		RevocationRegistryIDs []string `json:"rev_reg_ids"`
	}{}
	err := c.get("/revocation/registries/created", queryParams, &result)
	if err != nil {
		return nil, err
	}
	return result.RevocationRegistryIDs, nil
}

func (c *Client) GetRevocationRegistry(revocationRegistryID string) (RevocationRegistry, error) {
	var result = struct {
		RevocationRegistry RevocationRegistry `json:"result"`
	}{}
	err := c.get(fmt.Sprintf("/revocation/registry/%s", revocationRegistryID), nil, &result)
	if err != nil {
		return RevocationRegistry{}, err
	}
	return result.RevocationRegistry, nil
}

func (c *Client) UpdateRevocationRegistryTailsURI(revocationRegistryID string, tailsPublicURI string) (RevocationRegistry, error) {
	var result = struct {
		RevocationRegistry RevocationRegistry `json:"result"`
	}{}
	var body = struct {
		TailsPublicURI string `json:"tails_public_uri"`
	}{
		TailsPublicURI: tailsPublicURI,
	}
	err := c.patch(fmt.Sprintf("/revocation/registry/%s", revocationRegistryID), nil, body, &result)
	if err != nil {
		return RevocationRegistry{}, err
	}
	return result.RevocationRegistry, nil
}

func (c *Client) GetActiveRevocationRegistry(credentialDefinitionID string) (RevocationRegistry, error) {
	var revocationRegistry RevocationRegistry
	err := c.get(fmt.Sprintf("/revocation/active-registry/%s", credentialDefinitionID), nil, &revocationRegistry)
	if err != nil {
		return RevocationRegistry{}, err
	}
	return revocationRegistry, nil
}

func (c *Client) DownloadRegistryTailsFile(revocationRegistryID string) ([]byte, error) {
	tailsFile, err := c.getFile(fmt.Sprintf("/revocation/registry/%s/tails-file", revocationRegistryID))
	if err != nil {
		return nil, err
	}
	return tailsFile, nil
}

func (c *Client) UploadRegistryTailsFile(revocationRegistryID string) error {
	return c.put(fmt.Sprintf("/revocation/registry/%s/tails-file", revocationRegistryID))
}

func (c *Client) PublishRevocationRegistryDefinition(revocationRegistryID string) (RevocationRegistry, error) {
	var result = struct {
		RevocationRegistry RevocationRegistry `json:"result"`
	}{}
	err := c.post(fmt.Sprintf("/revocation/registry/%s/definition", revocationRegistryID), nil, nil, &result)
	if err != nil {
		return RevocationRegistry{}, err
	}
	return result.RevocationRegistry, nil
}

func (c *Client) PublishRevocationRegistryEntry(revocationRegistryID string) (RevocationRegistry, error) {
	var result = struct {
		RevocationRegistry RevocationRegistry `json:"result"`
	}{}
	err := c.post(fmt.Sprintf("/revocation/registry/%s/entry", revocationRegistryID), nil, nil, &result)
	if err != nil {
		return RevocationRegistry{}, err
	}
	return result.RevocationRegistry, nil
}

func (c *Client) SetRevocationRegistryState(revocationRegistryID string, state string) (RevocationRegistry, error) {
	var result = struct {
		RevocationRegistry RevocationRegistry `json:"result"`
	}{}
	var queryParams = map[string]string{
		"state": state,
	}
	err := c.patch(fmt.Sprintf("/revocation/registry/%s/set-state", revocationRegistryID), queryParams, nil, &result)
	if err != nil {
		return RevocationRegistry{}, err
	}
	return result.RevocationRegistry, nil
}

func (c *Client) RevokeIssuedCredential(credentialRevocationID string, revocationRegistryID string, publish bool) error {
	queryParams := map[string]string{
		"cred_rev_id": credentialRevocationID,
		"rev_reg_id":  revocationRegistryID,
		"publish":     strconv.FormatBool(publish),
	}
	return c.post("/revocation/revoke", queryParams, nil, nil)
}

// A map from revocation registry identifier to credential revocation identifiers
// For example:
// map[string][]string{
// 	"6i7GFi2cDx524ZNfxmGWcp:4:6i7GFi2cDx524ZNfxmGWcp:3:CL:165:default:CL_ACCUM:159875bc-a5c7-4d51-b5c0-b4782a01fb94": ["1"].
// }
type PendingRevocations map[string][]string

// PublishRevocations
// Pass nil in case you want to publish all pending revocations
func (c *Client) PublishRevocations(revocations PendingRevocations) error {
	var body = struct {
		Body PendingRevocations `json:"rrid2crid"`
	}{
		Body: revocations,
	}
	return c.post("/revocation/publish-revocations", nil, body, nil)
}

// ClearPendingRevocations
// Pass nil in case you want to clear all pending revocations
func (c *Client) ClearPendingRevocations(revocations PendingRevocations) (PendingRevocations, error) {
	var result = struct {
		Result PendingRevocations `json:"rrid2crid"`
	}{}
	var body = struct {
		Body map[string][]string `json:"purge"`
	}{
		Body: revocations,
	}

	err := c.post("/revocation/clear-pending-revocations", nil, body, &result)
	if err != nil {
		return nil, err
	}
	return result.Result, nil
}

func (c *Client) GetCredentialRevocationStatus(credentialExchangeID string, credentialRevocationID string, revocationRegistryID string) error {
	// TODO
	return nil
}

func (c *Client) GetNumberOfIssuedCredentials(revocationRegistryID string) (int, error) {
	// TODO
	return 0, nil
}
