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
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/ansys/aali-sharedtypes/pkg/aali_graphdb"
	"github.com/ansys/aali-sharedtypes/pkg/sharedtypes"
)

// JSONToGo converts a JSON data type to a Go data type.
//
// Parameters:
//
//	jsonType: The JSON data type to convert.
//
// Returns:
//
//	string: The Go data type.
//	error: An error if the JSON data type is not supported.
func JSONToGo(jsonType string) (string, error) {
	// Handle array types
	if strings.HasPrefix(jsonType, "array<") && strings.HasSuffix(jsonType, ">") {
		elementType := jsonType[6 : len(jsonType)-1]
		arrayType, err := JSONToGo(elementType)
		if err != nil {
			return "", err
		}

		return "[]" + arrayType, nil
	}

	// Handle dictionary types
	if strings.HasPrefix(jsonType, "dict[") && strings.HasSuffix(jsonType, "]") {
		// Extract the inner types of the dictionary
		inner := jsonType[5 : len(jsonType)-1]
		parts := strings.Split(inner, "][")
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid dictionary type: %s", jsonType)
		}

		keyType := parts[0]
		valueType := parts[1]

		// Convert the value type using JSONToGo
		goValueType, err := JSONToGo(valueType)
		if err != nil {
			return "", err
		}

		// Go maps always have string keys
		if keyType != "string" {
			return "", fmt.Errorf("unsupported key type for Go maps: %s (only string keys are allowed)", keyType)
		}

		return fmt.Sprintf("map[string]%s", goValueType), nil
	}

	switch {
	case jsonType == "string":
		return "string", nil
	case jsonType == "string(binary)":
		return "[]byte", nil
	case jsonType == "number":
		return "float64", nil // Default to float64 for general numeric values
	case jsonType == "integer":
		return "int", nil
	case jsonType == "boolean":
		return "bool", nil
	default:
		return "", fmt.Errorf("not supported JSON type: %s", jsonType)
	}
}

// GoToJSON converts a Go data type to a JSON data type.
//
// Parameters:
//
//	goType: The Go data type to convert.
//
// Returns:
//
//	string: The JSON data type.
func GoToJSON(goType string) string {
	if strings.HasPrefix(goType, "[]") && goType != "[]byte" {
		elementType := goType[2:]
		return "array<" + GoToJSON(elementType) + ">"
	}

	// Handle maps (map[string]T)
	if strings.HasPrefix(goType, "map[string]") {
		// Extract the value type (after "map[string]")
		valueType := goType[len("map[string]"):]
		return "dict[string][" + GoToJSON(valueType) + "]"
	}

	switch goType {
	case "string":
		return "string"
	case "float32", "float64":
		return "number"
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
		return "integer"
	case "bool":
		return "boolean"
	case "[]byte":
		return "string(binary)"
	default:
		return "object"
	}
}

