package acapy

// This thing is so big that it deserves it's own file
type CredentialExchangeRecordV2 struct {
	CredentialExchangeID       string                      `json:"cred_ex_id"`
	CredentialExchangeIDStored string                      `json:"cred_id_stored"`
	ConnectionID               string                      `json:"conn_id"`
	ThreadID                   string                      `json:"thread_id"`
	ParentThreadID             string                      `json:"parent_thread_id"`
	State                      string                      `json:"state"`
	CredentialPreview          CredentialPreviewV2         `json:"cred_preview"`
	CredentialOffer            CredentialOfferV2           `json:"cred_offer"`
	CredentialProposal         CredentialProposalV2        `json:"cred_proposal"`
	CredentialRequest          CredentialRequestV2         `json:"cred_request"`
	CredentialRequestMetadata  CredentialRequestMetadataV2 `json:"cred_request_metadata"`
	CredentialIssue            CredentialIssue             `json:"cred_issue"`
	Role                       string                      `json:"role"`
	Initiator                  string                      `json:"initiator"`
	CreatedAt                  string                      `json:"created_at"`
	UpdatedAt                  string                      `json:"updated_at"`
	ErrorMessage               string                      `json:"error_msg"`
	Trace                      bool                        `json:"trace"`
	AutoOffer                  bool                        `json:"auto_offer"`
	AutoIssue                  bool                        `json:"auto_issue"`
	AutoRemove                 bool                        `json:"auto_remove"`
}

type CredentialExchangeRecordResult struct {
	CredentialExchangeRecord CredentialExchangeRecordV2 `json:"cred_ex_record"`
	DIF                      CredentialExchangeDIF      `json:"dif"`
	Indy                     CredentialExchangeIndy     `json:"indy"`
}

func NewCredentialPreviewV2(attributes []CredentialPreviewAttributeV2) CredentialPreviewV2 {
	return CredentialPreviewV2{
		Type:       "issue-credential/2.0/credential-preview",
		Attributes: attributes,
	}
}

type CredentialPreviewV2 struct {
	Type       string                         `json:"@type"` // did:sov:BzCbsNYhMrjHiqZDTUASHg;spec/issue-credential/2.0/credential-preview
	Attributes []CredentialPreviewAttributeV2 `json:"attributes"`
}

type CredentialPreviewAttributeV2 struct {
	Name     string `json:"name"`
	MimeType string `json:"mime-type"` // optional
	Value    string `json:"value"`
}

type CredentialOfferV2 struct {
	Type              string              `json:"@type"` // did:sov:BzCbsNYhMrjHiqZDTUASHg;spec/issue-credential/2.0/offer-credential
	ID                string              `json:"@id"`
	Thread            Thread              `json:"~thread"`
	CredentialPreview CredentialPreviewV2 `json:"credential_preview"`
	OffersAttach      []AttachDecorator   `json:"offers~attach"`
	Formats           []Format            `json:"formats"`
}

type AttachDecorator struct {
	ID       string `json:"@id"`
	MimeType string `json:"mime-type"`
	Data     struct {
		Base64 string `json:"base64"`
	} `json:"data"`
}

type CredentialProposalV2 struct {
	ID                string              `json:"@id"`
	Type              string              `json:"@type"` // did:sov:BzCbsNYhMrjHiqZDTUASHg;spec/issue-credential/2.0/propose-credential
	FiltersAttach     []FilterAttach      `json:"filters~attach"`
	Comment           string              `json:"comment"`
	CredentialPreview CredentialPreviewV2 `json:"credential_proposal"`
	Formats           []Format            `json:"formats"`
}

type Format struct {
	AttachID string `json:"attach_id"` // dif or indy
	Format   string `json:"format"`    // dif/credential-manifest@v1.0 or hlindy-zkp-v1.0
}

type FilterAttach struct {
	ID       string `json:"@id"` // dif or indy
	MimeType string `json:"mime-type"`
	Data     struct {
		Base64 string `json:"base64"`
	} `json:"data"`
}

type CredentialRequestV2 struct {
	Type           string            `json:"@type"` // did:sov:BzCbsNYhMrjHiqZDTUASHg;spec/issue-credential/2.0/request-credential
	ID             string            `json:"@id"`
	Thread         Thread            `json:"~thread"`
	RequestsAttach []AttachDecorator `json:"requests~attach"`
	Formats        []Format          `json:"formats"`
}

type CredentialRequestMetadataV2 struct {
	Description string `json:"description"`
}

type CredentialIssue struct {
	Type              string            `json:"@type"` // did:sov:BzCbsNYhMrjHiqZDTUASHg;spec/issue-credential/2.0/issue-credential
	ID                string            `json:"@id"`
	CredentialsAttach []AttachDecorator `json:"credentials~attach"`
	Formats           []Format          `json:"formats"`
	Comment           string            `json:"comment"`
}

type CredentialExchangeDIF struct {
	CredentialExchangeDIFID string `json:"cred_ex_dif_id"`
	CreatedAt               string `json:"created_at"`
	CredentialExchangeID    string `json:"cred_ex_id"`
	Item                    string `json:"item"`
	State                   string `json:"state"`
	UpdatedAt               string `json:"updated_at"`
}

type CredentialExchangeIndy struct {
	CredentialExchangeIndyID  string `json:"cred_ex_indy_id"`
	CreatedAt                 string `json:"created_at"`
	UpdatedAt                 string `json:"updated_at"`
	CredentialExchangeID      string `json:"cred_ex_id"`
	RevocationRegistryID      string `json:"rev_reg_id"`
	CredentialRequestMetadata struct {
		MasterSecretBlindingData struct {
			VPrime  string      `json:"v_prime"`
			VrPrime interface{} `json:"vr_prime"`
		} `json:"master_secret_blinding_data"`
		Nonce            string `json:"nonce"`
		MasterSecretName string `json:"master_secret_name"`
	} `json:"cred_request_metadata"`
}
