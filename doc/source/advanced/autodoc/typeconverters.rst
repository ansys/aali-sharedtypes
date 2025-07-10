.. _typeconverters:

typeconverters
==============

This package provides functionality for typeconverters.

Functions
---------

.. function:: TestJSONToGo

   func TestJSONToGo(t \*testing.T) 


.. function:: TestGoToJSON

   func TestGoToJSON(t \*testing.T) 


.. function:: TestConvertStringToGivenType

   func TestConvertStringToGivenType(t \*testing.T) 


.. function:: TestDeepCopy

   func TestDeepCopy(t \*testing.T) 


.. function:: JSONToGo

   func JSONToGo(jsonType string) (string, error)

   JSONToGo converts a JSON data type to a Go data type.  Parameters:  jsonType: The JSON data type to convert.  Returns:  string: The Go data type. error: An error if the JSON data type is not supported.


.. function:: GoToJSON

   func GoToJSON(goType string) string

   GoToJSON converts a Go data type to a JSON data type.  Parameters:  goType: The Go data type to convert.  Returns:  string: The JSON data type.


.. function:: ConvertStringToGivenType

   func ConvertStringToGivenType(value string, goType string) (output interface

   ConvertStringToGivenType converts a string to a given Go type.  Parameters: - value: a string containing the value to convert - goType: a string containing the Go type to convert to  Returns: - output: an interface containing the converted value - err: an error containing the error message


.. function:: ConvertGivenTypeToString

   func ConvertGivenTypeToString(value interface\{\}, goType string) (output string, err error)

   ConvertGivenTypeToString converts a given Go type to a string.  Parameters: - value: an interface containing the value to convert - goType: a string containing the Go type to convert from  Returns: - string: a string containing the converted value - err: an error containing the error message


.. function:: DeepCopy

   func DeepCopy(src, dst interface\{\}) (err error)

   DeepCopy deep copies the source interface to the destination interface.  Parameters: - src: an interface containing the source - dst: an interface containing the destination  Returns: - err: an error containing the error message

