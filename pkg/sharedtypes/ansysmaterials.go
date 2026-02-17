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

// Represents a criterion returned from the llm
type MaterialLlmCriterion struct {
	AttributeName string `json:"attributeName"`
	Explanation   string `json:"explanation"`
	Confidence    int    `json:"confidence"`
}

// Represents a criterion with its GUID
type MaterialCriterionWithGuid struct {
	AttributeName string `json:"attributeName"`
	AttributeGuid string `json:"attributeGuid"`
	Explanation   string `json:"explanation"`
	Confidence    int    `json:"confidence"`
}

// Represents a defined material attribute with its name and GUID.
type MaterialAttribute struct {
	Name string `json:"name"`
	Guid string `json:"guid"`
}
