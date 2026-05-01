// Package kustomize provides kustomize integration for NixAI.
// This file embeds the OpenAPI schema for NixAI custom resource types,
// enabling strategic merge patches to work correctly.
package kustomize

import _ "embed"

//go:embed openapi.json
var openapiSchema []byte

// GetOpenAPISchema returns the embedded OpenAPI schema for NixAI types.
// This schema defines merge keys and patch strategies for custom arrays.
func GetOpenAPISchema() []byte {
	return openapiSchema
}

// GetOpenAPISchemaString returns the embedded OpenAPI schema as a string.
func GetOpenAPISchemaString() string {
	return string(openapiSchema)
}