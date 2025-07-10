.. _clients_flowkitpythonclient:

clients.flowkitpythonclient
===========================

This package provides functionality for clients.flowkitpythonclient.

Functions
---------

.. function:: ListFunctionsAndSaveToInteralStates

   func ListFunctionsAndSaveToInteralStates(url string, apiKey string) (err error)

   ListFunctionsAndSaveToInteralStates calls the FlowKit-Python API and saves the functions to internal states This function is used to get the list of available functions from the external function server and save them to internal states  Returns: - error: an error message if the API call fails


.. function:: RunFunction

   func RunFunction(functionName string, inputs map\[string\]sharedtypes.FilledInputOutput) (outputs map\[string\]sharedtypes.FilledInputOutput, err error)

   RunFunction calls the external function server and returns the outputs This function is used to run an external function  Parameters: - functionPath: the path of the function to run - inputs: the inputs to the function - outputDefinition: the definition of the outputs  Returns: - map\[string\]sharedtypes.FilledInputOutput: the outputs of the function - error: an error message if the API call fails
