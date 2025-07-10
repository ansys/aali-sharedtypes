.. _grpc_mapping:

=======================
GRPC Mapping
=======================

.. grid:: 1
   :gutter: 2

   .. grid-item-card:: Struct-to-Proto Translation
      :class-card: sd-shadow-sm sd-rounded-md
      :text-align: left

      Many structs are translated into proto messages before being sent across GRPC. Examples:

      - `HandlerRequest` → `llmpb.HandlerRequest`
      - `FunctionDefinition` → `flowkitpb.FunctionMeta`
      - `SessionContext` → `agentpb.SessionContext`

      The conversion logic is implemented in each service, keeping SharedTypes independent.
