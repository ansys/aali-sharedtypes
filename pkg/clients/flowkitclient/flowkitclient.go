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

package flowkitclient

import (
	"context"
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/ansys/aali-sharedtypes/pkg/aaliflowkitgrpc"
	"github.com/ansys/aali-sharedtypes/pkg/logging"
	"github.com/ansys/aali-sharedtypes/pkg/sharedtypes"
	"github.com/ansys/aali-sharedtypes/pkg/typeconverters"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// HealthCheck checks the health of the external function server
// This function is used to check if the external function server is running and reachable
//
// Parameters:
//   - url: the URL of the external function server
//   - apiKey: the API key to authenticate with the external function server
//
// Returns:
//   - err: an error message if the gRPC call fails
func HealthCheck(url string, apiKey string) (err error) {
	// Set up a connection to the server.
	c, conn, err := createClient(url, apiKey)
	if err != nil {
		return fmt.Errorf("unable to connect to external function gRPC: %v", err)
	}
	defer conn.Close()

	// Create a context with a cancel
	ctxWithCancel, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Call HealthCheck
	_, err = c.HealthCheck(ctxWithCancel, &aaliflowkitgrpc.HealthRequest{})
	if err != nil {
		return fmt.Errorf("error in external function gRPC HealthCheck: %v", err)
	}

	return nil
}

// GetVersion retrieves the version of the external function server
// This function is used to get the version of the external function server
//
// Parameters:
//   - url: the URL of the external function server
//   - apiKey: the API key to authenticate with the external function server
//
// Returns:
//   - version: the version of the external function server
//   - err: an error message if the gRPC call fails
func GetVersion(url string, apiKey string) (version string, err error) {
	// Set up a connection to the server.
	c, conn, err := createClient(url, apiKey)
	if err != nil {
		return "", fmt.Errorf("unable to connect to external function gRPC: %v", err)
	}
	defer conn.Close()

	// Create a context with a cancel
	ctxWithCancel, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Call GetVersion
	resp, err := c.GetVersion(ctxWithCancel, &aaliflowkitgrpc.VersionRequest{})
	if err != nil {
		return "", fmt.Errorf("error in external function gRPC GetVersion: %v", err)
	}

	return resp.Version, nil
}

// Global variable to store the available functions, types and categories
var AvailableFunctions map[string]*sharedtypes.FunctionDefinition
var AvailableTypes map[string]bool
var AvailableCategories map[string]bool

// ListFunctionsAndSaveToInteralStates calls the ListFunctions gRPC and saves the functions to internal states
// This function is used to get the list of available functions from the external function server
// and save them to internal states
//
// Parameters:
//   - url: the URL of the external function server
//   - apiKey: the API key to authenticate with the external function server
//
// Returns:
//   - error: an error message if the gRPC call fails
func ListFunctionsAndSaveToInteralStates(url string, apiKey string) (err error) {
	defer func() {
		r := recover()
		if r != nil {
			err = fmt.Errorf("panic occurred in ListFunctionsAndSaveToInteralStates: %v", r)
		}
	}()

	// Set up a connection to the server.
	c, conn, err := createClient(url, apiKey)
	if err != nil {
		return fmt.Errorf("unable to connect to external function gRPC: %v", err)
	}
	defer conn.Close()

	// Create a context with a cancel
	ctxWithCancel, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Call ListFunctions
	listResp, err := c.ListFunctions(ctxWithCancel, &aaliflowkitgrpc.ListFunctionsRequest{})
	if err != nil {
		return fmt.Errorf("error in external function gRPC ListFunctions: %v", err)
	}

	// Save the functions to internal states
	for _, function := range listResp.Functions {
		// convert inputs and outputs
		inputs := []sharedtypes.FunctionInput{}
		for _, inputParam := range function.Input {
			// check if options is nil
			if inputParam.Options == nil {
				inputParam.Options = []string{}
			}
			inputs = append(inputs, sharedtypes.FunctionInput{
				Name:    inputParam.Name,
				Type:    inputParam.Type,
				GoType:  inputParam.GoType,
				Options: inputParam.Options,
			})
			// add the type to available types
			if AvailableTypes != nil && inputParam.GoType != "" && inputParam.GoType != "any" {
				AvailableTypes[inputParam.GoType] = true
			}
		}
		outputs := []sharedtypes.FunctionOutput{}
		for _, outputParam := range function.Output {
			outputs = append(outputs, sharedtypes.FunctionOutput{
				Name:   outputParam.Name,
				Type:   outputParam.Type,
				GoType: outputParam.GoType,
			})
			// add the type to available types
			if AvailableTypes != nil && outputParam.GoType != "" && outputParam.GoType != "any" {
				AvailableTypes[outputParam.GoType] = true
			}
		}

		// Save the function to internal states
		AvailableFunctions[function.Name] = &sharedtypes.FunctionDefinition{
			Name:             function.Name,
			FlowkitUrl:       url,
			ApiKey:           apiKey,
			DisplayName:      function.DisplayName,
			Description:      function.Description,
			Category:         function.Category,
			DeprecatedParams: function.DeprecatedParams,
			Inputs:           inputs,
			Outputs:          outputs,
			Type:             "go",
		}
		// add the category to available categories
		if AvailableCategories != nil && function.Category != "" {
			AvailableCategories[function.Category] = true
		}
	}

	return nil
}

// RunFunction calls the RunFunction gRPC and returns the outputs
// This function is used to run an external function
//
// Parameters:
//   - functionName: the name of the function to run
//   - inputs: the inputs to the function
//
// Returns:
//   - map[string]sharedtypes.FilledInputOutput: the outputs of the function
//   - error: an error message if the gRPC call fails
func RunFunction(ctx *logging.ContextMap, functionName string, inputs map[string]sharedtypes.FilledInputOutput) (outputs map[string]sharedtypes.FilledInputOutput, err error) {
	defer func() {
		r := recover()
		if r != nil {
			err = fmt.Errorf("panic occurred in RunFunction: %v", r)
		}
	}()

	// Get function definition
	functionDef, ok := AvailableFunctions[functionName]
	if !ok {
		return nil, fmt.Errorf("function %s not found in available functions", functionName)
	}

	// Set up a connection to the server.
	c, conn, err := createClient(functionDef.FlowkitUrl, functionDef.ApiKey)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to external function gRPC: %v", err)
	}
	defer conn.Close()

	// Create a context with a cancel
	ctxWithCancel, cancel := context.WithCancel(context.Background())
	defer cancel()

	// get logging metadata from context
	ctxWithMetadata, err := logging.CreateMetaDataFromCtx(ctx, ctxWithCancel)
	if err != nil {
		return nil, fmt.Errorf("error adding metadata: %v", err)
	}

	// Convert inputs to gRPC format based on order from function definition
	grpcInputs := []*aaliflowkitgrpc.FunctionInput{}
	for _, inputDef := range functionDef.Inputs {
		// create grpc input
		grpcInput := &aaliflowkitgrpc.FunctionInput{
			Name:   inputDef.Name,
			GoType: inputDef.GoType,
		}

		// Get the input value
		value, ok := inputs[inputDef.Name]
		if ok {
			// found: convert value to string
			stringValue, exists, err := typeconverters.ConvertGivenTypeToString(value.Value, inputDef.GoType)
			if err != nil {
				return nil, fmt.Errorf("error converting input '%s' to string: %v", inputDef.Name, err)
			}
			if !exists {
				return nil, fmt.Errorf("type '%s' does not exist in typeconverters.ConvertGivenTypeToString", inputDef.Name)
			}
			grpcInput.Value = stringValue

		} else {
			// input discrepancy, set to null value
			grpcInput.Value = ""
		}

		// Append the grpc input to the list
		grpcInputs = append(grpcInputs, grpcInput)
	}

	// Call RunFunction
	runResp, err := c.RunFunction(ctxWithMetadata, &aaliflowkitgrpc.FunctionInputs{
		Name:   functionName,
		Inputs: grpcInputs,
	})
	if err != nil {
		return nil, fmt.Errorf("error in external function gRPC RunFunction: %v", err)
	}

	// convert outputs to map[string]sharedtypes.FilledInputOutput
	outputs = map[string]sharedtypes.FilledInputOutput{}
	for _, output := range runResp.Outputs {
		// convert value to Go type
		value, exists, err := typeconverters.ConvertStringToGivenType(output.Value, output.GoType)
		if err != nil {
			return nil, fmt.Errorf("error converting output %s (%v) to Go type: %v", output.Name, output.Value, err)
		}
		if !exists {
			return nil, fmt.Errorf("type '%s' does not exist in typeconverters.ConvertStringToGivenType", output.Name)
		}

		// Save the output to the map
		outputs[output.Name] = sharedtypes.FilledInputOutput{
			Name:   output.Name,
			GoType: output.GoType,
			Value:  value,
		}
	}

	return outputs, nil
}

