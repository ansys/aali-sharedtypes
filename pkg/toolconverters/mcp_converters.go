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

	"github.com/ansys/aali-sharedtypes/pkg/logging"
	"github.com/ansys/aali-sharedtypes/pkg/sharedtypes"
	"github.com/openai/openai-go"
	openaiv2 "github.com/openai/openai-go/v2"
	openaiv2shared "github.com/openai/openai-go/v2/shared"
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
//	[]openai.ChatCompletionToolParam: OpenAI-formatted tools.
func ConvertMCPToOpenAIFormat(
	ctx *logging.ContextMap,
	mcpTools []interface{},
) []openai.ChatCompletionToolParam {
	var openaiTools []openai.ChatCompletionToolParam

	for i, mcpTool := range mcpTools {
		// Convert interface{} to map for field access
		toolMap, ok := mcpTool.(map[string]interface{})
		if !ok {
			logging.Log.Warnf(ctx, "Skipping tool %d: not a valid object", i)
			continue
		}

		// Extract required fields
		name, nameOk := toolMap["name"].(string)
		if !nameOk || name == "" {
			logging.Log.Warnf(ctx, "Skipping tool %d: missing or invalid 'name' field", i)
			continue
		}

		// Extract description (optional)
		description, _ := toolMap["description"].(string)
		if description == "" {
			logging.Log.Warnf(ctx, "Tool '%s': missing description (recommended for better LLM understanding)", name)
		}

		// Extract inputSchema (optional)
		inputSchema, schemaOk := toolMap["inputSchema"].(map[string]interface{})
		if !schemaOk {
			logging.Log.Warnf(ctx, "Tool '%s': missing or invalid 'inputSchema' (LLM may not understand parameters)", name)
			// Create empty schema as fallback
			inputSchema = map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			}
		}

		// Convert to OpenAI format using openai.F() wrappers
		openaiTool := openai.ChatCompletionToolParam{
			Type: openai.F(openai.ChatCompletionToolTypeFunction),
			Function: openai.F(openai.FunctionDefinitionParam{
				Name:        openai.String(name),
				Description: openai.String(description),
				Parameters:  openai.F(openai.FunctionParameters(inputSchema)),
			}),
		}

		openaiTools = append(openaiTools, openaiTool)
		logging.Log.Debugf(ctx, "Converted MCP tool '%s' to OpenAI format", name)
	}

	if len(openaiTools) > 0 {
		logging.Log.Infof(ctx, "Converted %d MCP tools to OpenAI format", len(openaiTools))
	} else if len(mcpTools) > 0 {
		logging.Log.Warnf(ctx, "No valid tools converted from %d MCP tools provided", len(mcpTools))
	}

	return openaiTools
}

// ConvertMCPToOpenAIV2Format converts MCP tools to OpenAI v2 SDK format for Azure GPT.
//
// Parameters:
//
//	ctx: The logging context map.
//	mcpTools: Array of MCP tool definitions.
//
// Returns:
//
//	[]openaiv2.ChatCompletionToolUnionParam: OpenAI v2 formatted tools.
func ConvertMCPToOpenAIV2Format(
	ctx *logging.ContextMap,
	mcpTools []interface{},
) []openaiv2.ChatCompletionToolUnionParam {
	var openaiTools []openaiv2.ChatCompletionToolUnionParam

	for i, mcpTool := range mcpTools {
		// Convert interface{} to map for field access
		toolMap, ok := mcpTool.(map[string]interface{})
		if !ok {
			logging.Log.Warnf(ctx, "Skipping tool %d: not a valid object", i)
			continue
		}

		// Extract required fields
		name, nameOk := toolMap["name"].(string)
		if !nameOk || name == "" {
			logging.Log.Warnf(ctx, "Skipping tool %d: missing or invalid 'name' field", i)
			continue
		}

		// Extract description (optional)
		description, _ := toolMap["description"].(string)
		if description == "" {
			logging.Log.Warnf(ctx, "Tool '%s': missing description (recommended for better LLM understanding)", name)
		}

		// Extract inputSchema (optional)
		inputSchema, schemaOk := toolMap["inputSchema"].(map[string]interface{})
		if !schemaOk {
			logging.Log.Warnf(ctx, "Tool '%s': missing or invalid 'inputSchema' (LLM may not understand parameters)", name)
			// Create empty schema as fallback
			inputSchema = map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			}
		}

		// Convert to OpenAI v2 format
		functionDef := openaiv2shared.FunctionDefinitionParam{
			Name:        name,
			Description: openaiv2.String(description),
			Parameters:  openaiv2shared.FunctionParameters(inputSchema),
		}

		openaiTool := openaiv2.ChatCompletionFunctionTool(functionDef)
		openaiTools = append(openaiTools, openaiTool)
		logging.Log.Debugf(ctx, "Converted MCP tool '%s' to OpenAI v2 format", name)
	}

	if len(openaiTools) > 0 {
		logging.Log.Infof(ctx, "Converted %d MCP tools to OpenAI v2 format", len(openaiTools))
	} else if len(mcpTools) > 0 {
		logging.Log.Warnf(ctx, "No valid tools converted from %d MCP tools provided", len(mcpTools))
	}

	return openaiTools
}

