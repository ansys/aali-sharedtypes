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

// DefaultFields represents the default fields for the user query.
type AnsysGPTDefaultFields struct {
	QueryWord         string
	FieldName         string
	FieldDefaultValue string
}

// ACSSearchResponse represents the response from the ACS search.
type ACSSearchResponse struct {
	Physics             string  `json:"physics"`
	SourceTitleLvl3     string  `json:"sourceTitle_lvl3"`
	SourceURLLvl3       string  `json:"sourceURL_lvl3"`
	TokenSize           int     `json:"tokenSize"`
	SourceTitleLvl2     string  `json:"sourceTitle_lvl2"`
	Weight              float64 `json:"weight"`
	SourceURLLvl2       string  `json:"sourceURL_lvl2"`
	Product             string  `json:"product"`
	Content             string  `json:"content"`
	TypeOFasset         string  `json:"typeOFasset"`
	Version             string  `json:"version"`
	SearchScore         float64 `json:"@search.score"`
	SearchRerankerScore float64 `json:"@search.rerankerScore"`
	IndexName           string  `json:"indexName"`
}

// AnsysGPTCitation represents the citation from the AnsysGPT.
type AnsysGPTCitation struct {
	Title     string  `json:"Title"`
	URL       string  `json:"URL"`
	Relevance float64 `json:"Relevance"`
}

// RetrieverModuleChunk represents a chunk of data context received from the retriever module.
type AnsysGPTRetrieverModuleChunk struct {
	Captions        string  `json:"@search.captions"`
	Highlights      string  `json:"@search.highlights"`
	RerankerScore   float64 `json:"@search.reranker_score"`
	Score           float64 `json:"@search.score"`
	Content         string  `json:"content"`
	IndexName       string  `json:"index_name"`
	Physics         string  `json:"physics"`
	Product         string  `json:"product"`
	SourceTitleLvl2 string  `json:"sourceTitle_lvl2"`
	SourceTitleLvl3 string  `json:"sourceTitle_lvl3"`
	SourceUrlLvl2   string  `json:"sourceURL_lvl2"`
	SourceUrlLvl3   string  `json:"sourceURL_lvl3"`
	TokenSize       int     `json:"token_size"`
	TypeOfAsset     string  `json:"typeOFasset"`
	Version         string  `json:"version"`
	Weight          float64 `json:"weight"`
}
