package config

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
	"gopkg.in/yaml.v2"
)

////////////////////////////////////////////
// Standard Config init for Allie Go Modules
////////////////////////////////////////////

// InitConfig initializes the configuration for the Allie service.
//
// Parameters:
//   - requiredProperties: The list of required properties.
//   - optionalDefaultValues: The map of optional properties and their default values.
func InitConfig(requiredProperties []string, optionalDefaultValues map[string]interface{}) {
	// Get config file location
	// 1st option: read from environment variable
	configFile := os.Getenv("ALLIE_CONFIG_PATH")
	if configFile == "" {
		// 2nd option: read from default location... root directory
		configFile = "config.yaml"
	}

	// Get config properties from CLI
	err := CreateUpdateConfigFileFromCLI(configFile)
	if err != nil {
		pan := writeStringToFile("error in creating and/or updating configuration file from command line:")
		if pan != nil {
			panic(pan)
		}
		pan = writeInterfaceToFile(err)
		if pan != nil {
			panic(pan)
		}
		panic(err)
	}

	// Initialize config From File
	err = InitGlobalConfigFromFile(configFile, requiredProperties, optionalDefaultValues)
	if err != nil {
		pan := writeStringToFile("error in reading configuration values from configuration file:")
		if pan != nil {
			panic(pan)
		}
		pan = writeInterfaceToFile(err)
		if pan != nil {
			panic(pan)
		}
		panic(err)
	}

	// Optionally retrieve secrets from Azure Key Vault with Managed Identity (only works inside Azure Services)
	if GlobalConfig.EXTRACT_CONFIG_FROM_AZURE_KEY_VAULT {
		// Validate the required properties for Azure Key Vault are set
		err = ValidateConfig(*GlobalConfig, []string{"AZURE_KEY_VAULT_NAME", "AZURE_MANAGED_IDENTITY_ID"})
		if err != nil {
			pan := writeStringToFile("error in validating the mandatory configuration values for extracting configuration from Azure Key Vault:")
			if pan != nil {
				panic(pan)
			}
			pan = writeInterfaceToFile(err)
			if pan != nil {
				panic(pan)
			}
			panic(err)
		}

		// Initialize the config from Azure Key Vault
		err = InitGlobalConfigFromAzureKeyVault()
		if err != nil {
			pan := writeStringToFile("error in retrieving configuration values from Azure Key Vault:")
			if pan != nil {
				panic(pan)
			}
			pan = writeInterfaceToFile(err)
			if pan != nil {
				panic(pan)
			}
			panic(err)
		}
	}

	// Validate mandatory config properties
	err = ValidateConfig(*GlobalConfig, requiredProperties)
	if err != nil {
		pan := writeStringToFile("error in validating configuration variables:")
		if pan != nil {
			panic(pan)
		}
		pan = writeInterfaceToFile(err)
		if pan != nil {
			panic(pan)
		}
		panic(err)
	}
}

//////////////////////////////////////////
// Read Config variables from Config file
//////////////////////////////////////////

// InitGlobalConfigFromFile reads the configuration file and initializes the Config object.
//
// Parameters:
//   - fileName: The name of the configuration file.
//   - requiredProperties: The list of required properties.
//   - optionalDefaultValues: The map of optional properties and their default values.
//
// Returns:
//   - err: An error if there was an issue initializing the configuration.
func InitGlobalConfigFromFile(fileName string, requiredProperties []string, optionalDefaultValues map[string]interface{}) (err error) {
	var config Config
	configResult, err := readYaml(fileName, config)
	if err != nil {
		return err
	}

	// Assign to global config
	GlobalConfig = &configResult

	// Set optional properties if missing
	err = defineOptionalProperties(GlobalConfig, optionalDefaultValues)
	if err != nil {
		return err
	}

	return nil
}

