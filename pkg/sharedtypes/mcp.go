// Copyright (C) 2025 - 2026 ANSYS, Inc. and/or its affiliates.
// SPDX-License-Identifier: MIT
//
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package sharedtypes

import (
	"os"
	"strings"
)

// TLSConfig represents the TLS/HTTPS configuration for secure connections
type TLSConfig struct {
	Enabled            bool   `json:"enabled"`            // Whether TLS is enabled
	CertFile           string `json:"certFile,omitempty"` // Path to the TLS certificate file
	KeyFile            string `json:"keyFile,omitempty"`  // Path to the TLS key file
	CAFile             string `json:"caFile,omitempty"`   // Path to the CA certificate file
	InsecureSkipVerify bool   `json:"insecureSkipVerify"` // Whether to skip server certificate verification(dev)
}

// MCPConfig represents the configuration for MCP connections
type MCPConfig struct {
	ServerURL string     `json:"serverURL"`     // URL of the MCP server endpoint
	Transport string     `json:"transport"`     // Connection protocol: "stdio", "http", "websocket"
	AuthToken string     `json:"authToken"`     // Authentication token, supports ${ENV_VAR} syntax
	Timeout   int        `json:"timeout"`       // Connection timeout in seconds
	TLS       *TLSConfig `json:"tls,omitempty"` // TLS configuration for secure connections
}

// MCPTool represents a tool definition in the Model Context Protocol.
type MCPTool struct {
	Name         string                 `json:"name"`                   // Unique identifier for the tool
	OriginalName string                 `json:"originalName,omitempty"` // Original name of the tool
	Description  string                 `json:"description,omitempty"`  // Human-readable description of what the tool does
	InputSchema  map[string]interface{} `json:"inputSchema"`            // JSON Schema for the tool's parameters
	ServerURL    string                 `json:"serverURL,omitempty"`    // URL of the MCP server that provides this tool
}

// MCPContentItem represents a single content item in an MCP response.
type MCPContentItem struct {
	Type     string `json:"type"`               // Content type: "text", "image", "resource"
	Text     string `json:"text,omitempty"`     // Text content
	Data     string `json:"data,omitempty"`     // Base64-encoded data for images
	MimeType string `json:"mimeType,omitempty"` // MIME type for binary content
	URI      string `json:"uri,omitempty"`      // Resource URI for resource content
}

// MCPResource represents a resource available from an MCP server.
type MCPResource struct {
	URI         string `json:"uri"`                   // Unique resource identifier
	Name        string `json:"name"`                  // Human-readable resource name
	Description string `json:"description,omitempty"` // Description of the resource
	MimeType    string `json:"mimeType,omitempty"`    // MIME type of the resource content
}

// MCPPrompt represents a prompt template available from an MCP server.
type MCPPrompt struct {
	Name        string              `json:"name"`                  // Unique prompt identifier
	Description string              `json:"description,omitempty"` // Description of what the prompt does
	Arguments   []MCPPromptArgument `json:"arguments,omitempty"`   // Parameters the prompt accepts
}

// MCPPromptArgument represents an argument for a prompt template.
type MCPPromptArgument struct {
	Name        string `json:"name"`                  // Argument name
	Description string `json:"description,omitempty"` // Description of the argument
	Required    bool   `json:"required,omitempty"`    // Whether the argument is required
}

// GetAuthToken returns the authentication token, resolving environment variables if needed
// ${MCP_TOKEN} will return the value of the MCP_TOKEN environment variable
func (config *MCPConfig) GetAuthToken() string {
	if len(config.AuthToken) > 3 &&
		strings.HasPrefix(config.AuthToken, "${") &&
		strings.HasSuffix(config.AuthToken, "}") {
		envVar := config.AuthToken[2 : len(config.AuthToken)-1]
		return os.Getenv(envVar)
	}
	return config.AuthToken
}
