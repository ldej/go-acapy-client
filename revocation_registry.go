package acapy

import (
	"fmt"
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
	err := c.post(fmt.Sprintf("%s/revocation/create-registry", c.ACApyURL), nil, body, &result)
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
	err := c.get(fmt.Sprintf("%s/revocation/registries/created", c.ACApyURL), queryParams, &result)
	if err != nil {
		return nil, err
	}
	return result.RevocationRegistryIDs, nil
}

func (c *Client) GetRevocationRegistry(revocationRegistryID string) (RevocationRegistry, error) {
	var revocationRegistry RevocationRegistry
	err := c.get(fmt.Sprintf("%s/revocation/registry/%s", c.ACApyURL, revocationRegistryID), nil, &revocationRegistry)
	if err != nil {
		return RevocationRegistry{}, err
	}
	return revocationRegistry, nil
}

func (c *Client) UpdateRevocationRegistryTailsURI(revocationRegistryID string, tailsPublicURI string) (RevocationRegistry, error) {
	var revocationRegistry RevocationRegistry
	var body = struct {
		TailsPublicURI string `json:"tails_public_uri"`
	}{
		TailsPublicURI: tailsPublicURI,
	}
	err := c.patch(fmt.Sprintf("%s/revocation/registry/%s", c.ACApyURL, revocationRegistryID), nil, body, &revocationRegistry)
	if err != nil {
		return RevocationRegistry{}, err
	}
	return revocationRegistry, nil
}

func (c *Client) GetActiveRevocationRegistry(credentialDefinitionID string) (RevocationRegistry, error) {
	var revocationRegistry RevocationRegistry
	err := c.get(fmt.Sprintf("%s/revocation/active-registry/%s", c.ACApyURL, credentialDefinitionID), nil, &revocationRegistry)
	if err != nil {
		return RevocationRegistry{}, err
	}
	return revocationRegistry, nil
}

func (c *Client) DownloadRegistryTailsFile(revocationRegistryID string) ([]byte, error) {
	tailsFile, err := c.getFile(fmt.Sprintf("%s/revocation/registry/%s/tails_file", c.ACApyURL, revocationRegistryID))
	if err != nil {
		return nil, err
	}
	return tailsFile, nil
}

func (c *Client) UploadRegistryTailsFile(revocationRegistryID string, tailsFile []byte) error {
	// TODO
	// return c.put_file(fmt.Sprintf("%s/revocation/registry/%s/tails_file", c.ACApyURL, revocationRegistryID), nil, tailsFile)
	return nil
}

func (c *Client) PublishRevocationRegistryDefinition(revocationRegistryID string) (RevocationRegistry, error) {
	var revocationRegistry RevocationRegistry
	err := c.post(fmt.Sprintf("%s/revocation/registry/%s/definition", c.ACApyURL, revocationRegistryID), nil, nil, &revocationRegistry)
	if err != nil {
		return RevocationRegistry{}, err
	}
	return revocationRegistry, nil
}

func (c *Client) PublishRevocationRegistryEntry(revocationRegistryID string) error {
	// TODO
	return nil
}

func (c *Client) SetRevocationRegistryState(revocationRegistryID string, state string) error {
	// TODO
	return nil
}
