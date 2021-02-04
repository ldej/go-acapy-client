package acapy

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type PresentationExchangeRecord struct {
	PresentationExchangeID   string                  `json:"presentation_exchange_id"`
	ConnectionID             string                  `json:"connection_id"`
	ThreadID                 string                  `json:"thread_id"`
	State                    string                  `json:"state"`
	Initiator                string                  `json:"initiator"`
	Role                     string                  `json:"role"`
	PresentationProposalDict PresentationProposalMap `json:"presentation_proposal_dict"`
	PresentationRequest      PresentationRequest     `json:"presentation_request"`
	PresentationRequestDict  struct{}                `json:"presentation_request_dict"` // TODO ?
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
	ThreadID             string              `json:"@id"`
	Comment              string              `json:"comment"`
	PresentationProposal PresentationPreview `json:"presentation_proposal"`
}

// TODO
// type Proof struct {
// 	Proofs []struct {
// 		PrimaryProof struct {
// 			EqProof struct {
// 				RevealedAttrs struct {
// 					Name string `json:"name"`
// 				} `json:"revealed_attrs"`
// 				APrime string `json:"a_prime"`
// 				E      string `json:"e"`
// 				V      string `json:"v"`
// 				M      struct {
// 					MasterSecret string `json:"master_secret"`
// 				} `json:"m"`
// 				M2 string `json:"m2"`
// 			} `json:"eq_proof"`
// 			GeProofs []interface{} `json:"ge_proofs"`
// 		} `json:"primary_proof"`
// 		NonRevocProof struct {
// 			XList struct {
// 				Rho              string `json:"rho"`
// 				R                string `json:"r"`
// 				RPrime           string `json:"r_prime"`
// 				RPrimePrime      string `json:"r_prime_prime"`
// 				RPrimePrimePrime string `json:"r_prime_prime_prime"`
// 				O                string `json:"o"`
// 				OPrime           string `json:"o_prime"`
// 				M                string `json:"m"`
// 				MPrime           string `json:"m_prime"`
// 				T                string `json:"t"`
// 				TPrime           string `json:"t_prime"`
// 				M2               string `json:"m2"`
// 				S                string `json:"s"`
// 				C                string `json:"c"`
// 			} `json:"x_list"`
// 			CList struct {
// 				E string `json:"e"`
// 				D string `json:"d"`
// 				A string `json:"a"`
// 				G string `json:"g"`
// 				W string `json:"w"`
// 				S string `json:"s"`
// 				U string `json:"u"`
// 			} `json:"c_list"`
// 		} `json:"non_revoc_proof"`
// 	} `json:"proofs"`
// 	AggregatedProof struct {
// 		CHash string  `json:"c_hash"`
// 		CList [][]int `json:"c_list"`
// 	} `json:"aggregated_proof"`
// }

// TODO
type RequestedProof struct {
	RevealedAttrs      struct{}                `json:"revealed_attrs"`
	RevealedAttrGroups map[string]RevealedAttr `json:"revealed_attr_groups"`
	SelfAttestedAttrs  struct{}                `json:"self_attested_attrs"`
	UnrevealedAttrs    struct{}                `json:"unrevealed_attrs"`
	Predicates         struct{}                `json:"predicates"`
}

type RevealedAttr struct {
	SubProofIndex int `json:"sub_proof_index"`
	Values        struct {
		Name struct {
			Raw     string `json:"raw"`
			Encoded string `json:"encoded"`
		} `json:"name"`
	} `json:"values"`
}

