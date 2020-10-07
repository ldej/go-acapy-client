package acapy

import (
	"fmt"
	"strconv"
	"strings"
)

type PresentationExchange struct {
	PresentationExchangeID   string                  `json:"presentation_exchange_id"`
	ConnectionID             string                  `json:"connection_id"`
	ThreadID                 string                  `json:"thread_id"`
	State                    string                  `json:"state"`
	Initiator                string                  `json:"initiator"`
	Role                     string                  `json:"role"`
	PresentationProposalDict PresentationProposalMap `json:"presentation_proposal_dict"`
	PresentationRequest      PresentationRequest     `json:"presentation_request"`
	PresentationRequestDict  PresentationRequestMap  `json:"presentation_request_dict"` // TODO ?
	Presentation             Presentation            `json:"presentation"`
	Verified                 string                  `json:"verified"`
	CreatedAt                string                  `json:"created_at"`
	UpdatedAt                string                  `json:"updated_at"`
	ErrorMsg                 string                  `json:"error_msg"`
	AutoPresent              bool                    `json:"auto_present"`
	Trace                    bool                    `json:"trace"`
}

type PresentationProposalMap struct {
	Type                 string              `json:"@type"` // did:sov:BzCbsNYhMrjHiqZDTUASHg;spec/present-proof/1.0/propose-presentation
	ID                   string              `json:"@id"`   // thread id
	Comment              string              `json:"comment"`
	PresentationProposal PresentationPreview `json:"presentation_proposal"`
}

// TODO
type Proof struct {
	Proofs []struct {
		PrimaryProof struct {
			EqProof struct {
				RevealedAttrs struct {
					Name string `json:"name"`
				} `json:"revealed_attrs"`
				APrime string `json:"a_prime"`
				E      string `json:"e"`
				V      string `json:"v"`
				M      struct {
					MasterSecret string `json:"master_secret"`
				} `json:"m"`
				M2 string `json:"m2"`
			} `json:"eq_proof"`
			GeProofs []interface{} `json:"ge_proofs"`
		} `json:"primary_proof"`
		NonRevocProof struct {
			XList struct {
				Rho              string `json:"rho"`
				R                string `json:"r"`
				RPrime           string `json:"r_prime"`
				RPrimePrime      string `json:"r_prime_prime"`
				RPrimePrimePrime string `json:"r_prime_prime_prime"`
				O                string `json:"o"`
				OPrime           string `json:"o_prime"`
				M                string `json:"m"`
				MPrime           string `json:"m_prime"`
				T                string `json:"t"`
				TPrime           string `json:"t_prime"`
				M2               string `json:"m2"`
				S                string `json:"s"`
				C                string `json:"c"`
			} `json:"x_list"`
			CList struct {
				E string `json:"e"`
				D string `json:"d"`
				A string `json:"a"`
				G string `json:"g"`
				W string `json:"w"`
				S string `json:"s"`
				U string `json:"u"`
			} `json:"c_list"`
		} `json:"non_revoc_proof"`
	} `json:"proofs"`
	AggregatedProof struct {
		CHash string  `json:"c_hash"`
		CList [][]int `json:"c_list"`
	} `json:"aggregated_proof"`
}

// TODO
type RequestedProof struct {
	RevealedAttrs      struct{} `json:"revealed_attrs"`
	RevealedAttrGroups struct {
		ZeroNameUUID struct {
			SubProofIndex int `json:"sub_proof_index"`
			Values        struct {
				Name struct {
					Raw     string `json:"raw"`
					Encoded string `json:"encoded"`
				} `json:"name"`
			} `json:"values"`
		} `json:"0_name_uuid"`
	} `json:"revealed_attr_groups"`
	SelfAttestedAttrs struct{} `json:"self_attested_attrs"`
	UnrevealedAttrs   struct{} `json:"unrevealed_attrs"`
	Predicates        struct{} `json:"predicates"`
}

type Presentation struct {
	Proof          Proof          `json:"proof"`
	RequestedProof RequestedProof `json:"requested_proof"`
	Identifiers    []struct {
		SchemaID  string `json:"schema_id"`
		CredDefID string `json:"cred_def_id"`
		RevRegID  string `json:"rev_reg_id"`
		Timestamp int    `json:"timestamp"`
	} `json:"identifiers"`
}

type PresentationRequestMap struct {
}

type PresentationProposalRequest struct {
	Comment             string              `json:"comment"`
	AutoPresent         bool                `json:"auto_present"`
	PresentationPreview PresentationPreview `json:"presentation_proposal"`
	ConnectionID        string              `json:"connection_id"`
	Trace               bool                `json:"trace"`
}

