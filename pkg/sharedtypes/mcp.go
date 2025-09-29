// Copyright (C) 2025 ANSYS, Inc. and/or its affiliates.
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

// MCPConfig represents the configuration for MCP connections
type MCPConfig struct {
	ServerURL string `json:"serverURL"` // URL of the MCP server endpoint
	Transport string `json:"transport"` // Connection protocol: "stdio", "http", "websocket"
	AuthToken string `json:"authToken"` // Authentication token, supports ${ENV_VAR} syntax
	Timeout   int    `json:"timeout"`   // Connection timeout in seconds
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

