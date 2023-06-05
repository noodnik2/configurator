package configurator

import (
	"fmt"
	"io"
	"log"
	"os"
	"sort"
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
	return updateConfigFromMap(configFile, configMap)
}

// updateConfigFromMap updates both the written configuration and the environment
// to match the contents of the supplied map containing all configuration entries.
// Configuration entries with nil values will be removed from both targets.  NOTE:
// no transactional guarantees are provided; if an error is returned, partial
// update(s) may have been made.
func updateConfigFromMap(truncatedConfigFile io.Writer, fullConfigMap map[string]any) error {
	sortedEnvVarNames := make([]string, len(fullConfigMap))
	for envVarName := range fullConfigMap {
		sortedEnvVarNames = append(sortedEnvVarNames, envVarName)
	}
	sort.Strings(sortedEnvVarNames)

	// write the new configuration entries
	for _, envVarName := range sortedEnvVarNames {
		envVal := fullConfigMap[envVarName]
		if envVal == nil {
			continue
		}
		if _, printErr := fmt.Fprintf(truncatedConfigFile, "%s=%v\n", envVarName, envVal); printErr != nil {
			return printErr
		}
	}
	// update the environment
	cantUpdateVars := make(map[string][]string)
	for _, envVarName := range sortedEnvVarNames {
		envVal := fullConfigMap[envVarName]
		if envVal == nil {
			if unSetEnvErr := os.Unsetenv(envVarName); unSetEnvErr != nil {
				cantUpdateVars[unSetEnvErr.Error()] = append(cantUpdateVars[unSetEnvErr.Error()], envVarName)
			}
			continue
		}
		if setEnvErr := os.Setenv(envVarName, fmt.Sprintf("%v", envVal)); setEnvErr != nil {
			cantUpdateVars[setEnvErr.Error()] = append(cantUpdateVars[setEnvErr.Error()], envVarName)
		}
	}
	if len(cantUpdateVars) != 0 {
		return fmt.Errorf("couldn't update environment variable(s): %s", cantUpdateVars)
	}
	return nil
}
