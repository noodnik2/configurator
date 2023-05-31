package configurator

import (
	"context"
	"log"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

// LoadConfig loads uninitialized configuration values from the environment or from
// 'configFile', applying the defaults as specified in the 'config' structure's tags.
// Upon successful return, all environment values on publicly accessible, supported
// properties of the 'config' structure are loaded both into the config structure
// and into the environment.
func LoadConfig[T any](configFile string, config *T) error {
	if err := godotenv.Load(configFile); err != nil {
		log.Printf("NOTE: ignored %v", err)
	}

	ctx := context.Background()
	if err := envconfig.Process(ctx, config); err != nil {
		return err
	}
	return nil
}
