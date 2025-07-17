.. _sharedtypes_caveats:

====================
Integration Caveats
====================

.. grid:: 1
   :gutter: 2

   .. grid-item-card:: Common Pitfalls
      :class-card: sd-shadow-sm sd-rounded-md
      :text-align: left

      Watch for common issues when integrating SharedTypes into other services:

      - **Struct mismatch**: Ensure your proto definitions align with the corresponding SharedTypes structs.
      - **Import cycles**: Do not import service-specific logic back into SharedTypes.
      - **Marshalling errors**: Validate struct field tags and default values, especially for JSON or proto conversions.
      - **Missing fields in GRPC response**: Confirm that translation logic includes all necessary fields.