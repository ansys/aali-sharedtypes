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

// Package toolconverters provides utilities for converting MCP tool definitions to various LLM provider formats.
package toolconverters

import (
	"encoding/json"
	"fmt"

	"github.com/ansys/aali-sharedtypes/pkg/logging"
	"github.com/ansys/aali-sharedtypes/pkg/sharedtypes"
	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/shared"
)

// ConvertMCPToOpenAIFormat converts MCP tools to OpenAI function calling format.
//
// Parameters:
//
//	ctx: The logging context map.
//	mcpTools: Array of MCP tool definitions.
//
// Returns:
//
//	[]openai.ChatCompletionToolUnionParam: OpenAI formatted tools.
//	[]error: List of errors for tools that were skipped during conversion.
func ConvertMCPToOpenAIFormat(
	ctx *logging.ContextMap,
	mcpTools []interface{},
) ([]openai.ChatCompletionToolUnionParam, []error) {
	var openaiTools []openai.ChatCompletionToolUnionParam
	var errors []error

	for i, mcpTool := range mcpTools {
		// Convert interface{} to map for field access
		toolMap, ok := mcpTool.(map[string]interface{})
		if !ok {
			toolJSON, _ := json.Marshal(mcpTool)
			err := fmt.Errorf("tool at index %d is not a valid object, got type %T, value: %s", i, mcpTool, string(toolJSON))
			errors = append(errors, err)
			logging.Log.Errorf(ctx, "Skipping tool %d: not a valid object (type: %T, value: %s)", i, mcpTool, string(toolJSON))
			continue
		}

		// Extract required fields
		name, nameOk := toolMap["name"].(string)
		if !nameOk || name == "" {
			toolJSON, _ := json.Marshal(toolMap)
			err := fmt.Errorf("tool at index %d is missing or has invalid 'name' field, tool data: %s", i, string(toolJSON))
			errors = append(errors, err)
			logging.Log.Errorf(ctx, "Skipping tool %d: missing or invalid 'name' field, tool data: %s", i, string(toolJSON))
			continue
		}

		// Extract description (
		description, _ := toolMap["description"].(string)
		if description == "" {
			logging.Log.Warnf(ctx, "Tool '%s': missing description (recommended for better LLM understanding)", name)
		}

		// Extract inputSchema
		inputSchema, schemaOk := toolMap["inputSchema"].(map[string]interface{})
		if !schemaOk {
			logging.Log.Warnf(ctx, "Tool '%s': missing or invalid 'inputSchema' (LLM may not understand parameters)", name)
			// Create empty schema as fallback
			inputSchema = map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			}
		}

		// Convert to OpenAI format
		functionDef := shared.FunctionDefinitionParam{
			Name:        name,
			Description: openai.String(description),
			Parameters:  shared.FunctionParameters(inputSchema),
		}

		openaiTool := openai.ChatCompletionFunctionTool(functionDef)
		openaiTools = append(openaiTools, openaiTool)
		logging.Log.Debugf(ctx, "Converted MCP tool '%s' to OpenAI format", name)
	}

	if len(openaiTools) > 0 {
		logging.Log.Infof(ctx, "Converted %d MCP tools to OpenAI format", len(openaiTools))
	}
	if len(errors) > 0 {
		logging.Log.Errorf(ctx, "Failed to convert %d out of %d MCP tools (see detailed errors above)", len(errors), len(mcpTools))
	}

	return openaiTools, errors
}

// ConvertOpenAIToolCallsToSharedTypes converts OpenAI SDK tool calls to shared ToolCall format.
//
// Parameters:
//
//	ctx: The logging context map.
//	openaiToolCalls: Array of OpenAI tool call responses.
//
// Returns:
//
//	[]sharedtypes.ToolCall: Shared format tool calls.
//	[]error: List of errors for tool calls that were skipped during conversion.
func ConvertOpenAIToolCallsToSharedTypes(
	ctx *logging.ContextMap,
	openaiToolCalls []openai.ChatCompletionMessageToolCallUnion,
) ([]sharedtypes.ToolCall, []error) {
	var toolCalls []sharedtypes.ToolCall
	var errors []error

	for i, tc := range openaiToolCalls {
		// Skip tool calls with empty arguments
		if tc.Function.Arguments == "" {
			err := fmt.Errorf("tool call at index %d (ID: %s, Name: %s) has empty arguments", i, tc.ID, tc.Function.Name)
			errors = append(errors, err)
			logging.Log.Errorf(ctx, "Tool call at index %d (ID: %s, Name: %s) has empty arguments, skipping", i, tc.ID, tc.Function.Name)
			continue
		}

		// Parse arguments
		var args map[string]interface{}
		if err := json.Unmarshal([]byte(tc.Function.Arguments), &args); err != nil {
			parseErr := fmt.Errorf("failed to parse tool call at index %d (ID: %s, Name: %s): %w, raw arguments: %s",
				i, tc.ID, tc.Function.Name, err, tc.Function.Arguments)
			errors = append(errors, parseErr)
			logging.Log.Errorf(ctx, "Failed to parse tool call at index %d (ID: %s, Name: %s): %v, raw arguments: %s, skipping tool call",
				i, tc.ID, tc.Function.Name, err, tc.Function.Arguments)
			continue
		}

		// Only append valid tool calls
		toolCalls = append(toolCalls, sharedtypes.ToolCall{
			ID:    tc.ID,
			Type:  string(tc.Type),
			Name:  tc.Function.Name,
			Input: args,
		})
	}

	if len(toolCalls) > 0 {
		logging.Log.Infof(ctx, "Converted %d OpenAI tool calls to shared format", len(toolCalls))
	}
	if len(errors) > 0 {
		logging.Log.Errorf(ctx, "Failed to convert %d out of %d tool calls (see detailed errors above)", len(errors), len(openaiToolCalls))
	}

	return toolCalls, errors
}