// ConvertStringToGivenType converts a string to a given Go type.
//
// Parameters:
// - value: a string containing the value to convert
// - goType: a string containing the Go type to convert to
//
// Returns:
// - output: an interface containing the converted value
// - err: an error containing the error message
func ConvertStringToGivenType(value string, goType string) (output interface{}, err error) {
	defer func() {
		r := recover()
		if r != nil {
			err = fmt.Errorf("panic occured in convertStringToGivenType: %v", r)
		}
	}()

	switch goType {
	case "string":
		return value, nil
	case "float32":
		if value == "" {
			value = "0"
		}
		return strconv.ParseFloat(value, 32)
	case "float64":
		if value == "" {
			value = "0"
		}
		return strconv.ParseFloat(value, 64)
	case "int":
		if value == "" {
			value = "0"
		}
		return strconv.Atoi(value)
	case "uint32":
		if value == "" {
			value = "0"
		}
		valueUint64, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			return nil, err
		}
		return uint32(valueUint64), nil
	case "uint64":
		if value == "" {
			value = "0"
		}
		return strconv.ParseUint(value, 10, 64)
	case "bool":
		if value == "" {
			value = "false"
		}
		return strconv.ParseBool(value)
	case "interface{}", "any":
		var output interface{}
		if value == "" || value == "null" {
			output = nil
		} else {
			trimmed := strings.TrimSpace(value)

			// Check if it looks like a JSON object or array
			if (strings.HasPrefix(trimmed, "{") && strings.HasSuffix(trimmed, "}")) ||
				(strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]")) ||
				(strings.HasPrefix(trimmed, "\"") && strings.HasSuffix(trimmed, "\"")) {
				// Try to unmarshal as JSON
				err := json.Unmarshal([]byte(value), &output)
				if err != nil {
					return nil, err
				}
			} else if trimmed == "true" || trimmed == "false" {
				// Handle boolean values
				output = trimmed == "true"
			} else if num, err := strconv.ParseFloat(trimmed, 64); err == nil {
				// Try to parse as number (int or float)
				intNum, err := strconv.ParseInt(trimmed, 10, 64)
				if err == nil {
					output = intNum
				} else {
					output = num
				}
			} else {
				// Default: treat as plain string
				output = value
			}
		}
		return output, nil
	case "[]interface{}":
		if value == "" {
			value = "[]"
		}
		output := []interface{}{}
		err := json.Unmarshal([]byte(value), &output)
		if err != nil {
			return nil, err
		}
		return output, nil
	case "[]string":
		if value == "" {
			value = "[]"
		}
		output := []string{}
		err := json.Unmarshal([]byte(value), &output)
		if err != nil {
			return nil, err
		}
		return output, nil
	case "[]float32":
		if value == "" {
			value = "[]"
		}
		output := []float32{}
		err := json.Unmarshal([]byte(value), &output)
		if err != nil {
			return nil, err
		}
		return output, nil
	case "[]float64":
		if value == "" {
			value = "[]"
		}
		output := []float64{}
		err := json.Unmarshal([]byte(value), &output)
		if err != nil {
			return nil, err
		}
		return output, nil
	case "[]int":
		if value == "" {
			value = "[]"
		}
		output := []int{}
		err := json.Unmarshal([]byte(value), &output)
		if err != nil {
			return nil, err
		}
		return output, nil
	case "[]bool":
		if value == "" {
			value = "[]"
		}
		output := []bool{}
		err := json.Unmarshal([]byte(value), &output)
		if err != nil {
			return nil, err
		}
		return output, nil
	case "[]byte":
		if value == "" {
			value = "[]"
		}
		output := []byte{}
		err := json.Unmarshal([]byte(value), &output)
		if err != nil {
			return nil, err
		}
		return output, nil
	case "[][]float32":
		if value == "" {
			value = "[]"
		}
		output := [][]float32{}
		err := json.Unmarshal([]byte(value), &output)
		if err != nil {
			return nil, err
		}
		return output, nil
	case "*chan string":
		var output *chan string
		output = nil
		return output, nil
	case "*chan interface{}":
		var output *chan interface{}
		output = nil
		return output, nil
	case "map[string]string":
		if value == "" {
			value = "{}"
		}
		output := map[string]string{}
		err := json.Unmarshal([]byte(value), &output)
		if err != nil {
			return nil, err
		}
		return output, nil
	case "map[string]float64":
		if value == "" {
			value = "{}"
		}
		output := map[string]float64{}
		err := json.Unmarshal([]byte(value), &output)
		if err != nil {
			return nil, err
		}
		return output, nil
	case "map[string]int":
		if value == "" {
			value = "{}"
		}
		output := map[string]int{}
		err := json.Unmarshal([]byte(value), &output)
		if err != nil {
			return nil, err
		}
		return output, nil
	case "map[string]bool":
		if value == "" {
			value = "{}"
		}
		output := map[string]bool{}
		err := json.Unmarshal([]byte(value), &output)
		if err != nil {
			return nil, err
		}
		return output, nil
	case "map[string][]string":
		if value == "" {
			value = "{}"
		}
		output := map[string][]string{}
		err := json.Unmarshal([]byte(value), &output)
		if err != nil {
			return nil, err
		}
		return output, nil
	case "map[string]map[string]string":
		if value == "" {
			value = "{}"
		}
		output := map[string]map[string]string{}
		err := json.Unmarshal([]byte(value), &output)
		if err != nil {
			return nil, err
		}
		return output, nil
	case "ParameterMap":
		if value == "" {
			value = "{}"
		}
		output := aali_graphdb.ParameterMap{}
		err := json.Unmarshal([]byte(value), &output)
		if err != nil {
			return nil, err
		}
		return output, nil
	case "[]map[string]string":
		if value == "" {
			value = "[]"
		}
		output := []map[string]string{}
		err := json.Unmarshal([]byte(value), &output)
		if err != nil {
			return nil, err
		}
		return output, nil
	case "map[uint]float32":
		if value == "" {
			value = "{}"
		}
		output := map[uint]float32{}
		err := json.Unmarshal([]byte(value), &output)
		if err != nil {
			return nil, err
		}
		return output, nil
	case "[]map[uint]float32":
		if value == "" {
			value = "[]"
		}
		output := []map[uint]float32{}
		err := json.Unmarshal([]byte(value), &output)
		if err != nil {
			return nil, err
		}
		return output, nil
	case "[]map[string]interface{}", "[]map[string]any":
		if value == "" {
			value = "[]"
		}
		output := []map[string]interface{}{}
		err := json.Unmarshal([]byte(value), &output)
		if err != nil {
			return nil, err
		}
		return output, nil
	case "DbArrayFilter":
		if value == "" {
			value = "{}"
		}
		output := sharedtypes.DbArrayFilter{}
		err := json.Unmarshal([]byte(value), &output)
		if err != nil {
			return nil, err
		}
		return output, nil
	case "DbFilters":
		if value == "" {
			value = "{}"
		}
		output := sharedtypes.DbFilters{}
		err := json.Unmarshal([]byte(value), &output)
		if err != nil {
			return nil, err
		}
		return output, nil
	case "Feedback":
		if value == "" {
			value = "{}"
		}
		output := sharedtypes.Feedback{}
		err := json.Unmarshal([]byte(value), &output)
		if err != nil {
			return nil, err
		}
		return output, nil
	case "ModelOptions":
		if value == "" {
			value = "{}"
		}
		output := sharedtypes.ModelOptions{}
		err := json.Unmarshal([]byte(value), &output)
		if err != nil {
			return nil, err
		}
		return output, nil
	case "[]DbJsonFilter":
		if value == "" {
			value = "[]"
		}
		output := []sharedtypes.DbJsonFilter{}
		err := json.Unmarshal([]byte(value), &output)
		if err != nil {
			return nil, err
		}
		return output, nil
	case "[]DbResponse":
		if value == "" {
			value = "[]"
		}
		output := []sharedtypes.DbResponse{}
		err := json.Unmarshal([]byte(value), &output)
		if err != nil {
			return nil, err
		}
		return output, nil
	case "[]ApiDbResponse":
		if value == "" {
			value = "[]"
		}
		output := []sharedtypes.ApiDbResponse{}
		err := json.Unmarshal([]byte(value), &output)
		if err != nil {
			return nil, err
		}
		return output, nil
        case "[]ExampleDbResponse":
		if value == "" {
			value = "[]"
		}
		output := []sharedtypes.ExampleDbResponse{}
		err := json.Unmarshal([]byte(value), &output)
		if err != nil {
			return nil, err
		}
		return output, nil
	case "[]HistoricMessage":
		if value == "" {
			value = "[]"
		}
		output := []sharedtypes.HistoricMessage{}
		err := json.Unmarshal([]byte(value), &output)
		if err != nil {
			return nil, err
		}
		return output, nil

	case "[]AnsysGPTDefaultFields":
		if value == "" {
			value = "[]"
		}
		output := []sharedtypes.AnsysGPTDefaultFields{}
		err := json.Unmarshal([]byte(value), &output)
		if err != nil {
			return nil, err
		}
		return output, nil

	case "[]ACSSearchResponse":
		if value == "" {
			value = "[]"
		}
		output := []sharedtypes.ACSSearchResponse{}
		err := json.Unmarshal([]byte(value), &output)
		if err != nil {
			return nil, err
		}
		return output, nil

	case "[]AnsysGPTCitation":
		if value == "" {
			value = "[]"
		}
		output := []sharedtypes.AnsysGPTCitation{}
		err := json.Unmarshal([]byte(value), &output)
		if err != nil {
			return nil, err
		}
		return output, nil

	case "[]AnsysGPTRetrieverModuleChunk":
		if value == "" {
			value = "[]"
		}
		output := []sharedtypes.AnsysGPTRetrieverModuleChunk{}
		err := json.Unmarshal([]byte(value), &output)
		if err != nil {
			return nil, err
		}
		return output, nil

	case "[]DbData":
		if value == "" {
			value = "[]"
		}
		output := []sharedtypes.DbData{}
		err := json.Unmarshal([]byte(value), &output)
		if err != nil {
			return nil, err
		}
		return output, nil
	case "[]CodeGenerationElement":
		if value == "" {
			value = "[]"
		}
		output := []sharedtypes.CodeGenerationElement{}
		err := json.Unmarshal([]byte(value), &output)
		if err != nil {
			return nil, err
		}
		return output, nil
	case "[]CodeGenerationExample":
		if value == "" {
			value = "[]"
		}
		output := []sharedtypes.CodeGenerationExample{}
		err := json.Unmarshal([]byte(value), &output)
		if err != nil {
			return nil, err
		}
		return output, nil
	case "[]CodeGenerationUserGuideSection":
		if value == "" {
			value = "[]"
		}
		output := []sharedtypes.CodeGenerationUserGuideSection{}
		err := json.Unmarshal([]byte(value), &output)
		if err != nil {
			return nil, err
		}
		return output, nil
	case "[]MaterialLlmCriterion":
		if value == "" {
			value = "[]"
		}
		output := []sharedtypes.MaterialLlmCriterion{}
		err := json.Unmarshal([]byte(value), &output)
		if err != nil {
			return nil, err
		}
		return output, nil
	case "[]MaterialCriterionWithGuid":
		if value == "" {
			value = "[]"
		}
		output := []sharedtypes.MaterialCriterionWithGuid{}
		err := json.Unmarshal([]byte(value), &output)
		if err != nil {
			return nil, err
		}
		return output, nil
	case "[]MaterialAttribute":
		if value == "" {
			value = "[]"
		}
		output := []sharedtypes.MaterialAttribute{}
		err := json.Unmarshal([]byte(value), &output)
		if err != nil {
			return nil, err
		}
		return output, nil
	}

	return nil, fmt.Errorf("unsupported GoType: '%s'", goType)
}

