.. _sharedtypes_setup:

=====
Setup
=====

**Add to your Go project:**

.. code-block:: bash

   go get github.com/ansys/aali-sharedtypes

**Import in your code:**

.. code-block:: go

   import "github.com/ansys/aali-sharedtypes/pkg/sharedtypes"

**Basic usage:**

.. code-block:: go

   request := sharedtypes.HandlerRequest{
       Adapter:         "chat",
       InstructionGuid: "unique-id",
       Data:            "Hello world",
   }
