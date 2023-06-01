package configurator

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestApiLoadSuccess(t *testing.T) {

	const envFileVS1Value = "this is the value of V_S1"

	type testConfig struct {
		S1 string `env:"V_S1" desc:"First string"`
		S2 string `env:"V_S2,default=Maybe"`
		S3 string `env:"V_S3" secret:"true"`
		S4 string `env:"V_S4,default=shush" secret:"mask"`
		B5 bool   `env:"YESNO,default=true"`
	}

	requirer := require.New(t)

	envFileName, ctefErr := createTempEnvFileFromMap(t, map[string]any{"V_S1": envFileVS1Value})
	requirer.NoError(ctefErr)

	// confirm values are read from the filesystem
	config1 := testConfig{}
	requirer.NoError(LoadConfig(envFileName, &config1))
	requirer.Equal(envFileVS1Value, config1.S1)

	// confirm default values are applied
	requirer.True(config1.B5)
	requirer.Equal("Maybe", config1.S2)
	requirer.Equal("shush", config1.S4)

	// confirm initial values aren't overridden
	const initialS1Value = "some initial S1 value"
	config2 := testConfig{S1: initialS1Value}
	requirer.NoError(os.Setenv("V_S1", initialS1Value))
	requirer.NoError(LoadConfig(envFileName, &config2))
	requirer.Equal(initialS1Value, config2.S1)

	// confirm values in the environment take precedence over those read from the filesystem
	const newS1ValueInEnv = "new value for S1 in environment"
	config3 := testConfig{}
	defer func() {
		requirer.NoError(os.Unsetenv("V_S1"))
	}()
	requirer.NoError(os.Setenv("V_S1", newS1ValueInEnv))
	requirer.NoError(LoadConfig(envFileName, &config3))
	requirer.Equal(newS1ValueInEnv, config3.S1)
}

func TestApiLoadRequired(t *testing.T) {

	const envFileVS1Value = "this is the value of V_S1"

	type testConfig struct {
		S0  string `env:"V_S0,default=seen"`
		S1  string `env:"V_S1,required"`
		S1B string `env:"V_S1B,default=also_seen"`
		S2  string `env:"V_S2,required"`
		S3  string `env:"V_S3,default=not_seen"` // lack of required entry aborted the loop in envconfig.Process()
	}

	requirer := require.New(t)

	// only one of the two required values is provided
	envFileName, ctefErr := createTempEnvFileFromMap(t, map[string]any{"V_S1": envFileVS1Value})
	requirer.NoError(ctefErr)

	// confirm expected error is returned
	config := testConfig{}
	loadErr := LoadConfig(envFileName, &config)
	requirer.Error(loadErr)
	requirer.Contains(loadErr.Error(), "missing required value")
	requirer.Contains(loadErr.Error(), "V_S2")
	requirer.Equal("seen", config.S0)
	requirer.Equal("also_seen", config.S1B)
	requirer.Equal("", config.S3)
}
