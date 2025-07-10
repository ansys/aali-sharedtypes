.. _aaliagentgrpc:

aaliagentgrpc
=============

This package provides functionality for aaliagentgrpc.

Functions
---------

.. function:: NewWorkflowRunClient

   func NewWorkflowRunClient(cc grpc.ClientConnInterface) WorkflowRunClient


.. function:: RegisterWorkflowRunServer

   func RegisterWorkflowRunServer(s grpc.ServiceRegistrar, srv WorkflowRunServer)


Types
-----

.. type:: ClientMessage

   type ClientMessage struct

   ClientMessage is the message sent by the client to the server.


.. type:: ClientMessage_SessionContext

   type ClientMessage\_SessionContext struct


.. type:: ClientMessage_ClientRequest

   type ClientMessage\_ClientRequest struct


.. type:: ServerMessage

   type ServerMessage struct

   ServerMessage is the message sent by the server to the client.


.. type:: ServerMessage_ConnectionStatus

   type ServerMessage\_ConnectionStatus struct


.. type:: ServerMessage_AuthenticationStatus

   type ServerMessage\_AuthenticationStatus struct


.. type:: ServerMessage_ClientResponse

   type ServerMessage\_ClientResponse struct


.. type:: SessionContext

   type SessionContext struct

   SessionContext is the message to initiate a session with the server.


.. type:: ClientRequest

   type ClientRequest struct

   ClientRequest is the message to send a request to the server.


.. type:: WorkflowFeedback

   type WorkflowFeedback struct

   WorkflowFeedback is the message to send feedback to the server.


.. type:: ConnectionStatus

   type ConnectionStatus struct

   ConnectionStatus is the message to indicate the connection status after client sends a session context message.


.. type:: AuthenticationStatus

   type AuthenticationStatus struct

   AuthenticationStatus is the message to indicate failing authentication after client sends a session context message with authentication enabled.


.. type:: ClientResponse

   type ClientResponse struct

   ClientResponse is the message to send a response to the client.


.. type:: ErrorResponse

   type ErrorResponse struct

   ErrorResponse is the message to send an error response to the client.


.. type:: WorkflowRunClient

   type WorkflowRunClient interface

   WorkflowRunClient is the client API for WorkflowRun service.  For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab\=doc\#ClientConn.NewStream.  WorkflowRun is a gRPC service that allows for running a workflow.


.. type:: WorkflowRun_RunWorkflowClient

   type WorkflowRun\_RunWorkflowClient \= grpc.BidiStreamingClient\[ClientMessage, ServerMessage\]

   This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.


.. type:: WorkflowRunServer

   type WorkflowRunServer interface

   WorkflowRunServer is the server API for WorkflowRun service. All implementations must embed UnimplementedWorkflowRunServer for forward compatibility.  WorkflowRun is a gRPC service that allows for running a workflow.


.. type:: UnimplementedWorkflowRunServer

   type UnimplementedWorkflowRunServer struct

   UnimplementedWorkflowRunServer must be embedded to have forward compatible implementations.  NOTE: this should be embedded by value instead of pointer to avoid a nil pointer dereference when methods are called.


.. type:: UnsafeWorkflowRunServer

   type UnsafeWorkflowRunServer interface

   UnsafeWorkflowRunServer may be embedded to opt out of forward compatibility for this service. Use of this interface is not recommended, as added methods to WorkflowRunServer will result in compilation errors.


.. type:: WorkflowRun_RunWorkflowServer

   type WorkflowRun\_RunWorkflowServer \= grpc.BidiStreamingServer\[ClientMessage, ServerMessage\]

   This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.


Interfaces
----------

.. type:: WorkflowRunClient

   WorkflowRunClient is the client API for WorkflowRun service.  For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab\=doc\#ClientConn.NewStream.  WorkflowRun is a gRPC service that allows for running a workflow.

   **Methods:**

   * RunWorkflow(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient\[ClientMessage, ServerMessage\], error)


.. type:: WorkflowRunServer

   WorkflowRunServer is the server API for WorkflowRun service. All implementations must embed UnimplementedWorkflowRunServer for forward compatibility.  WorkflowRun is a gRPC service that allows for running a workflow.

   **Methods:**

   * RunWorkflow(grpc.BidiStreamingServer\[ClientMessage, ServerMessage\]) error
   * mustEmbedUnimplementedWorkflowRunServer()


.. type:: UnsafeWorkflowRunServer

   UnsafeWorkflowRunServer may be embedded to opt out of forward compatibility for this service. Use of this interface is not recommended, as added methods to WorkflowRunServer will result in compilation errors.

   **Methods:**

   * mustEmbedUnimplementedWorkflowRunServer()


Constants
---------

.. data:: WorkflowRun_RunWorkflow_FullMethodName

   const WorkflowRun_RunWorkflow_FullMethodName = "/aaliagentgrpc.WorkflowRun/RunWorkflow"


Variables
---------

.. data:: File_pkg_aaliagentgrpc_aali_agent_proto

   var File\_pkg\_aaliagentgrpc\_aali\_agent\_proto protoreflect.FileDescriptor