type PresentationAttribute struct {
	Name                   string `json:"name"`
	CredentialDefinitionID string `json:"cred_def_id"`
	MimeType               string `json:"mime-type"`
	Value                  string `json:"value"`
	Referent               string `json:"referent"`
}

type Predicate struct {
	Name                   string `json:"name"`
	CredentialDefinitionID string `json:"cred_def_id"`
	Predicate              string `json:"predicate"`
	Threshold              int    `json:"threshold"`
}

type PresentationPreview struct {
	Type       string                  `json:"@type"`
	Attributes []PresentationAttribute `json:"attributes"`
	Predicates []Predicate             `json:"predicates"`
}

func NewPresentationPreview(attributes []PresentationAttribute, predicates []Predicate) PresentationPreview {
	if attributes == nil {
		attributes = []PresentationAttribute{}
	}
	if predicates == nil {
		predicates = []Predicate{}
	}
	return PresentationPreview{
		Type:       "did:sov:BzCbsNYhMrjHiqZDTUASHg;spec/present-proof/1.0/presentation-preview",
		Attributes: attributes,
		Predicates: predicates,
	}
}

type PresentationRequest struct {
	Name                string                        `json:"name"`
	Version             string                        `json:"version"`
	Nonce               string                        `json:"nonce"`
	RequestedAttributes map[string]RequestedAttribute `json:"requested_attributes"`
	RequestedPredicates map[string]RequestedPredicate `json:"requested_predicates"`
}

type PresentationRequestRequest struct {
	Trace        bool         `json:"trace"`
	Comment      string       `json:"comment"`
	ConnectionID string       `json:"connection_id"`
	ProofRequest ProofRequest `json:"proof_request"`
}

type NonRevoked struct {
	From int64 `json:"from"`
	To   int64 `json:"to"`
}

type PredicateRestrictions struct {
	SchemaVersion          string `json:"schema_version"`
	CredentialDefinitionID string `json:"cred_def_id"`
	SchemaName             string `json:"schema_name"`
	SchemaIssuerDID        string `json:"schema_issuer_did"`
	IssuerDID              string `json:"issuer_did"`
	SchemaID               string `json:"schema_id"`
}

// TODO constructor
type RequestedPredicate struct {
	Restrictions PredicateRestrictions `json:"restrictions"`
	Name         string                `json:"name,omitempty"`  // XOR with Names
	Names        []string              `json:"names,omitempty"` // XOR with Name | Requires non-empty restrictions
	PType        string                `json:"p_type"`
	PValue       int                   `json:"p_value"`
	NonRevoked   NonRevoked            `json:"non_revoked"`
}

// TODO constructor
type RequestedAttribute struct {
	Restrictions []map[string]string `json:"restrictions"`    // valid key: cred_def_id
	Name         string              `json:"name,omitempty"`  // XOR with Names
	Names        []string            `json:"names,omitempty"` // XOR with Name | Requires non-empty restrictions
	NonRevoked   NonRevoked          `json:"non_revoked"`
}

// TODO constructor
type ProofRequest struct {
	Name                string                        `json:"name"`
	Nonce               string                        `json:"nonce"`
	RequestedPredicates map[string]RequestedPredicate `json:"requested_predicates"` // TODO cannot be nil
	RequestedAttributes map[string]RequestedAttribute `json:"requested_attributes"` // TODO cannot be nil
	Version             string                        `json:"version"`
	NonRevoked          NonRevoked                    `json:"non_revoked"`
}

type PresentationProofAttribute struct {
	Revealed     bool   `json:"revealed"`
	Timestamp    int64  `json:"timestamp"`
	CredentialID string `json:"cred_id"` // referent?
}

type PresentationProofPredicate struct {
	Timestamp    int64  `json:"timestamp"`
	CredentialID string `json:"cred_id"` // referent?
}

// TODO Create constructor to prevent nils
type PresentationProof struct {
	RequestedAttributes    map[string]PresentationProofAttribute `json:"requested_attributes"`     // TODO Cannot be nil
	RequestedPredicates    map[string]PresentationProofPredicate `json:"requested_predicates"`     // TODO Cannot be nil
	SelfAttestedAttributes map[string]string                     `json:"self_attested_attributes"` // TODO Cannot be nil
	Trace                  bool                                  `json:"trace"`
}

func (c *Client) SendPresentationProposal(request PresentationProposalRequest) (PresentationExchange, error) {
	var presentationExchange PresentationExchange
	err := c.post(fmt.Sprintf("%s/present-proof/send-proposal", c.ACApyURL), nil, request, &presentationExchange)
	if err != nil {
		return PresentationExchange{}, err
	}
	return presentationExchange, nil
}

