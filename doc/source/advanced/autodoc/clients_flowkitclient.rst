.. _clients_flowkitclient:

clients.flowkitclient
=====================

This package provides functionality for clients.flowkitclient.

Functions
---------

.. function:: ListFunctionsAndSaveToInteralStates

   func ListFunctionsAndSaveToInteralStates(url string, apiKey string) (err error)

   ListFunctionsAndSaveToInteralStates calls the ListFunctions gRPC and saves the functions to internal states This function is used to get the list of available functions from the external function server and save them to internal states  Returns: - error: an error message if the gRPC call fails


.. function:: RunFunction

   func RunFunction(functionName string, inputs map\[string\]sharedtypes.FilledInputOutput) (outputs map\[string\]sharedtypes.FilledInputOutput, err error)

   RunFunction calls the RunFunction gRPC and returns the outputs This function is used to run an external function  Parameters: - functionName: the name of the function to run - inputs: the inputs to the function  Returns: - map\[string\]sharedtypes.FilledInputOutput: the outputs of the function - error: an error message if the gRPC call fails


.. function:: StreamFunction

   func StreamFunction(ctx \*logging.ContextMap, functionName string, inputs map\[string\]sharedtypes.FilledInputOutput) (channel \*chan string, err error)

   StreamFunction calls the StreamFunction gRPC and returns a channel to stream the outputs This function is used to stream the outputs of an external function  Parameters: - functionName: the name of the function to run - inputs: the inputs to the function  Returns: - \*chan string: a channel to stream the output - error: an error message if the gRPC call fails


Variables
---------

.. data:: AvailableFunctions

   var AvailableFunctions map\[string\]\*sharedtypes.FunctionDefinition

   Global variable to store the available functions

