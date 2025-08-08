.. _api_reference:

API Reference
=============

.. grid:: 1
   :gutter: 2

   .. grid-item-card:: Public Structs and Packages
      :class-card: sd-shadow-sm sd-rounded-md
      :text-align: left

      SharedTypes exposes the following key packages:

      - `pkg/sharedtypes` — core workflow, function, and result types
      - `pkg/config` — runtime config definitions
      - `pkg/aali_graphdb` — logical types and graph values

   .. grid-item-card:: Example Struct: ExecutionResult
      :class-card: sd-shadow-sm sd-rounded-md
      :text-align: left

      .. code-block:: go

         type ExecutionResult struct {
             Status string
             Output string
             Error  *ErrorReport
         }

      Returned by agents and services to describe outcome of a task.

   .. grid-item-card:: Example Struct: FunctionDefinition
      :class-card: sd-shadow-sm sd-rounded-md
      :text-align: left

      .. code-block:: go

         type FunctionDefinition struct {
             Name        string
             Description string
             Input       []FunctionInput
             Output      []FunctionOutput
         }

      Used to describe callable functions shared across services.

.. note::
   This page is a placeholder for future auto-generated Go API documentation.