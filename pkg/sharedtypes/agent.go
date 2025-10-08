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
	"encoding/json"

	"github.com/ansys/aali-sharedtypes/pkg/logging"
)

// Message represents the JSON message you are expecting
type SessionContext struct {
	WorkflowId string            `json:"workflow_id,omitempty"` // Workflow ID, only relevant if "workflow_endpoint" is "custom"
	Variables  map[string]string `json:"variables,omitempty"`   // Variables to be passed to the workflow
	// Exec connection
	SessionType string `json:"session_type"`      // Type of session: "workflow", "exec"; default is "workflow" if not provided
	ExecId      string `json:"exec_id,omitempty"` // Unique identifier of connecting Exec, only relevant if "session_type" is "exec"
	// Authentication
	JwtToken string `json:"jwt_token"` // JWT token for authentication for "session_type" "workflow" (optional)
	ApiKey   string `json:"api_key"`   // API key for authentication, used if "session_type" is "exec" or can be used if "session_type" is "workflow" for authentication
	// Snapshot logic
	SnapshotId     string `json:"snapshot_id,omitempty"`     // Snapshot ID, only relevant if "session_type" is "workflow"; if defined, the given snapshot will retrived from the database
	WorkflowRunId  string `json:"workflow_run_id,omitempty"` // Workflow run ID, only relevant if "session_type" is "workflow"; if defined, mandatory if "snapshot_id" is defined in order to retrieve the snapshot from the database
	UserId         string `json:"user_id,omitempty"`         // User ID, only relevant if "session_type" is "workflow"; if defined, mandatory if "snapshot_id" is defined in order to retrieve the snapshot from the database
	StoreSnapshots bool   `json:"store_snapshots,omitempty"` // Store snapshots, only relevant if "session_type" is "workflow"; if true, all taken snapshots will be stored in the database
	// Other
	IpAddress string `json:"ip_address,omitempty"` // IP address of the client, required for some auth_types
}

// ConversationHistoryMessage is a structure that contains the message ID, role, content, and images of a conversation history message.
type ConversationHistoryMessage struct {
	MessageId        string   `json:"message_id"`
	Role             string   `json:"role"`
	Content          string   `json:"content"`
	Images           []string `json:"images"` // image in base64 format
	PositiveFeedback bool     `json:"positive_feedback"`
	NegativeFeedback bool     `json:"negative_feedback"`
	FeedbackText     string   `json:"feedback_text,omitempty"` // Optional feedback text
}

// Feedback is a structure that contains the conversation history, message ID, and feedback options of a workflow feedback.
type Feedback struct {
	ConversationHistory []ConversationHistoryMessage `json:"conversation"`
	MessageId           string                       `json:"message_id"`
	AddPositive         bool                         `json:"add_positive"`
	AddNegative         bool                         `json:"add_negative"`
	RemovePositive      bool                         `json:"remove_positive"`
	RemoveNegative      bool                         `json:"remove_negative"`
	FeedbackText        string                       `json:"feedback_text,omitempty"` // Optional feedback text
}

// SetSessionContext sets the SessionContext struct from the JSON payload
//
// Parameters:
//   - msg: the JSON payload
//
// Returns:
//   - SessionContext: the SessionContext struct
func ExtractSessionContext(ctx *logging.ContextMap, msg []byte) (SessionContext, error) {
	var SessionContext SessionContext

	// Unmarshal the JSON payload into your struct
	if err := json.Unmarshal(msg, &SessionContext); err != nil {
		logging.Log.Error(ctx, "Error decoding JSON:", err)
		return SessionContext, err
	} else {
		return SessionContext, nil
	}
}