// readYaml reads the yaml specified in `fileName` parameter and saves it to `config_struct`
//
// Parameters:
//   - fileName: The name of the configuration file.
//   - configStruct: Struct with the parameters of the YAML to read.
//
// Returns:
//   - extractedConfigStruct: The extracted configuration struct.
//   - err: An error if there was an issue reading the YAML file.
func readYaml(fileName string, configStruct Config) (extractedConfigStruct Config, err error) {
	// Read the YAML file into a byte slice
	data, err := os.ReadFile(fileName)
	if err != nil {
		message := "config.yaml file is missing from directory or not accessible"
		return Config{}, fmt.Errorf(message)
	}

	// Unmarshal the YAML data into the config object
	err = yaml.Unmarshal(data, &configStruct)
	if err != nil {
		// Create a new Config struct with field names and types
		configStruct := reflect.TypeOf(configStruct)
		var fieldList []string
		for i := 0; i < configStruct.NumField(); i++ {
			field := configStruct.Field(i)
			fieldList = append(fieldList, fmt.Sprintf("%q: %s", field.Name, field.Type.String()))
		}

		// Define error message
		message := fileName + " contains incorrect content. The allowed values are as follows: {"
		message += strings.Join(fieldList, ",")
		message = message + "}"

		// Write message and error to error file
		return Config{}, fmt.Errorf(message)
	}
	return configStruct, nil
}

// defineOptionalProperties sets optional properties for the configuration.
//
// Parameters:
//   - config: The configuration object to validate.
//   - optionalDefaultValues: The map of optional properties and their default values.
//
// Returns:
//   - err: An error if there was an issue setting the optional properties.
func defineOptionalProperties(config *Config, optionalDefaultValues map[string]interface{}) (err error) {
	defer func() {
		r := recover()
		if r != nil {
			// Write message to error file
			message := fmt.Sprintf("Error in defineOptionalProperties: %v", r)
			err = fmt.Errorf(message)
		}
	}()

	// Iterate over the optional default values
	for key, defaultValue := range optionalDefaultValues {
		// Use reflection to check if the field exists and is set to its zero value
		fieldValue := reflect.ValueOf(config).Elem().FieldByName(key)
		if fieldValue.IsValid() && isZeroValue(fieldValue) {
			// Set the default value using reflection
			if reflect.TypeOf(defaultValue) == fieldValue.Type() {
				fieldValue.Set(reflect.ValueOf(defaultValue))
			} else {
				// Handle type mismatch
				message := fmt.Sprintf("Type mismatch for key '%s': expected %v, got %v", key, fieldValue.Type(), reflect.TypeOf(defaultValue))
				return fmt.Errorf(message)
			}
		}
	}

	return nil
}

// isZeroValue checks if a reflect.Value is zero for its type.
//
// Parameters:
//   - v: The reflect.Value to check.
//
// Returns:
//   - bool: True if the value is zero, false otherwise.
func isZeroValue(v reflect.Value) bool {
	return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
}

/////////////////////////////////////////
// Create or update Config file from CLI
/////////////////////////////////////////

// CreateUpdateConfigFileFromCLI reads and updates the configuration file based on command-line arguments.
//
// Parameters:
//   - fileName: The name of the configuration file.
//
// Returns:
//   - err: An error if there was an issue creating or updating the configuration file.
func CreateUpdateConfigFileFromCLI(fileName string) (err error) {
	// Checking for any command-line arguments
	if len(os.Args) == 1 {
		log.Println("No command line options given; full config will be retrieved from existing config.yaml file and/or Azure Key Vault.")
		return
	}

	// Create a new config to store command-line options
	cliConfig := Config{}

	// Use reflection to create flags for each field in Config
	createFlags(reflect.ValueOf(&cliConfig).Elem(), "")

	// Parse the flags
	flag.Parse()

	// Checking if config.yaml file already exists
	_, err = os.Stat(fileName)

	if os.IsNotExist(err) {
		// If it doesn't exist, create a new one
		log.Println("config.yaml file does not exist. Creating a new one with provided CLI values...")

		file, _ := yaml.Marshal(cliConfig)
		_ = os.WriteFile(fileName, file, 0644)
	} else {
		// If it does exist, open and append to it
		log.Println("config.yaml file exists. Appending command line options with provided CLI values...")

		file, _ := os.ReadFile(fileName)
		config := Config{}
		err := yaml.Unmarshal(file, &config)
		if err != nil {
			message := fmt.Sprintf("Error in yaml.Unmarshal: %v", err)
			return fmt.Errorf(message)
		}

		// Use reflection to update fields that were set in the command line
		valCli := reflect.ValueOf(&cliConfig).Elem()
		valConfig := reflect.ValueOf(&config).Elem()
		t := reflect.TypeOf(cliConfig)

		for i := 0; i < t.NumField(); i++ {
			cliField := valCli.Field(i)
			configField := valConfig.Field(i)
			if !reflect.DeepEqual(cliField.Interface(), reflect.Zero(cliField.Type()).Interface()) {
				configField.Set(cliField)
			}
		}

		// Write back to the file
		file, _ = yaml.Marshal(config)
		_ = os.WriteFile(fileName, file, 0644)
	}

	return nil
}

