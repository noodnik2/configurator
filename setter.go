package configurator

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// SetConfigEnvItem allows setting in-place config values by the Name of their corresponding environment variable.
// See https://go.dev/blog/laws-of-reflection and https://research.swtch.com/interfaces
func SetConfigEnvItem[T any](config *T, envName, newValueAsString string) error {
	cfgStructType, cfgStructElements, getConfigInfoErr := getConfigStructInfo(config)
	if getConfigInfoErr != nil {
		return getConfigInfoErr
	}

	var isSet bool
	for fieldIndex := 0; fieldIndex < cfgStructType.NumField(); fieldIndex++ {

		cfgStructFieldTag := cfgStructType.Field(fieldIndex).Tag

		var ok bool
		var cfgStructFieldEnvTagValue string
		if cfgStructFieldEnvTagValue, ok = cfgStructFieldTag.Lookup(envTagKey); !ok {
			continue
		}

		tagParts := strings.Split(cfgStructFieldEnvTagValue, ",")
		if len(tagParts) == 0 || envName != tagParts[0] {
			continue
		}

		cfgStructFieldElement := cfgStructElements.Field(fieldIndex)
		if !cfgStructFieldElement.CanSet() {
			return fmt.Errorf("can't set(%s); not settable", envName)
		}

		cfgStructFieldElementKind := cfgStructFieldElement.Kind()
		switch cfgStructFieldElementKind {
		case reflect.String:
			cfgStructFieldElement.SetString(newValueAsString)
			isSet = true
		case reflect.Bool:
			parseBool, parseBoolErr := strconv.ParseBool(newValueAsString)
			if parseBoolErr != nil {
				return parseBoolErr
			}
			cfgStructFieldElement.SetBool(parseBool)
			isSet = true
		case reflect.Float64, reflect.Float32:
			parseFloat, parseFloatErr := strconv.ParseFloat(newValueAsString,
				map[reflect.Kind]int{reflect.Float64: 64, reflect.Float32: 32}[cfgStructFieldElementKind])
			if parseFloatErr != nil {
				return parseFloatErr
			}
			cfgStructFieldElement.SetFloat(parseFloat)
			isSet = true
		case reflect.Int, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8:
			parseInt, parseIntErr := strconv.ParseInt(newValueAsString, 10,
				map[reflect.Kind]int{reflect.Int: strconv.IntSize, reflect.Int64: 64, reflect.Int32: 32, reflect.Int16: 16, reflect.Int8: 8}[cfgStructFieldElementKind])
			if parseIntErr != nil {
				return parseIntErr
			}
			cfgStructFieldElement.SetInt(parseInt)
			isSet = true
		case reflect.Uint, reflect.Uint64, reflect.Uint32, reflect.Uint16, reflect.Uint8:
			parseUint, parseUintErr := strconv.ParseUint(newValueAsString, 10,
				map[reflect.Kind]int{reflect.Uint: strconv.IntSize, reflect.Uint64: 64, reflect.Uint32: 32, reflect.Uint16: 16, reflect.Uint8: 8}[cfgStructFieldElementKind])
			if parseUintErr != nil {
				return parseUintErr
			}
			cfgStructFieldElement.SetUint(parseUint)
			isSet = true
		default:
			return fmt.Errorf("unrecognized Kind(%v)", cfgStructFieldElementKind)
		}
		break
	}
	if !isSet {
		return fmt.Errorf("env value(%s) wan't set; not found", envName)
	}
	return nil
}