// ConvertOpenAIToolCallsToSharedTypes converts OpenAI ChatCompletionMessageToolCall responses to shared ToolCall format.
//
// Parameters:
//
//	ctx: The logging context map.
//	openaiToolCalls: Array of OpenAI tool call responses.
//
// Returns:
//
//	[]sharedtypes.ToolCall: Shared format tool calls.
func ConvertOpenAIToolCallsToSharedTypes(
	ctx *logging.ContextMap,
	openaiToolCalls []openai.ChatCompletionMessageToolCall,
) []sharedtypes.ToolCall {
	var toolCalls []sharedtypes.ToolCall

	for _, tc := range openaiToolCalls {
		// Skip tool calls with empty arguments
		if tc.Function.Arguments == "" {
			logging.Log.Warnf(ctx, "Tool call %s has empty arguments, skipping", tc.ID)
			continue
		}

		// Parse arguments
		var args map[string]interface{}
		if err := json.Unmarshal([]byte(tc.Function.Arguments), &args); err != nil {
			logging.Log.Warnf(ctx, "Failed to parse tool call arguments for %s: %v, skipping tool call", tc.Function.Name, err)
			continue
		}

		// Only append valid tool calls
		toolCalls = append(toolCalls, sharedtypes.ToolCall{
			ID:    tc.ID,
			Type:  "function",
			Name:  tc.Function.Name,
			Input: args,
		})
	}

	if len(toolCalls) > 0 {
		logging.Log.Infof(ctx, "Converted %d OpenAI tool calls to shared format", len(toolCalls))
	}

	return toolCalls
}

// ConvertOpenAIV2ToolCallsToSharedTypes converts OpenAI v2 SDK tool calls to shared ToolCall format.
//
// Parameters:
//
//	ctx: The logging context map.
//	openaiToolCalls: Array of OpenAI v2 tool call responses.
//
// Returns:
//
//	[]sharedtypes.ToolCall: Shared format tool calls.
func ConvertOpenAIV2ToolCallsToSharedTypes(
	ctx *logging.ContextMap,
	openaiToolCalls []openaiv2.ChatCompletionMessageToolCallUnion,
) []sharedtypes.ToolCall {
	var toolCalls []sharedtypes.ToolCall

	for _, tc := range openaiToolCalls {
		// Skip tool calls with empty arguments
		if tc.Function.Arguments == "" {
			logging.Log.Warnf(ctx, "Tool call %s has empty arguments, skipping", tc.ID)
			continue
		}

		// Parse arguments
		var args map[string]interface{}
		if err := json.Unmarshal([]byte(tc.Function.Arguments), &args); err != nil {
			logging.Log.Warnf(ctx, "Failed to parse tool call arguments for %s: %v, skipping tool call", tc.Function.Name, err)
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
		logging.Log.Infof(ctx, "Converted %d OpenAI v2 tool calls to shared format", len(toolCalls))
	}

	return toolCalls
}