// StreamFunction calls the StreamFunction gRPC and returns a channel to stream the outputs
// This function is used to stream the outputs of an external function
//
// Parameters:
//   - functionName: the name of the function to run
//   - inputs: the inputs to the function
//
// Returns:
//   - *chan string: a channel to stream the output
//   - error: an error message if the gRPC call fails
func StreamFunction(ctx *logging.ContextMap, functionName string, inputs map[string]sharedtypes.FilledInputOutput) (channel *chan string, err error) {
	defer func() {
		r := recover()
		if r != nil {
			err = fmt.Errorf("panic occured in StreamFunction: %v", r)
		}
	}()

	// Get function definition
	functionDef, ok := AvailableFunctions[functionName]
	if !ok {
		return nil, fmt.Errorf("function %s not found in available functions", functionName)
	}

	// Set up a connection to the server.
	c, conn, err := createClient(functionDef.FlowkitUrl, functionDef.ApiKey)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to external function gRPC: %v", err)
	}

	// Create a context with a cancel
	ctxWithCancel, cancel := context.WithCancel(context.Background())

	// get logging metadata from context
	ctxWithMetadata, err := logging.CreateMetaDataFromCtx(ctx, ctxWithCancel)
	if err != nil {
		conn.Close()
		cancel()
		return nil, fmt.Errorf("error adding metadata: %v", err)
	}

	// Convert inputs to gRPC format based on order from function definition
	grpcInputs := []*aaliflowkitgrpc.FunctionInput{}
	for _, inputDef := range functionDef.Inputs {
		// create grpc input
		grpcInput := &aaliflowkitgrpc.FunctionInput{
			Name:   inputDef.Name,
			GoType: inputDef.GoType,
		}

		// Get the input value
		value, ok := inputs[inputDef.Name]
		if ok {
			// found: convert value to string
			stringValue, exists, err := typeconverters.ConvertGivenTypeToString(value.Value, inputDef.GoType)
			if err != nil {
				conn.Close()
				cancel()
				return nil, fmt.Errorf("error converting input %s to string: %v", inputDef.Name, err)
			}
			if !exists {
				conn.Close()
				cancel()
				return nil, fmt.Errorf("type '%s' does not exist in typeconverters.ConvertGivenTypeToString", inputDef.Name)
			}
			grpcInput.Value = stringValue

		} else {
			// input discrepancy, set to null value
			grpcInput.Value = ""
		}

		// Append the grpc input to the list
		grpcInputs = append(grpcInputs, grpcInput)
	}

	// Call StreamFunction
	stream, err := c.StreamFunction(ctxWithMetadata, &aaliflowkitgrpc.FunctionInputs{
		Name:   functionName,
		Inputs: grpcInputs,
	})
	if err != nil {
		conn.Close()
		cancel()
		return nil, fmt.Errorf("error in external function gRPC StreamFunction: %v", err)
	}

	// Create a stream channel
	streamChannel := make(chan string, 400)

	// Receive the stream from the server
	go receiveStreamFromServer(ctx, stream, &streamChannel, conn, cancel)

	return &streamChannel, nil
}

