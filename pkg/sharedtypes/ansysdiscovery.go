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

// SimulationInput captures the schema for simulation validation payloads.
type DiscoverySimulationInput struct {
	SimulationName     string               `json:"simulationName"`
	SimulationType     string               `json:"simulationType"`
	Model              string               `json:"model"`
	Objective          string               `json:"objective"`
	UserID             string               `json:"userId"`
	OtherInformation   string               `json:"otherInformation,omitempty"`
	Dimensions         Dimensions           `json:"dimensions"`
	Materials          []Material           `json:"materials"`
	BoundaryConditions []BoundaryCondition  `json:"boundaryConditions"`
	Attachments        []Attachment         `json:"attachments,omitempty"`
}

// Dimensions defines spatial extents and their units.
type DiscoveryDimensions struct {
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
	Z     float64 `json:"z"`
	Units string  `json:"units"`
}

// Material describes a labeled material state.
type DiscoveryMaterial struct {
	Label string `json:"label"`
	Name  string `json:"name"`
	State string `json:"state"`
}

// BoundaryCondition represents physics constraints for the simulation.
type DiscoveryBoundaryCondition struct {
	Index          int                      `json:"index"`
	ProxyGuid      string                   `json:"proxyGuid"`
	ProxyLabel     string                   `json:"proxyLabel"`
	Type           string                   `json:"type"`
	Classification int                      `json:"classification"`
	Rationale      string                   `json:"rationale"`
	Details        map[string]interface{}   `json:"details"`
	Guids          []string                 `json:"guids"`
	Names          []string                 `json:"names"`
	EntityIdsNames []map[string]interface{} `json:"entityIdsNames"`
}

// Attachment holds auxiliary binary payloads (e.g., base64-encoded uploads).
type DiscoveryAttachment struct {
	FileName string `json:"fileName"`
	Data     []byte `json:"data"`
}
