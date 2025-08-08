.. _type_system:

==================
Type System Design
==================

.. grid:: 1
   :gutter: 2

   .. grid-item-card:: Design Principles
      :class-card: sd-shadow-sm sd-rounded-md
      :text-align: left

      The type system in SharedTypes follows these principles:

      - Avoid tight coupling to any backend (e.g., GRPC, protobuf)
      - Keep structs small and focused
      - Use typed enums and constants where possible
      - Avoid deep nesting

      This makes it easier to test, refactor, and extend data models across the system.