// receiveStreamFromServer receives the stream from the server and sends it to the channel
//
// Parameters:
//   - stream: the stream from the server
//   - streamChannel: the channel to send the stream to
func receiveStreamFromServer(ctx *logging.ContextMap, stream aaliflowkitgrpc.ExternalFunctions_StreamFunctionClient, streamChannel *chan string, conn *grpc.ClientConn, cancel context.CancelFunc) {
	defer func() {
		r := recover()
		if r != nil {
			logging.Log.Errorf(ctx, "Panic occured in receiveStreamFromServer: %v", r)
		}
	}()

	// Receive the stream from the server
	for {
		res, err := stream.Recv()
		if err != nil && err != io.EOF {
			logging.Log.Errorf(ctx, "error receiving stream: %v", err)
		}

		// Send the stream to the channel
		*streamChannel <- res.Value

		// end if isLast is true
		if res.IsLast {
			break
		}
	}

	// Close the channel
	conn.Close()
	cancel()
	close(*streamChannel)
}

// createClient creates a client to the external functions gRPC
//
// Returns:
//   - client: the client to the external functions gRPC
//   - connection: the connection to the external functions gRPC
//   - err: an error message if the client creation fails
func createClient(url string, apiKey string) (client aaliflowkitgrpc.ExternalFunctionsClient, connection *grpc.ClientConn, err error) {
	// Extract the scheme (http or https) from the EXTERNALFUNCTIONS_ENDPOINT
	var scheme string
	var address string
	switch {
	case strings.HasPrefix(url, "https://"):
		scheme = "https"
		address = strings.TrimPrefix(url, scheme+"://")
	case strings.HasPrefix(url, "http://"):
		scheme = "http"
		address = strings.TrimPrefix(url, scheme+"://")
	default:
		// legacy support for endpoint definition without http or https in front
		scheme = "http"
		address = url
	}

	// Set up the gRPC dial options
	var opts []grpc.DialOption

	// Add custom dialer with IPv4 first, fallback to IPv6
	opts = append(opts, grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
		d := &net.Dialer{}

		// Try IPv4 first
		conn, err := d.DialContext(ctx, "tcp4", addr)
		if err == nil {
			return conn, nil
		}

		// Fall back to IPv6 if IPv4 fails
		return d.DialContext(ctx, "tcp6", addr)
	}))

	if scheme == "https" {
		// Set up a secure connection with default TLS config
		creds := credentials.NewTLS(nil)
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		// Set up an insecure connection
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	// Add the API key if it is set
	if apiKey != "" {
		opts = append(opts, grpc.WithUnaryInterceptor(apiKeyInterceptor(apiKey)))
	}

	// Set max message size to 1GB
	opts = append(opts, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*1024)))

	// Set up a connection to the server
	conn, err := grpc.NewClient(address, opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to connect to external function gRPC: %v", err)
	}

	// Return the client
	c := aaliflowkitgrpc.NewExternalFunctionsClient(conn)
	return c, conn, nil
}

// apiKeyInterceptor is a gRPC client interceptor that adds an API key to the context metadata
// This interceptor is used to add the API key to the context metadata for all gRPC calls
//
// Parameters:
//   - apiKey: the API key to add to the context metadata
//
// Returns:
//   - grpc.UnaryClientInterceptor: the interceptor that adds the API key to the context metadata
func apiKeyInterceptor(apiKey string) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		// Get existing metadata from context (if any)
		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			// No existing metadata, create new
			md = metadata.MD{}
		} else {
			// Copy the metadata to avoid modifying the original
			md = md.Copy()
		}

		// Add API key to the existing metadata (this preserves other keys)
		md.Set("x-api-key", apiKey)

		// Create new context with MERGED metadata
		ctx = metadata.NewOutgoingContext(ctx, md)

		// Invoke the RPC with the modified context
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
