.. _sharedtypes_usage:

==============
Module Usage
==============

.. grid:: 1
   :gutter: 2

   .. grid-item-card:: Import SharedTypes
      :class-card: sd-shadow-sm sd-rounded-md
      :text-align: left

      Reference the shared structs from your Go code.

      .. code-block:: go

         import "github.com/ansys/aali-sharedtypes/pkg/sharedtypes"

   .. grid-item-card:: Example: ExecutionResult
      :class-card: sd-shadow-sm sd-rounded-md
      :text-align: left

      Create and populate shared types like `ExecutionResult`.

      .. code-block:: go

         result := sharedtypes.ExecutionResult{
             Status: "Success",
             Output: "Task completed",
         }

   .. grid-item-card:: Example: AgentConfig
      :class-card: sd-shadow-sm sd-rounded-md
      :text-align: left

      Use configuration structs from the same module.

      .. code-block:: go

         config := sharedtypes.AgentConfig{
             EnableTracing: true,
         }