// ConvertGivenTypeToString converts a given Go type to a string.
//
// Parameters:
// - value: an interface containing the value to convert
// - goType: a string containing the Go type to convert from
//
// Returns:
// - string: a string containing the converted value
// - err: an error containing the error message
func ConvertGivenTypeToString(value interface{}, goType string) (output string, err error) {
	defer func() {
		r := recover()
		if r != nil {
			err = fmt.Errorf("panic occured in ConvertGivenTypeToString: %v", r)
		}
	}()

	switch goType {
	case "string":
		return value.(string), nil
	case "float32":
		return strconv.FormatFloat(float64(value.(float32)), 'f', -1, 32), nil
	case "float64":
		return strconv.FormatFloat(value.(float64), 'f', -1, 64), nil
	case "int":
		return strconv.Itoa(value.(int)), nil
	case "uint32":
		return strconv.FormatUint(uint64(value.(uint32)), 10), nil
	case "uint64":
		return strconv.FormatUint(value.(uint64), 10), nil
	case "bool":
		return strconv.FormatBool(value.(bool)), nil
	case "interface{}", "any":
		switch v := value.(type) {
		case string:
			return v, nil
		default:
			output, err := json.Marshal(value)
			if err != nil {
				return "", err
			}
			return string(output), nil
		}
	case "[]string":
		output, err := json.Marshal(value.([]string))
		if err != nil {
			return "", err
		}
		return string(output), nil
	case "[]interface{}":
		output, err := json.Marshal(value.([]interface{}))
		if err != nil {
			return "", err
		}
		return string(output), nil
	case "[]float32":
		output, err := json.Marshal(value.([]float32))
		if err != nil {
			return "", err
		}
		return string(output), nil
	case "[]float64":
		output, err := json.Marshal(value.([]float64))
		if err != nil {
			return "", err
		}
		return string(output), nil
	case "[]int":
		output, err := json.Marshal(value.([]int))
		if err != nil {
			return "", err
		}
		return string(output), nil
	case "[]bool":
		output, err := json.Marshal(value.([]bool))
		if err != nil {
			return "", err
		}
		return string(output), nil
	case "[]byte":
		output, err := json.Marshal(value.([]byte))
		if err != nil {
			return "", err
		}
		return string(output), nil
	case "[][]float32":
		output, err := json.Marshal(value.([][]float32))
		if err != nil {
			return "", err
		}
		return string(output), nil
	case "*chan string":
		return "", nil
	case "*chan interface{}":
		return "", nil
	case "map[string]string":
		output, err := json.Marshal(value.(map[string]string))
		if err != nil {
			return "", err
		}
		return string(output), nil
	case "map[string]float64":
		output, err := json.Marshal(value.(map[string]float64))
		if err != nil {
			return "", err
		}
		return string(output), nil
	case "map[string]int":
		output, err := json.Marshal(value.(map[string]int))
		if err != nil {
			return "", err
		}
		return string(output), nil
	case "map[string]bool":
		output, err := json.Marshal(value.(map[string]bool))
		if err != nil {
			return "", err
		}
		return string(output), nil
	case "map[string][]string":
		output, err := json.Marshal(value.(map[string][]string))
		if err != nil {
			return "", err
		}
		return string(output), nil
	case "map[string]map[string]string":
		output, err := json.Marshal(value.(map[string]map[string]string))
		if err != nil {
			return "", err
		}
		return string(output), nil
	case "ParameterMap":
		output, err := json.Marshal(value.(aali_graphdb.ParameterMap))
		if err != nil {
			return "", err
		}
		return string(output), nil
	case "[]map[string]string":
		output, err := json.Marshal(value.([]map[string]string))
		if err != nil {
			return "", err
		}
		return string(output), nil
	case "[]map[string]interface{}", "[]map[string]any":
		output, err := json.Marshal(value.([]map[string]interface{}))
		if err != nil {
			return "", err
		}
		return string(output), nil
	case "map[uint]float32":
		output, err := json.Marshal(value.(map[uint]float32))
		if err != nil {
			return "", err
		}
		return string(output), nil
	case "[]map[uint]float32":
		output, err := json.Marshal(value.([]map[uint]float32))
		if err != nil {
			return "", err
		}
		return string(output), nil
	case "DbArrayFilter":
		output, err := json.Marshal(value.(sharedtypes.DbArrayFilter))
		if err != nil {
			return "", err
		}
		return string(output), nil
	case "DbFilters":
		output, err := json.Marshal(value.(sharedtypes.DbFilters))
		if err != nil {
			return "", err
		}
		return string(output), nil
	case "Feedback":
		output, err := json.Marshal(value.(sharedtypes.Feedback))
		if err != nil {
			return "", err
		}
		return string(output), nil
	case "ModelOptions":
		output, err := json.Marshal(value.(sharedtypes.ModelOptions))
		if err != nil {
			return "", err
		}
		return string(output), nil
	case "[]DbJsonFilter":
		output, err := json.Marshal(value.([]sharedtypes.DbJsonFilter))
		if err != nil {
			return "", err
		}
		return string(output), nil
	case "[]DbResponse":
		output, err := json.Marshal(value.([]sharedtypes.DbResponse))
		if err != nil {
			return "", err
		}
		return string(output), nil
	case "[]ApiDbResponse":
		output, err := json.Marshal(value.([]sharedtypes.ApiDbResponse))
		if err != nil {
			return "", err
		}
		return string(output), nil
	case "[]ExampleDbResponse":
		output, err := json.Marshal(value.([]sharedtypes.ExampleDbResponse))
		if err != nil {
			return "", err
		}
		return string(output), nil
        case "[]HistoricMessage":
		output, err := json.Marshal(value.([]sharedtypes.HistoricMessage))
		if err != nil {
			return "", err
		}
		return string(output), nil
	case "[]AnsysGPTDefaultFields":
		output, err := json.Marshal(value.([]sharedtypes.AnsysGPTDefaultFields))
		if err != nil {
			return "", err
		}
		return string(output), nil
	case "[]ACSSearchResponse":
		output, err := json.Marshal(value.([]sharedtypes.ACSSearchResponse))
		if err != nil {
			return "", err
		}
		return string(output), nil
	case "[]AnsysGPTCitation":
		output, err := json.Marshal(value.([]sharedtypes.AnsysGPTCitation))
		if err != nil {
			return "", err
		}
		return string(output), nil
	case "[]AnsysGPTRetrieverModuleChunk":
		output, err := json.Marshal(value.([]sharedtypes.AnsysGPTRetrieverModuleChunk))
		if err != nil {
			return "", err
		}
		return string(output), nil
	case "[]DbData":
		output, err := json.Marshal(value.([]sharedtypes.DbData))
		if err != nil {
			return "", err
		}
		return string(output), nil
	case "[]CodeGenerationElement":
		output, err := json.Marshal(value.([]sharedtypes.CodeGenerationElement))
		if err != nil {
			return "", err
		}
		return string(output), nil
	case "[]CodeGenerationExample":
		output, err := json.Marshal(value.([]sharedtypes.CodeGenerationExample))
		if err != nil {
			return "", err
		}
		return string(output), nil
	case "[]CodeGenerationUserGuideSection":
		output, err := json.Marshal(value.([]sharedtypes.CodeGenerationUserGuideSection))
		if err != nil {
			return "", err
		}
		return string(output), nil
	case "[]MaterialLlmCriterion":
		output, err := json.Marshal(value.([]sharedtypes.MaterialLlmCriterion))
		if err != nil {
			return "", err
		}
		return string(output), nil
	case "[]MaterialCriterionWithGuid":
		output, err := json.Marshal(value.([]sharedtypes.MaterialCriterionWithGuid))
		if err != nil {
			return "", err
		}
		return string(output), nil
	case "[]MaterialAttribute":
		output, err := json.Marshal(value.([]sharedtypes.MaterialAttribute))
		if err != nil {
			return "", err
		}
		return string(output), nil
	}

	return "", fmt.Errorf("unsupported GoType: '%s'", goType)
}

// DeepCopy deep copies the source interface to the destination interface.
//
// Parameters:
// - src: an interface containing the source
// - dst: an interface containing the destination
//
// Returns:
// - err: an error containing the error message
func DeepCopy(src, dst interface{}) (err error) {
	defer func() {
		r := recover()
		if r != nil {
			err = fmt.Errorf("panic occured in DeepCopy: %v", r)
		}
	}()

	bytes, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, dst)
}
