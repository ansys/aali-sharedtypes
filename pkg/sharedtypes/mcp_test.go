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
	"testing"
)

func TestGetAuthToken(t *testing.T) {
	tests := []struct {
		name      string
		authToken string
		envVar    string
		envValue  string
		expected  string
	}{
		{
			name:      "plain token",
			authToken: "my-secret-token",
			expected:  "my-secret-token",
		},
		{
			name:      "env var syntax",
			authToken: "${MCP_TEST_TOKEN}",
			envVar:    "MCP_TEST_TOKEN",
			envValue:  "token-from-env",
			expected:  "token-from-env",
		},
		{
			name:      "empty token",
			authToken: "",
			expected:  "",
		},
		{
			name:      "env var not set",
			authToken: "${UNSET_VAR}",
			expected:  "",
		},
		{
			name:      "partial syntax not resolved",
			authToken: "${INCOMPLETE",
			expected:  "${INCOMPLETE",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set env var if specified
			if tt.envVar != "" {
				os.Setenv(tt.envVar, tt.envValue)
				defer os.Unsetenv(tt.envVar)
			}

			config := &MCPConfig{AuthToken: tt.authToken}
			result := config.GetAuthToken()

			if result != tt.expected {
				t.Errorf("GetAuthToken() = %q, want %q", result, tt.expected)
			}
		})
	}
}
