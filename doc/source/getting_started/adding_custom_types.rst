Adding Custom Types
===================

When creating custom FlowKit functions, you may need to define custom data types that aren't already available in AALI Shared Types. This guide walks you through the process of adding new types and getting them integrated into AALI FlowKit and Agent.

When to Add Custom Types
------------------------

You need to add custom types when:

- Your FlowKit function requires a specific data structure not available in existing types
- You need to pass complex data between functions
- You're integrating with external systems that use specific data formats
- You want to ensure type safety and validation for your function inputs/outputs

Step 1: Navigate to Type Definitions
-------------------------------------

All type definitions are located in ``pkg/sharedtypes/``. Each file represents a category of types:

.. code-block:: text

   pkg/sharedtypes/
   ├── agent.go           # Agent-related types
   ├── exec.go            # Execution context types
   ├── functiondefinition.go  # Function definition types
   ├── knowledgedb.go     # Knowledge database types
   ├── llmhandler.go      # LLM handler types
   ├── ansysgpt.go        # Ansys GPT specific types
   ├── ansysmaterials.go  # Ansys Materials types
   ├── ansysdiscovery.go  # Ansys Discovery types
   └── dataextraction.go  # Data extraction types

Step 2: Choose or Create a Category
------------------------------------

Decide where your custom type belongs:

- If it fits an existing category, add it to the appropriate file
- If it's a new domain, create a new file following the naming pattern

Example: adding a custom type to an existing category:

.. code-block:: go

   // In pkg/sharedtypes/functiondefinition.go

   // CustomAnalysisResult is an EXAMPLE of a custom type you could add
   // This type does not exist in the codebase - it's shown here as a template
   type CustomAnalysisResult struct {
       AnalysisID   string                 `json:"analysis_id"`
       Status       string                 `json:"status"`
       Results      map[string]interface{} `json:"results"`
       Timestamp    time.Time              `json:"timestamp"`
       Confidence   float64                `json:"confidence"`
   }

Step 3: Define Your Custom Type
-------------------------------

Follow these guidelines when defining your type:

1. Use clear, descriptive names
2. Add JSON tags for serialization
3. Include godoc comments
4. Consider validation requirements

.. code-block:: go

   // MyCustomType represents data for a specific FlowKit function
   type MyCustomType struct {
       // ID is a unique identifier for this instance
       ID string `json:"id"`

       // Name is a human-readable name
       Name string `json:"name"`

       // Data contains the actual payload
       Data map[string]interface{} `json:"data"`

       // ProcessingOptions defines how this data should be processed
       ProcessingOptions ProcessingConfig `json:"processing_options,omitempty"`
   }

   // ProcessingConfig defines options for processing MyCustomType
   type ProcessingConfig struct {
       EnableValidation bool   `json:"enable_validation"`
       MaxRetries       int    `json:"max_retries"`
       TimeoutSeconds   int    `json:"timeout_seconds"`
   }

Step 4: Add Type Converters
----------------------------

Navigate to ``pkg/typeconverters/typeconverters.go`` and add conversion support for your type:

.. code-block:: go

   // In the ConvertStringToGivenType function, add a case for your type:
   case "MyCustomType":
       var result sharedtypes.MyCustomType
       err := json.Unmarshal([]byte(s), &result)
       if err != nil {
           return nil, err
       }
       return result, nil

   // In the ConvertGivenTypeToString function, add a case:
   case sharedtypes.MyCustomType:
       bytes, err := json.Marshal(v)
       if err != nil {
           return "", err
       }
       return string(bytes), nil


Step 5: Build and Verify
------------------------

Build the module to ensure your changes compile:

.. code-block:: bash

   go build ./...
   go test ./...

Step 6: Submit Your Changes
---------------------------

Once your custom type is working:

1. **Commit your changes** to a feature branch
2. **Create a pull request** to merge into the main branch
3. **After merge**, the shared types need to be updated in:

   - AALI FlowKit: Import the latest shared types version
   - AALI Agent: Import the latest shared types version

4. **Coordinate with the team** to ensure both services are updated

Example: Using Your Custom Type in FlowKit
-------------------------------------------

After your type is integrated, you can use it in FlowKit function definitions.
This example shows how you would use the custom types defined earlier:

.. code-block:: go

   // Example function using the custom types defined in this guide
   func MyCustomFunction(input sharedtypes.MyCustomType) (sharedtypes.CustomAnalysisResult, error) {
       // Process the custom type
       result := sharedtypes.CustomAnalysisResult{
           AnalysisID: generateID(),
           Status:     "completed",
           Results:    processData(input.Data),
           Timestamp:  time.Now(),
           Confidence: 0.95,
       }

       return result, nil
   }

Best Practices
--------------

- **Keep types focused**: Each type should have a single, clear purpose
- **Use standard Go conventions**: Follow Go naming and structure guidelines
- **Document thoroughly**: Include examples in comments
- **Consider backward compatibility**: Changes to existing types can break other services
- **Test edge cases**: Ensure your type handles null/empty values appropriately

Next Steps
----------

- Explore existing types in the ``pkg/sharedtypes/`` directory
- Learn about type conversion in ``pkg/typeconverters/typeconverters.go``
- Review gRPC definitions in ``pkg/aaliagentgrpc/`` and ``pkg/aaliflowkitgrpc/``
