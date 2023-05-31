package main

import (
	"flag"
	"log"
	"strings"

	"github.com/noodnik2/configurator"
)

// configFilename - name of the file containing the configuration values; for standalone
// (e.g., CLI) applications, suggest $(HOME)/.config/<your-program-name>
const configFilename = "example.env"

// Config - Your configuration structure.  Properties must be public and have the `env`
// tags as documented in https://github.com/sethvargo/go-envconfig
type Config struct {
	FlitsPerGazeebop int     `env:"FPG"`
	ConversionRate   float32 `env:"CONVERSION_RATE,default=3.14"`
	IsRundable       bool    `env:"RUNDABLE"`
	AccessKey        string  `env:"ACCESS_KEY,required" secret:"hide"`
	Last4Ssn         string  `env:"LAST4_SSN,required" secret:"mask"`
	NotAnEnv         string
	notSeen          string
}

// main - Example of using 'configurator' to load configuration values
func main() {

	var editConfigurator = flag.Bool("editConfigurator", false, "invoke configurator editor")
	flag.Parse()

	log.Println("Current Configuration:")
	showConfig()

	if editConfigurator != nil && *editConfigurator {
		log.Println()
		log.Println("Invoking Configurator Editor:")
		editConfig()
		log.Println()
		log.Println("Updated Configuration:")
		showConfig()
	}
}

func showConfig() {

	config := getConfig()

	log.Println()
	log.Println("Structure:")
	log.Printf("\t%#v\n", config)

	configEnvItems, getConfigErr := configurator.GetConfigEnvItems(config)
	if getConfigErr != nil {
		log.Fatalf("couldn't GetConfigEnvItems(%s): %v\n", configFilename, getConfigErr)
	}

	log.Println()
	log.Println("Items:")
	for _, configEnvItem := range configEnvItems {
		var val any
		if configEnvItem.Secret == "" {
			val = configEnvItem.Val
		} else {
			// dealing with these flags is currently the client's responsibility; encapsulating
			// support for handling these within 'configurator' is under consideration.
			if configEnvItem.Secret == "mask" {
				val = strings.Repeat("*", len(configEnvItem.Val.(string)))
			} else {
				val = "<suppressed>"
			}
		}
		log.Printf("\t%s: %v\n", configEnvItem.Name, val)
	}

}

// editConfig invokes the configurator editor on the configuration returned by getConfig
func editConfig() {
	config := getConfig()
	if editErr := configurator.EditConfig(&config); editErr != nil {
		log.Fatalf("couldn't EditConfig(): %v\n", editErr)
	}
	if saveErr := configurator.SaveConfig(configFilename, config); saveErr != nil {
		log.Fatalf("couldn't SaveConfig(%s): %v\n", configFilename, saveErr)
	}
}

// getConfig returns the configuration loaded from the config file or from the environment (overrides)
func getConfig() Config {
	var appConfig Config
	if loadConfigErr := configurator.LoadConfig(configFilename, &appConfig); loadConfigErr != nil {
		log.Fatalf("couldn't LoadConfig(%s): %v\n", configFilename, loadConfigErr)
	}
	return appConfig
}
