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

// HandlerRequest represents the client request for a specific chat or embeddings operation.
type HandlerRequest struct {
	Adapter             string            `json:"adapter"` // "chat", "embeddings"
	InstructionGuid     string            `json:"instructionGuid"`
	ModelIds            []string          `json:"modelIds"`                   // optional model ids to define a set of specific models to be used for this request
	ModelCategory       []string          `json:"modelCategory"`              // optional model category; define one or more categories to filter models; models of the specified categories from first to last will be used for this request if available
	Data                interface{}       `json:"data"`                       // for embeddings, this can be a string or []string; for chat, only string is allowed
	Images              []string          `json:"images"`                     // List of images in base64 format
	MCPTools            []MCPTool         `json:"mcpTools,omitempty"`         // MCP tool definitions for tool calling support
	ChatRequestType     string            `json:"chatRequestType"`            // "summary", "code", "keywords", "general"; only relevant if "adapter" is "chat"
	DataStream          bool              `json:"dataStream"`                 // only relevant if "adapter" is "chat"
	MaxNumberOfKeywords uint32            `json:"maxNumberOfKeywords"`        // only relevant if "chatRequestType" is "keywords"
	IsConversation      bool              `json:"isConversation"`             // only relevant if "chatRequestType" is "code"
	ConversationHistory []HistoricMessage `json:"conversationHistory"`        // only relevant if "isConversation" is true
	GeneralContext      string            `json:"generalContext"`             // any added context you might need
	MsgContext          string            `json:"msgContext"`                 // any added context you might need
	SystemPrompt        interface{}       `json:"systemPrompt"`               // only relevant if "chatRequestType" is "general"
	ModelOptions        ModelOptions      `json:"modelOptions,omitempty"`     // only relevant if "adapter" is "chat"
	EmbeddingOptions    EmbeddingOptions  `json:"embeddingOptions,omitempty"` // only relevant if "adapter" is "embeddings"
}

// HandlerResponse represents the LLM Handler response for a specific request.
type HandlerResponse struct {
	// Common properties
	InstructionGuid string `json:"instructionGuid"`
	Type            string `json:"type"` // "info", "error", "chat", "embeddings"

	// Chat properties
	IsLast           *bool      `json:"isLast,omitempty"`
	Position         *uint32    `json:"position,omitempty"`
	InputTokenCount  *int       `json:"inputTokenCount,omitempty"`
	OutputTokenCount *int       `json:"outputTokenCount,omitempty"`
	ChatData         *string    `json:"chatData,omitempty"`
	ToolCalls        []ToolCall `json:"toolCalls,omitempty"` // Structured tool calls from LLM

	// Embeddings properties
	EmbeddedData   interface{} `json:"embeddedData,omitempty"`   // []float32 or [][]float32; for BAAI/bge-m3 these are dense vectors
	LexicalWeights interface{} `json:"lexicalWeights,omitempty"` // map[uint]float32 or []map[uint]float32; only for BAAI/bge-m3
	ColbertVecs    interface{} `json:"colbertVecs,omitempty"`    // [][]float32 or [][][]float32; only for BAAI/bge-m3

	// Error properties
	Error *ErrorResponse `json:"error,omitempty"`

	// Info properties
	InfoMessage *string `json:"infoMessage,omitempty"`
}

// HasToolCalls returns true if the response contains tool calls.
func (hr *HandlerResponse) HasToolCalls() bool {
	return len(hr.ToolCalls) > 0
}

// ErrorResponse represents the error response sent to the client when something fails during the processing of the request.
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// TransferDetails holds communication channels for the websocket listener and writer.
type TransferDetails struct {
	ResponseChannel chan HandlerResponse
	RequestChannel  chan HandlerRequest
}

// HistoricMessage represents a past chat message.
type HistoricMessage struct {
	Role       string     `json:"role"`
	Content    string     `json:"content"`
	Images     []string   `json:"images"`               // image in base64 format
	ToolCallId *string    `json:"toolCallId,omitempty"` // Tool call ID for tool responses
	ToolCalls  []ToolCall `json:"toolCalls,omitempty"`  // Tool calls made by assistant
}

// ModelOptions represents options for provider-specific API calls.
type ModelOptions struct {
	FrequencyPenalty *float32 `json:"frequencyPenalty,omitempty" yaml:"FREQUENCY_PENALTY,omitempty"`
	MaxTokens        *int32   `json:"maxTokens,omitempty" yaml:"MAX_TOKENS,omitempty"`
	PresencePenalty  *float32 `json:"presencePenalty,omitempty" yaml:"PRESENCE_PENALTY,omitempty"`
	Stop             []string `json:"stop,omitempty" yaml:"STOP,omitempty"`
	Temperature      *float32 `json:"temperature,omitempty" yaml:"TEMPERATURE,omitempty"`
	TopP             *float32 `json:"topP,omitempty" yaml:"TOP_P,omitempty"`

	// Reasoning effort level
	ReasoningEffort *string `json:"reasoningEffort,omitempty" yaml:"REASONING_EFFORT,omitempty"`
	// Reasoning summary format
	ReasoningSummary *string `json:"reasoningSummary,omitempty" yaml:"REASONING_SUMMARY,omitempty"`
	Verbosity        *string `json:"verbosity,omitempty" yaml:"VERBOSITY,omitempty"` // "low" | "medium" | "high"
}

// EmbeddingOptions represents the options for an embeddings request.
type EmbeddingOptions struct {
	ReturnDense   *bool `json:"returnDense"`   // Include dense vectors in response
	ReturnSparse  *bool `json:"returnSparse"`  // Include lexical weights in response
	ReturnColbert *bool `json:"returnColbert"` // Include colbert vectors in response
}

// EmbeddingResult holds both dense and sparse embeddings
type EmbeddingResult struct {
	Dense  []float32
	Sparse map[uint]float32
}

// ToolCall represents a tool invocation from the model.
type ToolCall struct {
	ID    string                 `json:"id"`
	Type  string                 `json:"type"`
	Name  string                 `json:"name"`
	Input map[string]interface{} `json:"input"`
}

// ToolResult represents the result of a tool execution.
type ToolResult struct {
	ToolCallID   string           `json:"tool_call_id"`            // Matches the ID from the original tool call
	Content      string           `json:"content"`                 // Primary text content
	ContentItems []MCPContentItem `json:"content_items,omitempty"` // Content items for multi-modal support
	IsError      bool             `json:"is_error"`                // True if tool execution failed
}
