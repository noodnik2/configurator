package configurator

import (
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
