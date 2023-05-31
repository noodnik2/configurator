package configurator

import (
	"fmt"
	"log"
	"os"
)

// SaveConfig saves the current 'config' values into 'configFile'.
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
	configFile, openErr := os.OpenFile(configFileName, os.O_CREATE|os.O_WRONLY, 0644)
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

func writeConfigMap(configFile *os.File, configMap map[string]any) error {
	for k, v := range configMap {
		if _, printErr := fmt.Fprintf(configFile, "%s=%v\n", k, v); printErr != nil {
			return printErr
		}
	}
	return nil
}
