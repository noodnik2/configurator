package configurator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetConfigEnvItems(t *testing.T) {

	testCases := []struct {
		name       string
		config     any
		assertions func(*require.Assertions, []ConfigEnvItem)
	}{
		{
			name:   "empty config structure",
			config: struct{}{},
			assertions: func(requirer *require.Assertions, items []ConfigEnvItem) {
				requirer.Nil(items)
			},
		},
		{
			name: "non-tagged, private, public and non-env elements",
			config: struct {
				f1 string
				f2 string `env:"privateEnv"`
				F3 string `json:"f3"`
				F4 string `env:"f4"`
			}{},
			assertions: func(requirer *require.Assertions, items []ConfigEnvItem) {
				requirer.Equal(1, len(items))
				requirer.Equal("f4", items[0].Name)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			requirer := require.New(t)
			items, getterErr := GetConfigEnvItems(tc.config)
			requirer.NoError(getterErr)
			tc.assertions(requirer, items)
		})
	}

}
