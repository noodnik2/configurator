package configurator

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestApiSave(t *testing.T) {

	const envFileVS1Value = "this is the value of V_S1"

	type testConfig struct {
		S1 string `env:"V_S1" desc:"First string"`
		S2 string `env:"V_S2,default=Maybe"`
		S3 string `env:"V_S3" secret:"true"`
		S4 string `env:"V_S4,default=shush" secret:"mask"`
		B1 bool   `env:"YESNO,default=true"`
	}

	requirer := require.New(t)

	envFileName, ctefErr := createTempEnvFileFromMap(t, map[string]any{"V_S1": envFileVS1Value})
	requirer.NoError(ctefErr)

	// load & verify config value read from file system
	config1 := testConfig{}
	requirer.NoError(LoadConfig(envFileName, &config1))
	requirer.Equal(envFileVS1Value, config1.S1)

	// update the config and save it back to the filesystem
	newS2 := config1.S1 + " in an update"
	config1.S2 = newS2
	requirer.NoError(SaveConfig(envFileName, config1))

	// verify expected changes in the filesystem are seen when loading new config
	newConfig := testConfig{}
	requirer.NoError(LoadConfig(envFileName, &newConfig))
	requirer.Equal(envFileVS1Value, newConfig.S1)
	requirer.Equal(newS2, newConfig.S2)
}

func TestUpdateConfigFromMap(t *testing.T) {

	testCases := []struct {
		name           string
		configVars     map[string]any
		envVarsSet     []string
		expectedOutput string
	}{
		{
			name:           "basic",
			configVars:     map[string]any{"V1": "val1", "V2": "val2", "V3": ""},
			expectedOutput: "V1=val1\nV2=val2\nV3=\n",
		},
		{
			name:           "unset",
			configVars:     map[string]any{"V1": "val1", "V2": nil, "V3": ""},
			expectedOutput: "V1=val1\nV3=\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			requirer := require.New(t)

			// set the variables to something OTHER than what they will be set to
			for k := range tc.configVars {
				t.Setenv(k, fmt.Sprintf("some different value for '%s'", k))
			}

			// invoke the writer
			writer := &bytes.Buffer{}
			requirer.NoError(updateConfigFromMap(writer, tc.configVars))
			requirer.Equal(tc.expectedOutput, writer.String())

			// ensure the environment now has the correct values
			for k, v := range tc.configVars {
				if v != nil {
					requirer.Equal(v.(string), os.Getenv(k))
				} else {
					_, isFound := os.LookupEnv(k)
					requirer.False(isFound)
				}
			}
		})
	}

}
