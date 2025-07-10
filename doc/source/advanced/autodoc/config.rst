.. _config:

config
======

This package provides functionality for config.

Functions
---------

.. function:: InitConfig

   func InitConfig(requiredProperties \[\]string, optionalDefaultValues map\[string\]interface\{\})

   ////////////////////////////////////////// Standard Config init for Aali Go Modules ////////////////////////////////////////// InitConfig initializes the configuration for the Aali service.  Parameters: - requiredProperties: The list of required properties. - optionalDefaultValues: The map of optional properties and their default values.


.. function:: InitGlobalConfigFromFile

   func InitGlobalConfigFromFile(fileName string, requiredProperties \[\]string, optionalDefaultValues map\[string\]interface\{\}) (err error)

   //////////////////////////////////////// Read Config variables from Config file //////////////////////////////////////// InitGlobalConfigFromFile reads the configuration file and initializes the Config object.  Parameters: - fileName: The name of the configuration file. - requiredProperties: The list of required properties. - optionalDefaultValues: The map of optional properties and their default values.  Returns: - err: An error if there was an issue initializing the configuration.


.. function:: CreateUpdateConfigFileFromCLI

   func CreateUpdateConfigFileFromCLI(fileName string) (err error)

   /////////////////////////////////////// Create or update Config file from CLI /////////////////////////////////////// CreateUpdateConfigFileFromCLI reads and updates the configuration file based on command-line arguments.  Parameters: - fileName: The name of the configuration file.  Returns: - err: An error if there was an issue creating or updating the configuration file.


.. function:: InitGlobalConfigFromAzureKeyVault

   func InitGlobalConfigFromAzureKeyVault() (err error)

   /////////////////////////////////////////////// Extract Config variables from Azure Key Vault /////////////////////////////////////////////// InitGlobalConfigFromAzureKeyVault extracts the configuration from Azure Key Vault. It iterates over all secrets in the key vault and if the secret name matches a field in the Config struct, it sets the field to the value of the secret.  Returns: - err: An error if there was an issue extracting the configuration.


.. function:: ValidateConfig

   func ValidateConfig(config Config, requiredProperties \[\]string) (err error)

   ///////////////////// Helper Functions ///////////////////// ValidateConfig checks for mandatory entries in the configuration and validates chosen models.  Parameters: - config: The configuration object to validate. - requiredProperties: The list of required properties.  Returns: - err: An error if there was an issue validating the configuration.


.. function:: GetGlobalConfigAsJSON

   func GetGlobalConfigAsJSON() string

   GetGlobalConfigAsJSON returns the global configuration as a JSON string.  Returns: - string: The global configuration as a JSON string.


.. function:: HandleLegacyPortDefinition

   func HandleLegacyPortDefinition(configAddress string, legacyPort string) (webserverAddress string, err error)

   ////////////////////////// Legacy Config Converters ////////////////////////// HandleLegacyPortDefinition checks if the address is set, and if not, uses the legacy port to define the web server address. If both are empty, it returns an error.  Parameters: - address: The address to use for the web server. - legacyPort: The legacy port to use if the address is not set.  Returns: - webserverAddress: The web server address to use. - err: An error if both address and legacy port are empty.


.. function:: TestDefineOptionalProperties

   func TestDefineOptionalProperties(t \*testing.T)

   TestDefineOptionalProperties tests the defineOptionalProperties function


Types
-----

.. type:: Config

   type Config struct

   Config contains all the configuration settings for the Aali service.


.. type:: FlowkitConnection

   type FlowkitConnection struct

   FlowkitConnections contains the configuration for connecting to the FlowKit server.


Variables
---------

.. data:: GlobalConfig

   var GlobalConfig \*Config

   Initialize conifg dict
