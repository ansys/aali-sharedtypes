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
	ServerURL string `json:"serverURL"`
	Transport string `json:"transport"`
	AuthToken string `json:"authToken"`
	Timeout   int    `json:"timeout"`
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

// DiscoverServerResponse represents the response from MCP server discovery
type DiscoverServerResponse struct {
	ServerURL            string   `json:"serverURL"`
	Status               string   `json:"status"`
	RequiresAuth         bool     `json:"requiresAuth"`
	AvailableTransports  []string `json:"availableTransports"`
	HasTools             bool     `json:"hasTools"`
	ToolsCount           int      `json:"toolsCount"`
	HasResources         bool     `json:"hasResources"`
	ResourcesCount       int      `json:"resourcesCount"`
	HasPrompts           bool     `json:"hasPrompts"`
	PromptsCount         int      `json:"promptsCount"`
	RecommendedTimeout   int      `json:"recommendedTimeout"`
	RecommendedTransport string   `json:"recommendedTransport"`
	Error                string   `json:"error,omitempty"`
	Note                 string   `json:"note,omitempty"`
}