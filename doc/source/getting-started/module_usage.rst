.. _sharedtypes_usage:

=============
Module Usage
=============

.. grid:: 1
   :gutter: 2

   .. grid-item-card:: Import SharedTypes
      :class-card: sd-shadow-sm sd-rounded-md
      :text-align: left

      Reference the shared structs from your Go code.

      .. code-block:: go

         import "github.com/ansys/aali-sharedtypes/pkg/sharedtypes"

   .. grid-item-card:: Example: SessionContext
      :class-card: sd-shadow-sm sd-rounded-md
      :text-align: left

      Create session context for workflow or exec connections.

      .. code-block:: go

         session := sharedtypes.SessionContext{
             SessionType: "workflow",
             WorkflowId:  "my-workflow",
             Variables: map[string]string{
                 "input_param": "value",
             },
         }

   .. grid-item-card:: Example: ExecRequest
      :class-card: sd-shadow-sm sd-rounded-md
      :text-align: left

      Create execution requests for code or flowkit operations.

      .. code-block:: go

         request := sharedtypes.ExecRequest{
             Type:            "code",
             Action:          "execute",
             InstructionGuid: "unique-id",
             ExecutionInstruction: &sharedtypes.ExecutionInstruction{
                 CodeType: "python",
                 Code:     []string{"print('Hello World')"},
             },
         }
