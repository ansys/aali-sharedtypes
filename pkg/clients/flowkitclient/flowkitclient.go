package flowkitclient

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/ansys/allie-sharedtypes/pkg/allieflowkitgrpc"
	"github.com/ansys/allie-sharedtypes/pkg/logging"
	"github.com/ansys/allie-sharedtypes/pkg/structs"
	"github.com/ansys/allie-sharedtypes/pkg/typeconverters"

	"github.com/ansys/allie-sharedtypes/pkg/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// Global variable to store the available functions
var AvailableFunctions map[string]*structs.FunctionDefinition

// ListFunctionsAndSaveToInteralStates calls the ListFunctions gRPC and saves the functions to internal states
// This function is used to get the list of available functions from the external function server
// and save them to internal states
//
// Returns:
//   - error: an error message if the gRPC call fails
func ListFunctionsAndSaveToInteralStates(ctx *logging.ContextMap) error {
	defer func() {
		r := recover()
		if r != nil {
			logging.Log.Errorf(ctx, "Panic occured in ListFunctionsAndSaveToInteralStates: %v", r)
		}
	}()

	// Set up a connection to the server.
	c, conn, err := createClient()
	if err != nil {
		return fmt.Errorf("unable to connect to external function gRPC: %v", err)
	}
	defer conn.Close()

	// Create a context with a cancel
	ctxWithCancel, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Call ListFunctions
	listResp, err := c.ListFunctions(ctxWithCancel, &allieflowkitgrpc.ListFunctionsRequest{})
	if err != nil {
		return fmt.Errorf("error in external function gRPC ListFunctions: %v", err)
	}

	// Save the functions to internal states
	for _, function := range listResp.Functions {
		// convert inputs and outputs
		inputs := []structs.FunctionInput{}
		for _, inputParam := range function.Input {
			// check if options is nil
			if inputParam.Options == nil {
				inputParam.Options = []string{}
			}
			inputs = append(inputs, structs.FunctionInput{
				Name:    inputParam.Name,
				Type:    inputParam.Type,
				GoType:  inputParam.GoType,
				Options: inputParam.Options,
			})
		}
		outputs := []structs.FunctionOutput{}
		for _, outputParam := range function.Output {
			outputs = append(outputs, structs.FunctionOutput{
				Name:   outputParam.Name,
				Type:   outputParam.Type,
				GoType: outputParam.GoType,
			})
		}

		// Save the function to internal states
		AvailableFunctions[function.Name] = &structs.FunctionDefinition{
			Name:        function.Name,
			DisplayName: function.DisplayName,
			Description: function.Description,
			Category:    function.Category,
			Inputs:      inputs,
			Outputs:     outputs,
			Type:        "go",
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
//   - map[string]structs.FilledInputOutput: the outputs of the function
//   - error: an error message if the gRPC call fails
func RunFunction(ctx *logging.ContextMap, functionName string, inputs map[string]structs.FilledInputOutput) (map[string]structs.FilledInputOutput, error) {
	defer func() {
		r := recover()
		if r != nil {
			logging.Log.Errorf(ctx, "Panic occured in RunFunction: %v", r)
		}
	}()

	// Set up a connection to the server.
	c, conn, err := createClient()
	if err != nil {
		return nil, fmt.Errorf("unable to connect to external function gRPC: %v", err)
	}
	defer conn.Close()

	// Create a context with a cancel
	ctxWithCancel, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Get function definition
	functionDef, ok := AvailableFunctions[functionName]
	if !ok {
		return nil, fmt.Errorf("function %s not found in available functions", functionName)
	}

	// Convert inputs to gRPC format based on order from function definition
	input := []*allieflowkitgrpc.FunctionInput{}
	for _, inputDef := range functionDef.Inputs {
		// Get the input value
		value, ok := inputs[inputDef.Name]
		if !ok {
			return nil, fmt.Errorf("Input %s not found in inputs", inputDef.Name)
		}

		// convert value to string
		stringValue, err := typeconverters.ConvertGivenTypeToString(value.Value, inputDef.GoType)
		if err != nil {
			return nil, fmt.Errorf("Error converting input %s to string: %v", inputDef.Name, err)
		}

		// Convert the input value to gRPC format
		input = append(input, &allieflowkitgrpc.FunctionInput{
			Name:   inputDef.Name,
			GoType: inputDef.GoType,
			Value:  stringValue,
		})
	}

	// Call RunFunction
	runResp, err := c.RunFunction(ctxWithCancel, &allieflowkitgrpc.FunctionInputs{
		Name:   functionName,
		Inputs: input,
	})
	if err != nil {
		return nil, fmt.Errorf("error in external function gRPC RunFunction: %v", err)
	}

	// convert outputs to map[string]structs.FilledInputOutput
	outputs := map[string]structs.FilledInputOutput{}
	for _, output := range runResp.Outputs {
		// convert value to Go type
		value, err := typeconverters.ConvertStringToGivenType(output.Value, output.GoType)
		if err != nil {
			return nil, fmt.Errorf("error converting output %s to Go type: %v", output.Name, err)
		}

		// Save the output to the map
		outputs[output.Name] = structs.FilledInputOutput{
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
func StreamFunction(ctx *logging.ContextMap, functionName string, inputs map[string]structs.FilledInputOutput) (*chan string, error) {
	defer func() {
		r := recover()
		if r != nil {
			logging.Log.Errorf(ctx, "Panic occured in RunFunction: %v", r)
		}
	}()

	// Set up a connection to the server.
	c, conn, err := createClient()
	if err != nil {
		return nil, fmt.Errorf("unable to connect to external function gRPC: %v", err)
	}

	// Create a context with a cancel
	ctxWithCancel, cancel := context.WithCancel(context.Background())

	// Get function definition
	functionDef, ok := AvailableFunctions[functionName]
	if !ok {
		conn.Close()
		cancel()
		return nil, fmt.Errorf("function %s not found in available functions", functionName)
	}

	// Convert inputs to gRPC format based on order from function definition
	input := []*allieflowkitgrpc.FunctionInput{}
	for _, inputDef := range functionDef.Inputs {
		// Get the input value
		value, ok := inputs[inputDef.Name]
		if !ok {
			conn.Close()
			cancel()
			return nil, fmt.Errorf("input %s not found in inputs", inputDef.Name)
		}

		// convert value to string
		stringValue, err := typeconverters.ConvertGivenTypeToString(value.Value, inputDef.GoType)
		if err != nil {
			conn.Close()
			cancel()
			return nil, fmt.Errorf("error converting input %s to string: %v", inputDef.Name, err)
		}

		// Convert the input value to gRPC format
		input = append(input, &allieflowkitgrpc.FunctionInput{
			Name:   inputDef.Name,
			GoType: inputDef.GoType,
			Value:  stringValue,
		})
	}

	// Call StreamFunction
	stream, err := c.StreamFunction(ctxWithCancel, &allieflowkitgrpc.FunctionInputs{
		Name:   functionName,
		Inputs: input,
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
func receiveStreamFromServer(ctx *logging.ContextMap, stream allieflowkitgrpc.ExternalFunctions_StreamFunctionClient, streamChannel *chan string, conn *grpc.ClientConn, cancel context.CancelFunc) {
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
func createClient() (client allieflowkitgrpc.ExternalFunctionsClient, connection *grpc.ClientConn, err error) {
	// Extract the scheme (http or https) from the EXTERNALFUNCTIONS_ENDPOINT
	var scheme string
	var address string
	switch {
	case strings.HasPrefix(config.GlobalConfig.EXTERNALFUNCTIONS_ENDPOINT, "https://"):
		scheme = "https"
		address = strings.TrimPrefix(config.GlobalConfig.EXTERNALFUNCTIONS_ENDPOINT, scheme+"://")
	case strings.HasPrefix(config.GlobalConfig.EXTERNALFUNCTIONS_ENDPOINT, "http://"):
		scheme = "http"
		address = strings.TrimPrefix(config.GlobalConfig.EXTERNALFUNCTIONS_ENDPOINT, scheme+"://")
	default:
		// legacy support for endpoint definition without http or https in front
		scheme = "http"
		address = config.GlobalConfig.EXTERNALFUNCTIONS_ENDPOINT
	}

	// Set up the gRPC dial options
	var opts []grpc.DialOption
	if scheme == "https" {
		// Set up a secure connection with default TLS config
		creds := credentials.NewTLS(nil)
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		// Set up an insecure connection
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	// Add the API key if it is set
	if config.GlobalConfig.FLOWKIT_API_KEY != "" {
		opts = append(opts, grpc.WithUnaryInterceptor(apiKeyInterceptor(config.GlobalConfig.FLOWKIT_API_KEY)))
	}

	// Set max message size to 1GB
	opts = append(opts, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*1024)))

	// Set up a connection to the server.
	conn, err := grpc.NewClient(address, opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to connect to external function gRPC: %v", err)
	}

	// Return the client
	c := allieflowkitgrpc.NewExternalFunctionsClient(conn)
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
		// Add API key to the context metadata
		md := metadata.Pairs("x-api-key", apiKey)
		ctx = metadata.NewOutgoingContext(ctx, md)

		// Invoke the RPC with the modified context
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
