package configurator

import (
	"fmt"
	"io"
	"log"
	"os"
)

// SaveConfig saves the current 'config' values into 'configFile', and
// updates the values of the corresponding environment variables.
func SaveConfig[T any](configFileName string, config T) error {
	envItems, getterErr := GetConfigEnvItems(config)
	if getterErr != nil {
		return getterErr
	}
	configMap := make(map[string]any, len(envItems))
	for _, envItem := range envItems {
		configMap[envItem.Name] = envItem.Val
	}
	return SaveConfigMap(configFileName, configMap)
}

// SaveConfigMap saves the map of environment name: environment value entries into 'configFile'.
func SaveConfigMap(configFileName string, configMap map[string]any) error {
	configFile, openErr := os.OpenFile(configFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if openErr != nil {
		return openErr
	}
	defer func() {
		if closeErr := configFile.Close(); closeErr != nil {
			log.Printf("NOTE: error closing %s: %v\n", configFileName, closeErr)
		}
	}()
	return writeConfigMap(configFile, configMap)
}

func writeConfigMap(configFile io.Writer, configMap map[string]any) error {
	// update the configuration file with the new value(s)
	for envVarName, envVarVal := range configMap {
		if _, printErr := fmt.Fprintf(configFile, "%s=%v\n", envVarName, envVarVal); printErr != nil {
			return printErr
		}
	}
	// update the environment with the new value(s)
	couldntSetVars := make(map[string][]string)
	for envVarName, envVarVal := range configMap {
		if _, envVarFound := os.LookupEnv(envVarName); envVarFound {
			if setEnvErr := os.Setenv(envVarName, fmt.Sprintf("%v", envVarVal)); setEnvErr != nil {
				couldntSetVars[setEnvErr.Error()] = append(couldntSetVars[setEnvErr.Error()], envVarName)
			}
		}
	}
	if len(couldntSetVars) != 0 {
		return fmt.Errorf("couldn't set environment variable(s): %s\n", couldntSetVars)
	}
	return nil
}
