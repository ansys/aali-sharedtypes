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
	"github.com/ansys/aali-sharedtypes/pkg/logging"
	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/packages/param"
	"github.com/openai/openai-go"
)

// ConvertMCPToOpenAIFormat converts MCP tools to OpenAI function calling format.
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

// ConvertMCPToMistralFormat converts MCP tools to Mistral function calling format.
func ConvertMCPToMistralFormat(
	ctx *logging.ContextMap,
	mcpTools []interface{},
) []map[string]interface{} {
	var mistralTools []map[string]interface{}

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

		// Convert to Mistral format (OpenAI-compatible)
		mistralTool := map[string]interface{}{
			"type": "function",
			"function": map[string]interface{}{
				"name":        name,
				"description": description,
				"parameters":  inputSchema,
			},
		}

		mistralTools = append(mistralTools, mistralTool)
		logging.Log.Debugf(ctx, "Converted MCP tool '%s' to Mistral format", name)
	}

	if len(mistralTools) > 0 {
		logging.Log.Infof(ctx, "Converted %d MCP tools to Mistral format", len(mistralTools))
	} else if len(mcpTools) > 0 {
		logging.Log.Warnf(ctx, "No valid tools converted from %d MCP tools provided", len(mcpTools))
	}

	return mistralTools
}

// ConvertMCPToAnthropicFormat converts MCP tools to Anthropic function calling format.
func ConvertMCPToAnthropicFormat(
	ctx *logging.ContextMap,
	mcpTools []interface{},
) []anthropic.ToolParam {
	var anthropicTools []anthropic.ToolParam

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

		// Extract required fields from inputSchema
		var required []string
		if req, ok := inputSchema["required"].([]interface{}); ok {
			for _, r := range req {
				if rStr, ok := r.(string); ok {
					required = append(required, rStr)
				}
			}
		}

		// Extract properties from inputSchema
		properties, _ := inputSchema["properties"]

		// Convert to Anthropic format
		anthropicTool := anthropic.ToolParam{
			Name: name,
			InputSchema: anthropic.ToolInputSchemaParam{
				Type:       "object",
				Properties: properties,
				Required:   required,
			},
		}

		// Add description if present
		if description != "" {
			anthropicTool.Description = param.NewOpt(description)
		}

		anthropicTools = append(anthropicTools, anthropicTool)
		logging.Log.Debugf(ctx, "Converted MCP tool '%s' to Anthropic format", name)
	}

	if len(anthropicTools) > 0 {
		logging.Log.Infof(ctx, "Converted %d MCP tools to Anthropic format", len(anthropicTools))
	} else if len(mcpTools) > 0 {
		logging.Log.Warnf(ctx, "No valid tools converted from %d MCP tools provided", len(mcpTools))
	}

	return anthropicTools
}