// CreateFlags initializes command-line flags for configuration.
//
// Parameters:
//   - val: The value to create flags for.
//   - prefix: The prefix to use for the flags.
func createFlags(val reflect.Value, prefix string) {
	t := val.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		name := strings.ToUpper(prefix + field.Name)
		switch field.Type.Kind() {
		case reflect.String:
			flag.StringVar(val.Field(i).Addr().Interface().(*string), name, "", "config option")
		case reflect.Int:
			flag.IntVar(val.Field(i).Addr().Interface().(*int), name, 0, "config option")
		case reflect.Bool:
			flag.BoolVar(val.Field(i).Addr().Interface().(*bool), name, false, "config option")
		case reflect.Slice:
			if field.Type.Elem().Kind() == reflect.String {
				flag.Var((*flagStringSlice)(val.Field(i).Addr().Interface().(*[]string)), name, "config option")
			}
		case reflect.Struct:
			createFlags(val.Field(i), name+"_")
		}
	}
}

/////////////////////////////////////////////////
// Extract Config variables from Azure Key Vault
/////////////////////////////////////////////////

// InitGlobalConfigFromAzureKeyVault extracts the configuration from Azure Key Vault.
// It iterates over all secrets in the key vault and if the secret name matches a field in the Config struct,
// it sets the field to the value of the secret.
//
// Returns:
//   - err: An error if there was an issue extracting the configuration.
func InitGlobalConfigFromAzureKeyVault() (err error) {
	// log
	log.Println("Extracting configuration from Azure Key Vault...")

	// get environment variables
	azureManagedIdentity := os.Getenv(GlobalConfig.AZURE_MANAGED_IDENTITY_ID)
	azureKeyVaultName := os.Getenv(GlobalConfig.AZURE_KEY_VAULT_NAME)

	// check if all required environment variables are set
	if azureManagedIdentity == "" {
		return fmt.Errorf("environment variable for %v is not set", azureManagedIdentity)
	}
	if azureKeyVaultName == "" {
		return fmt.Errorf("environment variable for %v is not set", azureKeyVaultName)
	}

	// create key vault URL
	keyVaultUrl := fmt.Sprintf("https://%s.vault.azure.net/", azureKeyVaultName)

	// create Managed Identity credential
	cred, err := azidentity.NewManagedIdentityCredential(&azidentity.ManagedIdentityCredentialOptions{
		ID: azidentity.ClientID(azureManagedIdentity),
	})
	if err != nil {
		return fmt.Errorf("failed to get Managed Identity credential: %w", err)
	}

	// Test the managed id by getting a token
	scope := "https://vault.azure.net/.default" // Scope for Azure Key Vault
	_, err = cred.GetToken(context.TODO(), policy.TokenRequestOptions{
		Scopes: []string{scope},
	})
	if err != nil {
		return fmt.Errorf("failed to get token from managed ID: %w", err)
	}

	// Reflect on the struct
	GlobalConfigValue := reflect.ValueOf(GlobalConfig).Elem()
	GlobalConfigType := GlobalConfigValue.Type()

	// create azsecrets client
	clientSecrets, err := azsecrets.NewClient(keyVaultUrl, cred, nil)
	if err != nil {
		return err
	}

	// list all secrets
	pagerSecerts := clientSecrets.NewListSecretPropertiesPager(nil)
	// iterate over all secrets
	for pagerSecerts.More() {
		page, err := pagerSecerts.NextPage(context.TODO())
		if err != nil {
			return err
		}
		for _, secret := range page.Value {
			// iterate over all fields in the struct
			for i := 0; i < GlobalConfigValue.NumField(); i++ {
				// Get the field
				field := GlobalConfigValue.Field(i)

				// Get the YAML tag
				fieldType := GlobalConfigType.Field(i)
				yamlTag := fieldType.Tag.Get("json")

				// Check if the field name matches the target field name
				if yamlTag == secret.ID.Name() {
					// Get the key value
					resp, err := clientSecrets.GetSecret(context.TODO(), secret.ID.Name(), "", nil)
					if err != nil {
						return err
					}
					// Check if the field is settable
					if field.CanSet() {
						// Set the field to the new value
						switch field.Kind() {
						case reflect.String:
							field.SetString(*resp.Value)
						case reflect.Bool:
							// Convert string to bool
							value, err := strconv.ParseBool(*resp.Value)
							if err != nil {
								return err
							}
							field.SetBool(value)
						case reflect.Int:
							// Convert string to int
							value, err := strconv.Atoi(*resp.Value)
							if err != nil {
								return err
							}
							field.SetInt(int64(value))
						case reflect.Slice:
							// Convert string to string slice
							var value []string
							err := json.Unmarshal([]byte(*resp.Value), &value)
							if err != nil {
								return err
							}
							field.Set(reflect.ValueOf(value))
						default:
							return fmt.Errorf("unsupported field type: %v", field.Kind())
						}
					}
				}
			}
		}
	}

	return nil
}

