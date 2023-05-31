# configurator

Golang CLI configuration library targeting functionalities:
- Make it easy to load and save configuration values in a file
- Allow configuration value overrides via environment variables
- Enable the user to create or modify the configuration without a code editor

## Requirements

Since `configurator` uses [Generics](https://go.dev/doc/tutorial/generics), `go` version `1.18`
or greater is required.

## Usage

### Include `configurator` Into Your Project 
```shell
$ go get github.com/noodnik2/configurator
```

### Main APIs

- `LoadConfig[T any](configFile string, config *T) error` - loads configuration from a file
- `SaveConfig[T any](configFileName string, config T) error` - saves configuration to a file
- `EditConfig[T any](config *T) error` - invokes a user dialog to set or update the configuration

### Second-Level APIs
- `GetConfigEnvItems[T any](config T) ([]ConfigEnvItem, error)` - gets a list of configuration items
- `SetConfigEnvItem[T any](config *T, envName, newValueAsString string) error` - updates a single configuration item

See the source code for details.

## Credits

`configurator` builds on top of and extends three related foundational libraries:
- [godotenv](https://github.com/joho/godotenv) - sets environment variables from a configuration file
- [go-envconfig](https://github.com/sethvargo/go-envconfig) - populates struct field values based on environment variables
- [promptui](https://github.com/manifoldco/promptui) - orchestrates a console-based dialog to set / modify configuration values







