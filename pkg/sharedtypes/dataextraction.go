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

import "github.com/google/uuid"

// DataExtractionDocumentData represents the data extracted from a document.
type DataExtractionDocumentData struct {
	DocumentName      string    `json:"documentName"`
	DocumentId        string    `json:"documentId"`
	Guid              uuid.UUID `json:"guid"`
	Level             string    `json:"level"`
	ChildIds          []string  `json:"childIds"`
	ParentId          string    `json:"parentId"`
	PreviousSiblingId string    `json:"previousSiblingId"`
	NextSiblingId     string    `json:"nextSiblingId"`
	LastChildId       string    `json:"lastChildId"`
	FirstChildId      string    `json:"firstChildId"`
	Text              string    `json:"text"`
	Keywords          []string  `json:"keywords"`
	Summary           string    `json:"summary"`
	Embedding         []float32 `json:"embedding"`
}

type CodeGenerationElement struct {
	Guid uuid.UUID          `json:"guid"`
	Type CodeGenerationType `json:"type"`

	NamePseudocode string `json:"name_pseudocode"` // Function name without dependencies
	NameFormatted  string `json:"name_formatted"`  // Name of the function with spaces and without parameters
	Description    string `json:"description"`

	Name              string   `json:"name"`
	Dependencies      []string `json:"dependencies"`
	Summary           string   `json:"summary"`
	ReturnType        string   `json:"return"`
	ReturnElementList []string `json:"return_element_list"`
	ReturnDescription string   `json:"return_description"` // Return description
	Remarks           string   `json:"remarks"`

	// Only for type "method"
	PyaedtParents     []string `json:"inheritsfrom"`
	// Only for type "class"
	PyaedtGroup       string   `json:"typeof"`
	
	// Only for type "function" or "method"
	Parameters []XMLMemberParam `json:"parameters"`
	Example    XMLMemberExample `json:"example"`

	// Only for type "enum"
	EnumValues []string `json:"enum_values"`

	// Metadata for databases
	VectorDBMetadata any `json:"vector_db_metadata,omitempty"` // Optional metadata for vector databases
	GraphDBMetadata  any `json:"graph_db_metadata,omitempty"`  // Optional metadata for graph databases

	PyaedtGroup      string `json:"typeof"`
	PyaedtParents  []string  `json:"inheritsfrom"`
}

// Enum values for CodeGenerationType
type CodeGenerationType string

const (
	Function  CodeGenerationType = "Function"
	Method    CodeGenerationType = "Method"
	Class     CodeGenerationType = "Class"
	Parameter CodeGenerationType = "Parameter"
	Enum      CodeGenerationType = "Enum"
	Module    CodeGenerationType = "Module"
)

type XMLMemberExample struct {
	Description string               `xml:",chardata" json:"description"` // Text content of <example>
	Code        XMLMemberExampleCode `xml:"code,omitempty" json:"code"`   // Optional <code> element
}

type XMLMemberExampleCode struct {
	Type string `xml:"type,attr" json:"type"` // Attribute for <code>
	Text string `xml:",chardata" json:"text"` // Text content of <code>
}

type XMLMemberParam struct {
	Name        string `xml:"name" json:"name"`             // Attribute for <param>
	Type        string `xml:"type,omitempty" json:"type"`   // Attribute for <param>
	Description string `xml:",chardata" json:"description"` // Text content of <param>
}

type CodeGenerationExample struct {
	Guid                   uuid.UUID         `json:"guid"`
	Name                   string            `json:"name"`
	Dependencies           []string          `json:"dependencies"`
	DependencyEquivalences map[string]string `json:"dependency_equivalences"`
	Chunks                 []string          `json:"chunks"`
}

type CodeGenerationUserGuideSection struct {
	Name            string   `json:"name"`
	Title           string   `json:"title"`
	IsFirstChild    bool     `json:"is_first_child"`
	NextSibling     string   `json:"next_sibling"`
	NextParent      string   `json:"next_parent"`
	DocumentName    string   `json:"document_name"`
	Parent          string   `json:"parent"`
	Content         string   `json:"content"`
	Level           int      `json:"level"`
	Link            string   `json:"link"`
	ReferencedLinks []string `json:"referenced_links"`
	Chunks          []string `json:"chunks"`
}
