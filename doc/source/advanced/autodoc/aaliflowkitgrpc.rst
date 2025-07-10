.. _aaliflowkitgrpc:

aaliflowkitgrpc
===============

This package provides functionality for aaliflowkitgrpc.

Functions
---------

.. function:: NewExternalFunctionsClient

   func NewExternalFunctionsClient(cc grpc.ClientConnInterface) ExternalFunctionsClient


.. function:: RegisterExternalFunctionsServer

   func RegisterExternalFunctionsServer(s grpc.ServiceRegistrar, srv ExternalFunctionsServer)


Types
-----

.. type:: ListFunctionsRequest

   type ListFunctionsRequest struct

   ListFunctionsRequest is the input message for the ListFunctions method. As no input is required, this message is empty.


.. type:: ListFunctionsResponse

   type ListFunctionsResponse struct

   ListFunctionsResponse is the output message for the ListFunctions method. It contains a map of function names to their definitions.


.. type:: FunctionDefinition

   type FunctionDefinition struct

   FunctionDefinition is the definition of an individual function. It contains the name, description, inputs and outputs of the function.


.. type:: FunctionInputDefinition

   type FunctionInputDefinition struct

   FunctionInputDefinition is the definition of an input for a function. It contains the name, type, Go language type and options for the input.


.. type:: FunctionOutputDefinition

   type FunctionOutputDefinition struct

   FunctionOutputDefinition is the definition of an output for a function. It contains the name, type and Go language type of the output.


.. type:: FunctionInputs

   type FunctionInputs struct

   FunctionInputs is the input message for the RunFunction method. It contains the name of the function to run and a list of inputs.


.. type:: FunctionInput

   type FunctionInput struct

   Single input for a function.


.. type:: FunctionOutputs

   type FunctionOutputs struct

   FunctionOutputs is the output message for the RunFunction method. It contains the name of the function that was run and a list of outputs.


.. type:: FunctionOutput

   type FunctionOutput struct

   FunctionOutput is a single output from a function. It contains the name, Go language type and value of the output.


.. type:: StreamOutput

   type StreamOutput struct

   StreamOutput is the output message for the StreamFunction method. It contains the message counter, a flag indicating if this is the last message


.. type:: ExternalFunctionsClient

   type ExternalFunctionsClient interface

   ExternalFunctionsClient is the client API for ExternalFunctions service.  For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab\=doc\#ClientConn.NewStream.  ExternalFunctions is a gRPC service that allows for listing and running the function available in the externalfunctions package.


.. type:: ExternalFunctions_StreamFunctionClient

   type ExternalFunctions\_StreamFunctionClient \= grpc.ServerStreamingClient\[StreamOutput\]

   This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.


.. type:: ExternalFunctionsServer

   type ExternalFunctionsServer interface

   ExternalFunctionsServer is the server API for ExternalFunctions service. All implementations must embed UnimplementedExternalFunctionsServer for forward compatibility.  ExternalFunctions is a gRPC service that allows for listing and running the function available in the externalfunctions package.


.. type:: UnimplementedExternalFunctionsServer

   type UnimplementedExternalFunctionsServer struct

   UnimplementedExternalFunctionsServer must be embedded to have forward compatible implementations.  NOTE: this should be embedded by value instead of pointer to avoid a nil pointer dereference when methods are called.


.. type:: UnsafeExternalFunctionsServer

   type UnsafeExternalFunctionsServer interface

   UnsafeExternalFunctionsServer may be embedded to opt out of forward compatibility for this service. Use of this interface is not recommended, as added methods to ExternalFunctionsServer will result in compilation errors.


.. type:: ExternalFunctions_StreamFunctionServer

   type ExternalFunctions\_StreamFunctionServer \= grpc.ServerStreamingServer\[StreamOutput\]

   This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.


Interfaces
----------

.. type:: ExternalFunctionsClient

   ExternalFunctionsClient is the client API for ExternalFunctions service.  For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab\=doc\#ClientConn.NewStream.  ExternalFunctions is a gRPC service that allows for listing and running the function available in the externalfunctions package.

   **Methods:**

   * ListFunctions(ctx context.Context, in \*ListFunctionsRequest, opts ...grpc.CallOption) (\*ListFunctionsResponse, error)
   * RunFunction(ctx context.Context, in \*FunctionInputs, opts ...grpc.CallOption) (\*FunctionOutputs, error)
   * StreamFunction(ctx context.Context, in \*FunctionInputs, opts ...grpc.CallOption) (grpc.ServerStreamingClient\[StreamOutput\], error)


.. type:: ExternalFunctionsServer

   ExternalFunctionsServer is the server API for ExternalFunctions service. All implementations must embed UnimplementedExternalFunctionsServer for forward compatibility.  ExternalFunctions is a gRPC service that allows for listing and running the function available in the externalfunctions package.

   **Methods:**

   * ListFunctions(context.Context, \*ListFunctionsRequest) (\*ListFunctionsResponse, error)
   * RunFunction(context.Context, \*FunctionInputs) (\*FunctionOutputs, error)
   * StreamFunction(\*FunctionInputs, grpc.ServerStreamingServer\[StreamOutput\]) error
   * mustEmbedUnimplementedExternalFunctionsServer()


.. type:: UnsafeExternalFunctionsServer

   UnsafeExternalFunctionsServer may be embedded to opt out of forward compatibility for this service. Use of this interface is not recommended, as added methods to ExternalFunctionsServer will result in compilation errors.

   **Methods:**

   * mustEmbedUnimplementedExternalFunctionsServer()


Constants
---------

.. data:: ExternalFunctions_ListFunctions_FullMethodName

   const ExternalFunctions_ListFunctions_FullMethodName = "/aaliflowkitgrpc.ExternalFunctions/ListFunctions"


.. data:: ExternalFunctions_RunFunction_FullMethodName

   const ExternalFunctions_RunFunction_FullMethodName = "/aaliflowkitgrpc.ExternalFunctions/RunFunction"


.. data:: ExternalFunctions_StreamFunction_FullMethodName

   const ExternalFunctions_StreamFunction_FullMethodName = "/aaliflowkitgrpc.ExternalFunctions/StreamFunction"


Variables
---------

.. data:: File_pkg_aaliflowkitgrpc_aali_flowkit_proto

   var File\_pkg\_aaliflowkitgrpc\_aali\_flowkit\_proto protoreflect.FileDescriptor
