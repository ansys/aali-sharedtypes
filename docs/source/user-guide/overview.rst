.. _sharedtypes_overview:

========
Overview
========

.. grid:: 1
   :gutter: 2

   .. grid-item-card:: Shared Data Models
      :class-card: sd-shadow-sm sd-rounded-md
      :text-align: left

      SharedTypes provides a central set of Go structs used across AALI components. These include workflows, execution results, function definitions, graph values, and configuration types. Structs are designed to work independently from protobufs and GRPC.

      This module is imported by Agent, Flowkit, LLM, and others to maintain a consistent internal representation of logic and data.