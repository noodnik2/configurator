package configurator

import (
	"fmt"
	"reflect"
	"strings"
)

// ConfigEnvItem contains properties related to a tagged environment item found within a passed structure
type ConfigEnvItem struct {
	Name   string
	Val    any
	Secret string
	Kind   reflect.Kind
}

const envTagKey = "env"

// GetConfigEnvItems gets a list of 'ConfigEnvItem' values from 'config'
// elements tagged as environment items.  See https://go.dev/blog/laws-of-reflection
func GetConfigEnvItems[T any](config T) ([]ConfigEnvItem, error) {
	cfgStructType, cfgStructElements, getConfigInfoErr := getConfigStructInfo(&config)
	if getConfigInfoErr != nil {
		return nil, getConfigInfoErr
	}

	var cfgTagItems []ConfigEnvItem
	for fieldIndex := 0; fieldIndex < cfgStructType.NumField(); fieldIndex++ {

		cfgStructFieldTag := cfgStructType.Field(fieldIndex).Tag

		var ok bool
		var cfgStructFieldEnvTagValue string
		if cfgStructFieldEnvTagValue, ok = cfgStructFieldTag.Lookup(envTagKey); !ok {
			continue
		}

		tagParts := strings.Split(cfgStructFieldEnvTagValue, ",")
		if len(tagParts) == 0 {
			continue
		}

		cfgStructFieldElement := cfgStructElements.Field(fieldIndex)
		if !cfgStructFieldElement.CanInterface() {
			// e.g., private visibility
			continue
		}

		envItem := ConfigEnvItem{Name: tagParts[0], Kind: cfgStructFieldElement.Kind()}
		if secretTagVal, okS := cfgStructFieldTag.Lookup("secret"); okS {
			envItem.Secret = secretTagVal
		}

		envItem.Val = reflect.ValueOf(cfgStructFieldElement.Interface()).Interface()
		cfgTagItems = append(cfgTagItems, envItem)
	}

	return cfgTagItems, nil
}

func getConfigStructInfo[T any](config *T) (reflect.Type, reflect.Value, error) {
	cfgStructType := reflect.TypeOf(*config)
	cfgStructElements := reflect.ValueOf(config).Elem()

	if cfgStructElements.Kind() == reflect.Interface {
		cfgStructElements = cfgStructElements.Elem()
	}
	if cfgStructElements.Kind() != reflect.Struct {
		return nil, reflect.Value{}, fmt.Errorf("unsupported config kind(%d)\n", cfgStructElements.Kind())
	}
	return cfgStructType, cfgStructElements, nil
}