type Presentation struct {
	// Proof          Proof          `json:"proof"` // TODO
	RequestedProof RequestedProof `json:"requested_proof"`
	Identifiers    []struct {
		SchemaID               string `json:"schema_id"`
		CredentialDefinitionID string `json:"cred_def_id"`
		RevocationRegistryID   string `json:"rev_reg_id"`
		Timestamp              int    `json:"timestamp"`
	} `json:"identifiers"`
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

type PredicateType string

const (
	PredicateLT  PredicateType = "<"
	PredicateLTE PredicateType = "<="
	PredicateGT  PredicateType = ">"
	PredicateGTE PredicateType = ">="
)

type Predicate struct {
	Name                   string        `json:"name"`
	CredentialDefinitionID string        `json:"cred_def_id"`
	Predicate              PredicateType `json:"predicate"`
	Threshold              int           `json:"threshold"`
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

type Restrictions struct {
	SchemaVersion          string `json:"schema_version,omitempty"`
	CredentialDefinitionID string `json:"cred_def_id,omitempty"`
	SchemaName             string `json:"schema_name,omitempty"`
	SchemaIssuerDID        string `json:"schema_issuer_did,omitempty"`
	IssuerDID              string `json:"issuer_did,omitempty"`
	SchemaID               string `json:"schema_id,omitempty"`
	// TODO support `{"attr::attr1::value": "<some-value>"}`
}

func (r Restrictions) IsEmpty() bool {
	return r.SchemaVersion == "" &&
		r.CredentialDefinitionID == "" &&
		r.SchemaName == "" &&
		r.SchemaIssuerDID == "" &&
		r.IssuerDID == "" &&
		r.SchemaID == ""
}

func NewRequestedPredicate(
	restrictions *Restrictions,
	name string,
	names []string,
	ptype PredicateType,
	pvalue int,
	nonRevoked NonRevoked,
) (RequestedPredicate, error) {

	restrictionsSlice := make([]Restrictions, 0)

	if name != "" && len(names) > 0 {
		return RequestedPredicate{}, errors.New("use either 'name' or 'names', but not both")
	}
	if len(names) > 0 {
		if restrictions == nil || restrictions.IsEmpty() {
			return RequestedPredicate{}, errors.New("restrictions cannot be empty when using 'names'")
		}
		restrictionsSlice = append(restrictionsSlice, *restrictions)
	}

	return RequestedPredicate{
		Restrictions: restrictionsSlice,
		Name:         name,
		Names:        names,
		PType:        ptype,
		PValue:       pvalue,
		NonRevoked:   nonRevoked,
	}, nil
}

type RequestedPredicate struct {
	Restrictions []Restrictions `json:"restrictions"`    // Required when using Names, otherwise empty slice instead of nil
	Name         string         `json:"name,omitempty"`  // XOR with Names
	Names        []string       `json:"names,omitempty"` // XOR with Name | Requires non-empty restrictions
	PType        PredicateType  `json:"p_type"`
	PValue       int            `json:"p_value"`
	NonRevoked   NonRevoked     `json:"non_revoked"` // Optional
}

func NewRequestedAttribute(
	restrictions *Restrictions,
	name string,
	names []string,
	nonRevoked NonRevoked,
) (RequestedAttribute, error) {

	restrictionsSlice := make([]Restrictions, 0)

	if name != "" && len(names) > 0 {
		return RequestedAttribute{}, errors.New("use either 'name' or 'names', but not both")
	}
	if len(names) > 0 {
		if restrictions == nil || restrictions.IsEmpty() {
			return RequestedAttribute{}, errors.New("restrictions cannot be empty when using 'names'")
		}
		restrictionsSlice = append(restrictionsSlice, *restrictions)
	}

	return RequestedAttribute{
		Restrictions: restrictionsSlice,
		Name:         name,
		Names:        names,
		NonRevoked:   nonRevoked,
	}, nil
}

type RequestedAttribute struct {
	Restrictions []Restrictions `json:"restrictions"`    // Required when using Names, otherwise empty slice instead of nil
	Name         string         `json:"name,omitempty"`  // XOR with Names
	Names        []string       `json:"names,omitempty"` // XOR with Name | Requires non-empty restrictions
	NonRevoked   NonRevoked     `json:"non_revoked"`     // Optional
}

func NewProofRequest(
	name string,
	nonce string,
	requestedPredicates map[string]RequestedPredicate,
	requestedAttributes map[string]RequestedAttribute,
	version string,
	nonRevoked *NonRevoked,
) ProofRequest {

	if requestedPredicates == nil {
		requestedPredicates = map[string]RequestedPredicate{}
	}
	if requestedAttributes == nil {
		requestedAttributes = map[string]RequestedAttribute{}
	}
	return ProofRequest{
		Name:                name,
		Nonce:               nonce,
		RequestedPredicates: requestedPredicates,
		RequestedAttributes: requestedAttributes,
		Version:             version,
		NonRevoked:          nonRevoked,
	}
}

type ProofRequest struct {
	Name                string                        `json:"name"`
	Nonce               string                        `json:"nonce"`                // TODO what is this nonce
	RequestedPredicates map[string]RequestedPredicate `json:"requested_predicates"` // cannot be nil
	RequestedAttributes map[string]RequestedAttribute `json:"requested_attributes"` // cannot be nil
	Version             string                        `json:"version"`
	NonRevoked          *NonRevoked                   `json:"non_revoked,omitempty"`
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

func NewPresentationProof(
	requestedAttributes map[string]PresentationProofAttribute,
	requestedPredicates map[string]PresentationProofPredicate,
	selfAttestedAttributes map[string]string,
) PresentationProof {
	if requestedAttributes == nil {
		requestedAttributes = map[string]PresentationProofAttribute{}
	}
	if requestedPredicates == nil {
		requestedPredicates = map[string]PresentationProofPredicate{}
	}
	if selfAttestedAttributes == nil {
		selfAttestedAttributes = map[string]string{}
	}
	return PresentationProof{
		RequestedAttributes:    requestedAttributes,
		RequestedPredicates:    requestedPredicates,
		SelfAttestedAttributes: selfAttestedAttributes,
	}
}

type PresentationProof struct {
	RequestedAttributes    map[string]PresentationProofAttribute `json:"requested_attributes"`
	RequestedPredicates    map[string]PresentationProofPredicate `json:"requested_predicates"`
	SelfAttestedAttributes map[string]string                     `json:"self_attested_attributes"`
	Trace                  bool                                  `json:"trace"`
}

func (c *Client) SendPresentationProposal(request PresentationProposalRequest) (PresentationExchangeRecord, error) {
	var presentationExchange PresentationExchangeRecord
	err := c.post("/present-proof/send-proposal", nil, request, &presentationExchange)
	if err != nil {
		return PresentationExchangeRecord{}, err
	}
	return presentationExchange, nil
}

func (c *Client) CreatePresentationRequest(request PresentationRequestRequest) (PresentationExchangeRecord, error) {
	var presentationExchange PresentationExchangeRecord
	err := c.post("/present-proof/create-request", nil, request, &presentationExchange)
	if err != nil {
		return PresentationExchangeRecord{}, err
	}
	return presentationExchange, nil
}

func (c *Client) SendPresentationRequest(request PresentationRequestRequest) (PresentationExchangeRecord, error) {
	var presentationExchange PresentationExchangeRecord
	err := c.post("/present-proof/send-request", nil, request, &presentationExchange)
	if err != nil {
		return PresentationExchangeRecord{}, err
	}
	return presentationExchange, nil
}

type QueryPresentationExchangeParams struct {
	ConnectionID string
	Role         string
	State        string
	ThreadID     string
}

func (c *Client) QueryPresentationExchange(params QueryPresentationExchangeParams) ([]PresentationExchangeRecord, error) {
	var result = struct {
		PresentationExchanges []PresentationExchangeRecord `json:"results"`
	}{}
	queryParams := map[string]string{
		"connection_id": params.ConnectionID,
		"role":          params.Role,
		"state":         params.State,
		"thread_id":     params.ThreadID,
	}
	err := c.get("/present-proof/records", queryParams, &result)
	if err != nil {
		return nil, err
	}
	return result.PresentationExchanges, nil
}

func (c *Client) GetPresentationExchangeByID(presentationExchangeID string) (PresentationExchangeRecord, error) {
	var presentationExchange PresentationExchangeRecord
	err := c.get(fmt.Sprintf("/present-proof/records/%s", presentationExchangeID), nil, &presentationExchange)
	if err != nil {
		return PresentationExchangeRecord{}, nil
	}
	return presentationExchange, nil
}

func (c *Client) SendPresentationRequestByID(presentationExchangeID string, request PresentationRequestRequest) (PresentationExchangeRecord, error) {
	var presentationExchange PresentationExchangeRecord
	err := c.post(fmt.Sprintf("/present-proof/records/%s/send-request", presentationExchangeID), nil, request, &presentationExchange)
	if err != nil {
		return PresentationExchangeRecord{}, nil
	}
	return presentationExchange, nil
}

func (c *Client) SendPresentationByID(presentationExchangeID string, proof PresentationProof) (PresentationExchangeRecord, error) {
	var presentationExchange PresentationExchangeRecord
	err := c.post(fmt.Sprintf("/present-proof/records/%s/send-presentation", presentationExchangeID), nil, proof, &presentationExchange)
	if err != nil {
		return PresentationExchangeRecord{}, nil
	}
	return presentationExchange, nil
}

func (c *Client) VerifyPresentationByID(presentationExchangeID string) (PresentationExchangeRecord, error) {
	var presentationExchange PresentationExchangeRecord
	err := c.post(fmt.Sprintf("/present-proof/records/%s/verify-presentation", presentationExchangeID), nil, nil, &presentationExchange)
	if err != nil {
		return PresentationExchangeRecord{}, nil
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

func (c *Client) GetPresentationCredentialsByID(presentationExchangeID string, count int, wql string, proofRequestReferents []string, start int) ([]PresentationCredential, error) {
	var result []PresentationCredential
	queryParams := map[string]string{
		"extra_query": wql,
		"referent":    strings.Join(proofRequestReferents, ","),
	}
	if count > 0 {
		queryParams["count"] = strconv.Itoa(count)
	}
	if start > 0 {
		queryParams["start"] = strconv.Itoa(start)
	}
	err := c.get(fmt.Sprintf("/present-proof/records/%s/credentials", presentationExchangeID), queryParams, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) RemovePresentationExchangeByID(presentationExchangeID string) error {
	return c.delete(fmt.Sprintf("/present-proof/records/%s", presentationExchangeID))
}
