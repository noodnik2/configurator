package configurator

import (
	"bytes"
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

	envFileName, cleanupFn, ctefErr := createTempEnvFileFromMap(t, map[string]any{"V_S1": envFileVS1Value})
	requirer.NoError(ctefErr)
	defer cleanupFn()

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

func TestWriteConfigMap(t *testing.T) {

	testCases := []struct {
		name           string
		configVars     map[string]any
		envVarsSet     []string
		envVarsNotSet  []string
		expectedOutput string
	}{
		{
			name:           "basic",
			configVars:     map[string]any{"V1": "val1", "V2": "val2"},
			envVarsSet:     []string{"V1"},
			envVarsNotSet:  []string{"V2"},
			expectedOutput: "V1=val1\nV2=val2\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			requirer := require.New(t)

			// clean up test case environment variables
			defer func() {
				for _, v := range append(tc.envVarsSet, tc.envVarsNotSet...) {
					_ = os.Unsetenv(v)
				}
			}()

			// set the "set" variables to something OTHER than what they will be set to
			for _, v := range tc.envVarsSet {
				requirer.NoError(os.Setenv(v, tc.configVars[v].(string)+"some changed value"))
			}
			// ensure the "not set" variables aren't set
			for _, v := range tc.envVarsNotSet {
				requirer.NoError(os.Unsetenv(v))
			}

			// invoke the writer
			writer := &bytes.Buffer{}
			requirer.NoError(writeConfigMap(writer, tc.configVars))
			requirer.Equal(tc.expectedOutput, writer.String())

			// ensure the "set" variables now have the correct value
			for _, v := range tc.envVarsSet {
				requirer.Equal(tc.configVars[v].(string), os.Getenv(v))
			}
			// ensure the "not set" variables are still not set
			for _, v := range tc.envVarsNotSet {
				_, envVarIsSet := os.LookupEnv(v)
				requirer.False(envVarIsSet)
			}
		})
	}

}
