package configurator

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSetConfigItem(t *testing.T) {

	type testConfig struct {
		S1   string  `env:"S1"`
		B2   bool    `env:"B2"`
		F64  float64 `env:"F64"`
		I32  int32   `env:"I32"`
		UI16 uint16  `env:"UI16"`
		SP6  *string `env:"SP6"`
	}

	const s1Value = "S1"
	const b2Value = true
	const sp6Value = "SP6"

	testCases := []struct {
		name       string
		config     testConfig
		settings   map[string]string
		assertions func(*require.Assertions, any, []error)
	}{
		{
			name:   "empty config structure",
			config: testConfig{},
			settings: map[string]string{
				"nonexistent": "doesn't matter",
			},
			assertions: func(requirer *require.Assertions, newConfig any, errors []error) {
				requirer.Equal(1, len(errors))
				requirer.Contains(errors[0].Error(), "not found")
				requirer.Equal(testConfig{}, newConfig)
			},
		},
		{
			name:   "all supported types",
			config: testConfig{},
			settings: map[string]string{
				"S1":   fmt.Sprintf("%v", s1Value),
				"B2":   fmt.Sprintf("%v", b2Value),
				"F64":  fmt.Sprintf("%v", math.MaxFloat64),
				"I32":  fmt.Sprintf("%v", math.MinInt32),
				"UI16": fmt.Sprintf("%v", math.MaxUint16),
				"SP6":  fmt.Sprintf("%v", sp6Value),
			},
			assertions: func(requirer *require.Assertions, newConfig any, errors []error) {

				expectedConfig := testConfig{
					S1:   s1Value,
					B2:   b2Value,
					F64:  math.MaxFloat64,
					I32:  math.MinInt32,
					UI16: math.MaxUint16,
				}

				requirer.Equal(1, len(errors), errors)
				requirer.Contains(errors[0].Error(), "unrecognized Kind")
				requirer.Equal(expectedConfig, newConfig)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			requirer := require.New(t)
			updatedConfig := tc.config
			var errors []error
			for envName, newVal := range tc.settings {
				if setErr := SetConfigEnvItem(&updatedConfig, envName, newVal); setErr != nil {
					errors = append(errors, setErr)
				}
			}
			tc.assertions(requirer, updatedConfig, errors)
		})
	}

}
