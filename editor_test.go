package configurator

import (
	"os"
	"testing"

	"github.com/manifoldco/promptui"
	"github.com/stretchr/testify/require"
)

// TestApiEditEnv tests the synchronization between the configuration file and the environment
// when an update made by the editor is saved using SaveConfig.  NOTE: this test case was used
// to demonstrate a problem noted during an integration test, but was fixed in the internal
// "writeConfigMap" function (tested in TestWriteConfigMap).  Therefore, due to its redundancy,
// this test should be removed.
func TestApiEditEnv(t *testing.T) {

	const envFileVS1Value = "initial value of V_S1"

	type testConfig struct {
		S1 string `env:"V_S1" desc:"First string"`
	}

	requirer := require.New(t)

	envFileName, cleanupFn, ctefErr := createTempEnvFileFromMap(t, map[string]any{"V_S1": envFileVS1Value})
	requirer.NoError(ctefErr)
	defer cleanupFn()

	// load & verify config value read from file system
	config1 := testConfig{}
	requirer.NoError(LoadConfig(envFileName, &config1))
	requirer.Equal(envFileVS1Value, config1.S1)

	seam := &promptUiTestSeam{
		pr: mockPr{
			mockedResponses: map[int]string{
				1: "y", // user says "y" the second time through the dialog
			},
		},
	}
	requirer.NoError(editConfig(&config1, seam, 2))

	// confirm expected user interface dialog
	requirer.Equal(2, len(seam.prompters))

	// save the changes
	requirer.NoError(SaveConfig(envFileName, config1))

	// load the updated values into a new config structure
	// NOTE: the updated values will likely be coming from the
	// updated environment rather than the updated config file
	config2 := testConfig{}
	requirer.NoError(LoadConfig(envFileName, &config2))

	// confirms we're able to see the edited value after re-loading
	requirer.Equal(config1.S1, config2.S1)

}

func TestApiEdit(t *testing.T) {

	const envFileVS1Value = "this is the value of V_S1"

	type testConfig struct {
		S1 string `env:"V_S1" desc:"First string"`
		S2 string `env:"V_S2,default=Maybe"`
		S3 string `env:"V_S3" secret:"true"`
		S4 string `env:"V_S4,default=shush" secret:"mask"`
		S5 string `env:"V_S5" secret:"hide"`
		B6 bool   `env:"YESNO,default=true"`
	}

	requirer := require.New(t)

	envFileName, cleanupFn, ctefErr := createTempEnvFileFromMap(t, map[string]any{"V_S1": envFileVS1Value})
	requirer.NoError(ctefErr)
	defer cleanupFn()

	// load & verify config value read from file system
	config1 := testConfig{}
	requirer.NoError(LoadConfig(envFileName, &config1))
	requirer.Equal(envFileVS1Value, config1.S1)

	seam := &promptUiTestSeam{
		pr: mockPr{
			mockedResponses: map[int]string{
				11: "y", // user says "y" the second time through the dialog
			},
		},
	}
	requirer.NoError(editConfig(&config1, seam, 3))

	// verify expected user interface dialog
	requirer.Equal(14, len(seam.prompters))
	prompt1 := seam.prompters[0].(*promptui.Prompt)
	requirer.Equal("V_S1", prompt1.Label)
	prompt2 := seam.prompters[1].(*promptui.Prompt)
	requirer.Equal(int32(0), prompt2.Mask)
	requirer.Equal("V_S2", prompt2.Label)
	prompt4 := seam.prompters[3].(*promptui.Prompt)
	requirer.False(prompt4.AllowEdit)
	requirer.True(prompt4.HideEntered)
	requirer.Equal('*', prompt4.Mask)
	prompt5 := seam.prompters[4].(*promptui.Prompt)
	requirer.False(prompt5.AllowEdit)
	requirer.True(prompt5.HideEntered)
	requirer.Zero(prompt5.Mask)
	prompt6 := seam.prompters[5].(*promptui.Select)
	requirer.Equal("YESNO", prompt6.Label)
	requirer.Equal(2, len(prompt6.Items.([]string)))
	requirer.Equal(1, prompt6.CursorPos) // value defaults to "true", which is at offset 1
	prompt7 := seam.prompters[6].(*promptui.Prompt)
	requirer.Equal("Done", prompt7.Label)
	requirer.Equal("n", prompt7.Default)

}

type mockPr struct {
	responseCount   int
	mockedResponses map[int]string
}

func (m *mockPr) Run() (string, error) {
	if m.mockedResponses != nil {
		defer func() { m.responseCount++ }()
		if mockResponse, ok := m.mockedResponses[m.responseCount]; ok {
			return mockResponse, nil
		}
	}
	return "mock prompt response", nil
}

type mockSr struct {
	responseCount   int
	mockedResponses map[int]string
}

func (m *mockSr) Run() (int, string, error) {
	if m.mockedResponses != nil {
		defer func() { m.responseCount++ }()
		if mockResponse, ok := m.mockedResponses[m.responseCount]; ok {
			return 0, mockResponse, nil
		}
	}
	return 0, "mock select response", nil
}

type promptUiTestSeam struct {
	prompters []any
	pr        mockPr
	sr        mockSr
}

func (pu *promptUiTestSeam) getPrompter(pr promptRunner) promptRunner {
	pu.prompters = append(pu.prompters, pr)
	return &pu.pr
}

func (pu *promptUiTestSeam) getSelector(sr selectRunner) selectRunner {
	pu.prompters = append(pu.prompters, sr)
	return &pu.sr
}

func createTempEnvFileFromMap(t *testing.T, envMap map[string]any) (envFileName string, cleanupFn func(), err error) {
	var envFile *os.File
	if envFile, err = os.CreateTemp(t.TempDir(), t.Name()); err != nil {
		return
	}
	defer func() {
		require.NoError(t, envFile.Close())
	}()

	cleanupFn = func() {
		require.NoError(t, os.Remove(envFile.Name()))
	}
	if err = writeConfigMap(envFile, envMap); err != nil {
		cleanupFn()
		return
	}

	envFileName = envFile.Name()
	return
}
