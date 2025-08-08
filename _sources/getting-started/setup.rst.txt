.. _sharedtypes_setup:

=======================
Installation and Setup
=======================

.. grid:: 1
   :gutter: 2

   .. grid-item-card:: Clone the Repository
      :class-card: sd-shadow-sm sd-rounded-md
      :text-align: left

      Clone the source code for local development.

      .. code-block:: bash

         git clone https://github.com/ansys/aali-sharedtypes.git
         cd aali-sharedtypes

   .. grid-item-card:: Add as Dependency
      :class-card: sd-shadow-sm sd-rounded-md
      :text-align: left

      Use Go modules to install the sharedtypes package in your own service.

      .. code-block:: bash

         go get github.com/ansys/aali-sharedtypes@latest