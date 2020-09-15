package acapy

type Schema struct {
	Version    string   `json:"schema_version"`
	Name       string   `json:"schema_name"`
	Attributes []string `json:"attributes"`
}

type SchemaResponse struct {
	SchemaID string `json:"schema_id"`
	Schema   struct {
		Ver       string   `json:"ver"`
		ID        string   `json:"id"`
		Name      string   `json:"name"`
		Version   string   `json:"version"`
		AttrNames []string `json:"attr_names"`
	} `json:"schema"`
	SeqNo int `json:"seqNo,omitempty"`
}

func (c *Client) RegisterSchema(name string, version string, attributes []string) (SchemaResponse, error) {
	var schema = Schema{
		Name:       name,
		Version:    version,
		Attributes: attributes,
	}
	var schemaResponse SchemaResponse
	err := c.post(c.AcapyURL+"/schemas", nil, schema, &schemaResponse)
	return schemaResponse, err
}
