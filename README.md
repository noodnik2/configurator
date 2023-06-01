# configurator - `golang` CLI Configuration Library

Target functionality:
- Make it easy to load and save application configuration in a file (e.g., `$HOME/.config/`_<your-app-name>_)
- Allow configuration value overrides via environment variables
- Enable a user to create or modify the configuration without a code editor

## Statistics
- [![Go Coverage](https://github.com/noodnik2/configurator/wiki/coverage.svg)](https://raw.githack.com/wiki/noodnik2/configurator/coverage.html)

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

### Examples

See example code that uses `configurator` in the [examples](./examples) sub-folder.

## Credits

`configurator` builds on top of and extends three related foundational libraries:
- [GoDotEnv] - sets environment variables from a configuration file
- [Envconfig] - populates struct field values based on environment variables
- [promptui] - orchestrates a console-based dialog to set / modify configuration values

## Caveats
- Many cool features of [Envconfig] (such as the use of [Prefixes](https://github.com/sethvargo/go-envconfig/tree/main#prefix),
  [Complex Types](https://github.com/sethvargo/go-envconfig/tree/main#complex-types) 
  and [Nested Structs](https://github.com/sethvargo/go-envconfig/tree/main#structs)) aren't supported as of this writing,
  but could be added if / when needed.

[GoDotEnv]: https://github.com/joho/godotenv
[Envconfig]: https://github.com/sethvargo/go-envconfig
[promptui]: https://github.com/manifoldco/promptui







