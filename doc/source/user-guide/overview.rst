.. _sharedtypes_overview:

========
Overview
========

AALI SharedTypes provides common Go data structures used across AALI services to ensure consistency and interoperability.

Key Packages
------------

**pkg/sharedtypes**
   Core request/response types for LLM Handler, Agent communication

**pkg/config**
   Configuration management utilities

**pkg/logging**
   Structured logging with context support

**pkg/aali_graphdb**
   Graph database value types and logical operations

**pkg/typeconverters**
   Utilities for converting between different type representations

Common Usage Patterns
----------------------

**LLM Handler Requests:**

.. code-block:: go

   import "github.com/ansys/aali-sharedtypes/pkg/sharedtypes"

   request := sharedtypes.HandlerRequest{
       Adapter:         "chat",
       InstructionGuid: "request-123",
       Data:            "Your question here",
       ChatRequestType: "summary",
   }

**Configuration Loading:**

.. code-block:: go

   import "github.com/ansys/aali-sharedtypes/pkg/config"

   config.InitConfig([]string{}, map[string]interface{}{
       "SERVICE_NAME": "my-service",
       "LOG_LEVEL":    "info",
   })

**Structured Logging:**

.. code-block:: go

   import "github.com/ansys/aali-sharedtypes/pkg/logging"

   ctx := &logging.ContextMap{}
   logging.Log.Infof(ctx, "Processing request: %s", requestID)

**Function Definitions:**

.. code-block:: go

   import "github.com/ansys/aali-sharedtypes/pkg/sharedtypes"

   funcDef := sharedtypes.FunctionDefinition{
       Name:        "calculate_sum",
       DisplayName: "Calculate Sum",
       Description: "Adds two numbers together",
       Category:    "generic",
       Type:        "go",
       Inputs: []sharedtypes.FunctionInput{
           {Name: "a", Type: "number", GoType: "float64"},
           {Name: "b", Type: "number", GoType: "float64"},
       },
       Outputs: []sharedtypes.FunctionOutput{
           {Name: "result", Type: "number", GoType: "float64"},
       },
   }

When to Use
-----------

- **Building AALI services** - Use these types to maintain compatibility
- **Integrating with LLM Handler** - Use HandlerRequest/HandlerResponse
- **Working with Agent** - Use SessionContext and related types
- **Managing configuration** - Use the config package for consistent setup

For complete type definitions and function signatures, see the :doc:`../advanced/autodoc/index`.
