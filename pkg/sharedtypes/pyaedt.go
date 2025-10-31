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

// DesignContext represents the design context structure.
type DesignContext struct {
	AedtVersion        string                   `json:"aedtVersion"`
	PyaedtVersion      string                   `json:"pyaedtVersion"`
	Application        string                   `json:"application"`
	Design             string                   `json:"design"`
	ProjectName        string                   `json:"projectName"`
	Selections         []interface{}            `json:"selections"`
	Units              Units                    `json:"units"`
	CoordinateSystem   string                   `json:"coordinateSystem"`
	ObjectList         []Object                 `json:"objectList"`
	Planes             []string                 `json:"planes"`
	Materials          []string                 `json:"materials"`
	BoundaryConditions map[string]interface{}   `json:"boundaryConditions"`
	Excitations        []string                 `json:"excitations"`
	SolutionType       string                   `json:"solutionType"`
	GeneratedMesh      string                   `json:"generatedMesh"`
	AvailableSetups    map[string]Setup         `json:"availableSetups"`
	OutputVariables    []string                 `json:"outputVariables"`
	Sweeps             map[string][]interface{} `json:"sweeps"`
}

// Units represents the units structure in the design context for generic mode.
type Units struct {
	Angle        string `json:"Angle"`
	AngularSpeed string `json:"Angular Speed"`
	Capacitance  string `json:"Capacitance"`
	Conductance  string `json:"Conductance"`
	Current      string `json:"Current"`
	Frequency    string `json:"Frequency"`
	Inductance   string `json:"Inductance"`
	Length       string `json:"Length"`
	Mass         string `json:"Mass"`
	Power        string `json:"Power"`
	Resistance   string `json:"Resistance"`
	Speed        string `json:"Speed"`
	Temperature  string `json:"Temperature"`
	Time         string `json:"Time"`
	Voltage      string `json:"Voltage"`
}

// Setup represents the setup configuration
type Setup struct {
	ID                     int     `json:"ID"`
	SetupType              string  `json:"SetupType"`
	SolveType              string  `json:"SolveType"`
	Frequency              string  `json:"Frequency"`
	MaxDeltaE              float64 `json:"MaxDeltaE"`
	MaximumPasses          int     `json:"MaximumPasses"`
	MinimumPasses          int     `json:"MinimumPasses"`
	MinimumConvergedPasses int     `json:"MinimumConvergedPasses"`
	PercentRefinement      int     `json:"PercentRefinement"`
	IsEnabled              bool    `json:"IsEnabled"`
	MeshLink               struct {
		ImportMesh bool `json:"ImportMesh"`
	} `json:"MeshLink"`
	BasisOrder                     int     `json:"BasisOrder"`
	DoLambdaRefine                 bool    `json:"DoLambdaRefine"`
	DoMaterialLambda               bool    `json:"DoMaterialLambda"`
	SetLambdaTarget                bool    `json:"SetLambdaTarget"`
	Target                         float64 `json:"Target"`
	UseMaxTetIncrease              bool    `json:"UseMaxTetIncrease"`
	DrivenSolverType               string  `json:"DrivenSolverType"`
	EnhancedLowFreqAccuracy        bool    `json:"EnhancedLowFreqAccuracy"`
	EnhancedFEBIPreconditioner     bool    `json:"EnhancedFEBIPreconditioner"`
	SaveRadFieldsOnly              bool    `json:"SaveRadFieldsOnly"`
	SaveAnyFields                  bool    `json:"SaveAnyFields"`
	IESolverType                   string  `json:"IESolverType"`
	LambdaTargetForIESolver        float64 `json:"LambdaTargetForIESolver"`
	UseDefaultLambdaTgtForIESolver bool    `json:"UseDefaultLambdaTgtForIESolver"`
	IESolverAccuracy               string  `json:"IE Solver Accuracy"`
	InfiniteSphereSetup            int     `json:"InfiniteSphereSetup"`
	MaxPass                        int     `json:"MaxPass"`
	MinPass                        int     `json:"MinPass"`
	MinConvPass                    int     `json:"MinConvPass"`
	PerError                       int     `json:"PerError"`
	PerRefine                      int     `json:"PerRefine"`
	Sweeps                         struct {
		NextUniqueID  int  `json:"NextUniqueID"`
		MoveBackwards bool `json:"MoveBackwards"`
	} `json:"Sweeps"`
}

// Object represents a single object in the design context object list.
type Object struct {
	ID           int     `json:"id"`
	MaterialName string  `json:"material_name"`
	Name         string  `json:"name"`
	SolveInside  bool    `json:"solve_inside"`
	Transparency float64 `json:"transparency"`
}