///////////////////////
// Helper Functions
///////////////////////

// ValidateConfig checks for mandatory entries in the configuration and validates chosen models.
//
// Parameters:
//   - config: The configuration object to validate.
//   - requiredProperties: The list of required properties.
//
// Returns:
//   - err: An error if there was an issue validating the configuration.
func ValidateConfig(config Config, requiredProperties []string) (err error) {
	// Check if all mandatory properties are present
	configValue := reflect.ValueOf(config)

	for _, property := range requiredProperties {
		field := configValue.FieldByName(property)

		if !field.IsValid() || field.IsZero() {
			message := fmt.Sprintf("config.yaml is missing mandatory property '%v': ", property)
			return fmt.Errorf(message)
		}
	}

	return nil
}

// GetGlobalConfigAsJSON returns the global configuration as a JSON string.
//
// Returns:
//   - string: The global configuration as a JSON string.
func GetGlobalConfigAsJSON() string {
	jsonData, err := json.Marshal(GlobalConfig)
	if err != nil {
		return ""
	}
	return string(jsonData)
}

///////////////////////
// Error file creator
///////////////////////

// writeInterfaceToFile writes interface data to an error log file.
//
// Parameters:
//   - data: The data to write to the file.
//
// Returns:
//   - error: An error if there was an issue writing to the file.
func writeInterfaceToFile(data interface{}) error {
	var file *os.File
	var err error

	// Get file name
	filename := "error.log"

	// Create file
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// If the file does not exist, create a new file.
		file, err = os.Create(filename)
		if err != nil {
			return err
		}
	} else {
		// If the file already exists, open it in append mode.
		file, err = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
	}
	defer file.Close()

	// Write to file
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Add time
	timestamp := timeToString(time.Now())

	// Write to file
	line := fmt.Sprintf("%s: %s\n", timestamp, string(jsonData))
	_, err = file.Write([]byte(line))
	if err != nil {
		return err
	}

	return nil
}

// writeStringToFile writes string data to an error log file.
//
// Parameters:
//   - data: The data to write to the file.
//
// Returns:
//   - error: An error if there was an issue writing to the file.
func writeStringToFile(data string) error {
	var file *os.File
	var err error

	// Get file name
	filename := "error.log"

	// Add time
	timestamp := timeToString(time.Now())

	// Change string
	data = timestamp + ": " + data

	if _, err = os.Stat(filename); os.IsNotExist(err) {
		// File does not exist, create a new file
		file, err = os.Create(filename)
		if err != nil {
			return err
		}
	} else {
		// File exists, open it for appending
		file, err = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
	}
	defer file.Close()

	// Append data to file with a new line
	_, err = fmt.Fprintln(file, data)
	return err
}

// timeToString converts a time.Time value to a formatted string.
//
// Parameters:
//   - t: The time.Time value to convert.
//
// Returns:
//   - string: The formatted string.
func timeToString(t time.Time) string {
	layout := "2006-01-02 15:04:05.000"
	return t.Format(layout)
}