func (c *Client) CreatePresentationRequest(request PresentationRequestRequest) (PresentationExchange, error) {
	var presentationExchange PresentationExchange
	err := c.post(fmt.Sprintf("%s/present-proof/create-request", c.ACApyURL), nil, request, &presentationExchange)
	if err != nil {
		return PresentationExchange{}, err
	}
	return presentationExchange, nil
}

func (c *Client) SendPresentationRequest(request PresentationRequestRequest) (PresentationExchange, error) {
	var presentationExchange PresentationExchange
	err := c.post(fmt.Sprintf("%s/present-proof/send-request", c.ACApyURL), nil, request, &presentationExchange)
	if err != nil {
		return PresentationExchange{}, err
	}
	return presentationExchange, nil
}

type QueryPresentationExchangeParams struct {
	ConnectionID string
	Role         string
	State        string
	ThreadID     string
}

func (c *Client) QueryPresentationExchange(params QueryPresentationExchangeParams) ([]PresentationExchange, error) {
	var result = struct {
		PresentationExchanges []PresentationExchange `json:"results"`
	}{}
	queryParams := map[string]string{
		"connection_id": params.ConnectionID,
		"role":          params.Role,
		"state":         params.State,
		"thread_id":     params.ThreadID,
	}
	err := c.get(fmt.Sprintf("%s/present-proof/records", c.ACApyURL), queryParams, &result)
	if err != nil {
		return nil, err
	}
	return result.PresentationExchanges, nil
}

func (c *Client) GetPresentationExchangeByID(presentationExchangeID string) (PresentationExchange, error) {
	var presentationExchange PresentationExchange
	err := c.get(fmt.Sprintf("%s/present-proof/records/%s", c.ACApyURL, presentationExchangeID), nil, &presentationExchange)
	if err != nil {
		return PresentationExchange{}, nil
	}
	return presentationExchange, nil
}

func (c *Client) SendPresentationRequestByID(presentationExchangeID string, request PresentationRequestRequest) (PresentationExchange, error) {
	var presentationExchange PresentationExchange
	err := c.post(fmt.Sprintf("%s/present-proof/records/%s/send-request", c.ACApyURL, presentationExchangeID), nil, request, &presentationExchange)
	if err != nil {
		return PresentationExchange{}, nil
	}
	return presentationExchange, nil
}

func (c *Client) SendPresentationByID(presentationExchangeID string, proof PresentationProof) (PresentationExchange, error) {
	var presentationExchange PresentationExchange
	err := c.post(fmt.Sprintf("%s/present-proof/records/%s/send-presentation", c.ACApyURL, presentationExchangeID), nil, proof, &presentationExchange)
	if err != nil {
		return PresentationExchange{}, nil
	}
	return presentationExchange, nil
}

func (c *Client) VerifyPresentationByID(presentationExchangeID string) (PresentationExchange, error) {
	var presentationExchange PresentationExchange
	err := c.post(fmt.Sprintf("%s/present-proof/records/%s/verify-presentation", c.ACApyURL, presentationExchangeID), nil, nil, &presentationExchange)
	if err != nil {
		return PresentationExchange{}, nil
	}
	return presentationExchange, nil
}

type PresentationCredential struct {
	CredentialInfo struct {
		Referent               string            `json:"referent"`
		Attrs                  map[string]string `json:"attrs"`
		SchemaID               string            `json:"schema_id"`
		CredentialDefinitionID string            `json:"cred_def_id"`
		RevocationRegistryID   string            `json:"rev_reg_id"`
		CredentialRevocationID string            `json:"cred_rev_id"`
	} `json:"cred_info"`
	Interval struct {
		From int `json:"from"`
		To   int `json:"to"`
	} `json:"interval"`
	PresentationReferents []string `json:"presentation_referents"`
}

func (c *Client) GetPresentationCredentialsByID(presentationExchangeID string, count int, extraQuery string, proofRequestReferents []string, start int) ([]PresentationCredential, error) {
	var result []PresentationCredential
	queryParams := map[string]string{
		"count":       strconv.Itoa(count),
		"extra_query": extraQuery,
		"start":       strconv.Itoa(start),
		"referent":    strings.Join(proofRequestReferents, ","),
	}
	err := c.get(fmt.Sprintf("%s/present-proof/records/%s", c.ACApyURL, presentationExchangeID), queryParams, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) RemovePresentationExchangeByID(presentationExchangeID string) error {
	return c.post(fmt.Sprintf("%s/present-proof/records/%s/remove", c.ACApyURL, presentationExchangeID), nil, nil, nil)
}
