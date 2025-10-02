// Copyright (C) 2025 ANSYS, Inc. and/or its affiliates.
// SPDX-License-Identifier: MIT

// Package aaliagentgrpc provides gRPC service definitions and protocol buffers for the AALI Agent.
//
// The AALI Agent is responsible for running workflows in the AALI distributed system.
// This package contains the generated Protocol Buffer and gRPC code from the aali-agent.proto file.
//
// Proto Definition
//
// The proto definition file can be found at:
// https://github.com/ansys/aali-sharedtypes/blob/main/pkg/aaliagentgrpc/aali-agent.proto
//
// Usage
//
// This package provides the WorkflowRun service which allows bidirectional streaming
// communication between clients and the AALI Agent server for workflow execution.
//
// Source Repository
//
// https://github.com/ansys/aali-sharedtypes
package aaliagentgrpc