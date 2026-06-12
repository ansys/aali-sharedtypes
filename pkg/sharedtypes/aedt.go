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
	"github.com/google/uuid"
)

// AedtEmbeddingOptions represents the options for an embeddings request.
type AedtEmbeddingOptions struct {
	ReturnDense   *bool `json:"returnDense"`   // Include dense vectors in response
	ReturnSparse  *bool `json:"returnSparse"`  // Include lexical weights in response
	ReturnColbert *bool `json:"returnColbert"` // Include colbert vectors in response
	IsPrompt      *bool `json:"isPrompt"`      // Is the query passage a prompt or code. For Nomic Embed Code embedding model
}

// Checked
// AedtCodeGenerationType represents the type of code generation element.
type AedtCodeGenerationElement struct {
	Guid              uuid.UUID          `json:"guid"`
	Type              CodeGenerationType `json:"type"`
	NamePseudocode    string             `json:"name_pseudocode"` // Function name without dependencies
	NameFormatted     string             `json:"name_formatted"`  // Name of the function with spaces and without parameters
	Description       string             `json:"description"`
	Name              string             `json:"name"`
	Dependencies      []string           `json:"dependencies"`
	Summary           string             `json:"summary"`
	ReturnType        string             `json:"return"`
	ReturnElementList []string           `json:"return_element_list"`
	ReturnDescription string             `json:"return_description"` // Return description
	Remarks           string             `json:"remarks"`

	// Only for type "method"
	PyaedtParents []string `json:"inheritsfrom"`
	// Only for type "class"
	PyaedtGroup string `json:"typeof"`

	// Only for type "function" or "method"
	Parameters []XMLMemberParam `json:"parameters"`
	Example    XMLMemberExample `json:"example"`
	// Only for type "enum"
	EnumValues []string `json:"enum_values"`
	// Metadata for databases
	VectorDBMetadata any `json:"vector_db_metadata,omitempty"` // Optional metadata for vector databases
	GraphDBMetadata  any `json:"graph_db_metadata,omitempty"`  // Optional metadata for graph databases
}

// Checked
// AedtApiDbResponse represents the response from the database.
// ApiDbResponse is now AedtApiDbResponse
// for remaining DbResponse, use the standard DbResponse
type AedtApiDbResponse struct {
	Guid           uuid.UUID `json:"guid"`
	Name           string    `json:"name"`
	Type           string    `json:"type"`
	ParentClass    string    `json:"parent_class"`
	PyaedtGroup    string    `json:"typeof,omitempty"`
	Summary        string    `json:"summary, omitempty"`
	NameFormatted  string    `json:"name_formatted"`
	NamePseudocode string    `json:"name_pseudocode"`
	//ParentId          *uuid.UUID             `json:"parent_id"`
	//Siblings  []DbData `json:"siblings,omitempty"`
}

// Checked
// AedtElementContextsTuple represents a tuple of element contexts.
type AedtElementContextsTuple struct {
	Params      string `json:"params"`
	Return      string `json:"return"`
	Example     string `json:"example"`
	Instruction string `json:"instruction"`
}

// Checked
// AedtCodeGenerationExample represents an example for code generation.
type AedtCodeGenerationExample struct {
	Guid    uuid.UUID `json:"guid"`
	Name    string    `json:"name"`
	Designs []string  `json:"designs"`
	Chunks  []string  `json:"chunks"`
}
