package acapy

// This thing is so big that it deserves it's own file
type CredentialExchangeRecord struct {
	CredentialExchangeID      string                    `json:"credential_exchange_id"`
	CredentialDefinitionID    string                    `json:"credential_definition_id"`
	ConnectionID              string                    `json:"connection_id"`
	ThreadID                  string                    `json:"thread_id"`
	ParentThreadID            string                    `json:"parent_thread_id"`
	SchemaID                  string                    `json:"schema_id"`
	RevocationID              string                    `json:"revocation_id"`
	RevocationRegistryID      string                    `json:"revoc_reg_id"`
	State                     string                    `json:"state"`
	CredentialOffer           CredentialOffer           `json:"credential_offer"`
	CredentialOfferMap        CredentialOfferMap        `json:"credential_offer_dict"`
	CredentialProposalMap     CredentialProposal        `json:"credential_proposal_dict"`
	CredentialRequest         CredentialRequest         `json:"credential_request"`
	CredentialRequestMetadata CredentialRequestMetadata `json:"credential_request_metadata"`
	Credential                Credential                `json:"credential"`
	RawCredential             RawCredential             `json:"raw_credential"`
	Role                      string                    `json:"role"`
	Initiator                 string                    `json:"initiator"`
	CreatedAt                 string                    `json:"created_at"`
	UpdatedAt                 string                    `json:"updated_at"`
	ErrorMessage              string                    `json:"error_msg"`
	Trace                     bool                      `json:"trace"`
	AutoOffer                 bool                      `json:"auto_offer"`
	AutoIssue                 bool                      `json:"auto_issue"`
	AutoRemove                bool                      `json:"auto_remove"`
}

type CredentialOffer struct {
	SchemaID               string `json:"schema_id"`
	CredentialDefinitionID string `json:"cred_def_id"`
	KeyCorrectnessProof    struct {
		C     string     `json:"c"`
		XzCap string     `json:"xz_cap"`
		XrCap [][]string `json:"xr_cap"`
	} `json:"key_correctness_proof"`
	Nonce string `json:"nonce"`
}

type OfferAttach struct {
	ID       string `json:"@id"` // libindy-cred-offer-0
	MimeType string `json:"mime-type"`
	Data     struct {
		Base64 string `json:"base64"`
	} `json:"data"`
}

type CredentialOfferMap struct {
	Type              string            `json:"@type"` // did:sov:BzCbsNYhMrjHiqZDTUASHg;spec/issue-credential/1.0/offer-credential
	ID                string            `json:"@id"`
	Thread            Thread            `json:"~thread"`
	CredentialPreview CredentialPreview `json:"credential_preview"`
	Comment           string            `json:"comment"`
	OffersAttach      []OfferAttach     `json:"offers~attach"`
}

type CredentialProposal struct {
	ID                     string            `json:"@id"`
	Type                   string            `json:"@type"` // did:sov:BzCbsNYhMrjHiqZDTUASHg;spec/issue-credential/1.0/propose-credential
	CredentialDefinitionID string            `json:"cred_def_id"`
	SchemaID               string            `json:"schema_id"`
	IssuerDID              string            `json:"issuer_did"`
	SchemaName             string            `json:"schema_name"`
	SchemaIssuerDID        string            `json:"schema_issuer_did"`
	SchemaVersion          string            `json:"schema_version"`
	Comment                string            `json:"comment"`
	CredentialPreview      CredentialPreview `json:"credential_proposal"`
}

type CredentialRequest struct {
	CredentialDefinitionID string `json:"cred_def_id"`
	ProverDID              string `json:"prover_did"`
	BlindedMs              struct {
		U                   string   `json:"u"`
		Ur                  string   `json:"ur"`
		HiddenAttributes    []string `json:"hidden_attributes"`
		CommittedAttributes struct {
		} `json:"committed_attributes"` // TODO
	} `json:"blinded_ms"`
	BlindedMsCorrectnessProof struct {
		C        string `json:"c"`
		VDashCap string `json:"v_dash_cap"`
		MCaps    struct {
			MasterSecret string `json:"master_secret"`
		} `json:"m_caps"`
		RCaps struct {
		} `json:"r_caps"`
	} `json:"blinded_ms_correctness_proof"`
	Nonce string `json:"nonce"`
}

type CredentialRequestMetadata struct {
	Description string `json:"description"`
}

type RawCredential struct {
	Description string `json:"raw_credential"`
}
