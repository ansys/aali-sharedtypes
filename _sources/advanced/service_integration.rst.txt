.. _service_integration:

===================
Service Integration
===================

.. grid:: 1
   :gutter: 2

   .. grid-item-card:: SharedTypes Across Services
      :class-card: sd-shadow-sm sd-rounded-md
      :text-align: left

      SharedTypes is used by Agent, Flowkit, LLM, and Chat. Each of these services imports the module to:

      - Reuse workflow and execution structs
      - Access common function metadata
      - Standardize input/output handling

      This approach eliminates duplication and ensures consistency between services.