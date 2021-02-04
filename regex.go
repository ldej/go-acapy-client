package acapy

import "regexp"

const (
	reDID                    = `^(did:sov:)?[123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz]{21,22}$`
	reCredentialDefinitionID = `^([123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz]{21,22}):3:CL:(([1-9][0-9]*)|([123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz]{21,22}:2:.+:[0-9.]+)):(.+)?$`
	reSchemaID               = `^[123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz]{21,22}:2:.+:[0-9.]+$`
)

var (
	compiledDID                  = regexp.MustCompile(reDID)
	compiledCredentialDefinition = regexp.MustCompile(reCredentialDefinitionID)
	compiledSchemaID             = regexp.MustCompile(reSchemaID)
)
