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

package typeconverters

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/ansys/aali-sharedtypes/pkg/sharedtypes"
)

func TestJSONToGo(t *testing.T) {
	tests := []struct {
		jsonType  string
		goType    string
		expectErr bool
	}{
		{"string", "string", false},
		{"string(binary)", "[]byte", false},
		{"number", "float64", false},
		{"integer", "int", false},
		{"boolean", "bool", false},
		{"array<string>", "[]string", false},
		{"array<number>", "[]float64", false},
		{"array<integer>", "[]int", false},
		{"array<boolean>", "[]bool", false},
		{"array<object>", "", true},
		{"dict[string][string]", "map[string]string", false},
		{"dict[string][integer]", "map[string]int", false},
		{"dict[string][number]", "map[string]float64", false},
		{"unsupportedType", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.jsonType, func(t *testing.T) {
			got, err := JSONToGo(tt.jsonType)
			if (err != nil) != tt.expectErr {
				t.Errorf("JSONToGo(%s) error = %v, expectErr %v", tt.jsonType, err, tt.expectErr)
				return
			}
			if got != tt.goType {
				t.Errorf("JSONToGo(%s) = %v, want %v", tt.jsonType, got, tt.goType)
			}
		})
	}
}

func TestGoToJSON(t *testing.T) {
	tests := []struct {
		goType   string
		jsonType string
	}{
		{"string", "string"},
		{"[]byte", "string(binary)"},
		{"float32", "number"},
		{"float64", "number"},
		{"int", "integer"},
		{"int8", "integer"},
		{"int16", "integer"},
		{"int32", "integer"},
		{"int64", "integer"},
		{"uint", "integer"},
		{"uint8", "integer"},
		{"uint16", "integer"},
		{"uint32", "integer"},
		{"uint64", "integer"},
		{"bool", "boolean"},
		{"map[string]string", "dict[string][string]"},
		{"map[string]int", "dict[string][integer]"},
		{"map[string]float64", "dict[string][number]"},
		{"map[string]interface{}", "dict[string][object]"},
		{"interface{}", "object"},
		{"[]string", "array<string>"},
		{"[]float64", "array<number>"},
		{"[]int", "array<integer>"},
		{"unknown", "object"},
	}

	for _, test := range tests {
		result := GoToJSON(test.goType)
		if result != test.jsonType {
			t.Errorf("GoToJSON(%q) = %q; want %q", test.goType, result, test.jsonType)
		}
	}
}

func TestConvertStringToGivenType(t *testing.T) {
	tests := []struct {
		value       string
		goType      string
		expected    interface{}
		expectedErr error
	}{
		{"42", "int", 42, nil},
		{"true", "bool", true, nil},
		{"3.14", "float64", 3.14, nil},
		{`["a","b","c"]`, "[]string", []string{"a", "b", "c"}, nil},
		{"", "[]int", []int{}, nil},
		{"{\"key\":\"value\"}", "map[string]string", map[string]string{"key": "value"}, nil},
		{"", "map[string]float64", map[string]float64{}, nil},
		{"{}", "map[string]bool", map[string]bool{}, nil},
		{"[]", "[]DbJsonFilter", []sharedtypes.DbJsonFilter{}, nil},
		{"", "*chan string", (*chan string)(nil), nil},
		{`{"serverURL":"http://localhost:8080","transport":"http","authToken":"secret123","timeout":30}`, "MCPConfig", sharedtypes.MCPConfig{ServerURL: "http://localhost:8080", Transport: "http", AuthToken: "secret123", Timeout: 30}, nil},
		{`[{"serverURL":"http://localhost:8080","transport":"http","authToken":"secret123","timeout":30}]`, "[]MCPConfig", []sharedtypes.MCPConfig{{ServerURL: "http://localhost:8080", Transport: "http", AuthToken: "secret123", Timeout: 30}}, nil},
		// Add more test cases as needed for each supported type
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%s_%s", test.goType, test.value), func(t *testing.T) {
			output, _, err := ConvertStringToGivenType(test.value, test.goType)
			if err != nil && test.expectedErr == nil {
				t.Errorf("Expected no error, got: %v", err)
			}
			if err == nil && test.expectedErr != nil {
				t.Errorf("Expected error: %v, got nil", test.expectedErr)
			}
			if fmt.Sprintf("%v", output) != fmt.Sprintf("%v", test.expected) {
				t.Errorf("Expected output: %v, got: %v", test.expected, output)
			}
		})
	}
}

func TestDeepCopy(t *testing.T) {
	type TestData struct {
		Name string
		Age  int
	}

	src := TestData{Name: "John", Age: 30}
	dst := new(TestData)

	err := DeepCopy(src, dst)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(src, *dst) {
		t.Errorf("deep copy failed, got: %v, want: %v", *dst, src)
	}
}
