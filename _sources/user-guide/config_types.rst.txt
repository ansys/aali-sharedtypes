.. _config_types:

===================
Configuration Types
===================

.. grid:: 1
   :gutter: 2

   .. grid-item-card:: Shared Config Structs
      :class-card: sd-shadow-sm sd-rounded-md
      :text-align: left

      Configuration models defined in SharedTypes are used to set runtime behavior. This includes:

      - `AgentConfig`
      - `LoggingConfig`
      - `LLMConfig`

      These structs are reused by various services to avoid duplicating config logic and ensure consistent setup.