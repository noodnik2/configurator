package configurator

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/manifoldco/promptui"
)

// EditConfig invokes a user dialog to present and optionally
// change the current values in the 'config' structure
func EditConfig[T any](config *T) error {
	return editConfig(config, &promptUiSeamNoop{}, 100)
}

// editConfig provides a testable version of EditConfig
func editConfig[T any](config *T, seam promptUiSeam, maxTimes int) error {

	loopCounter := 0
	for {
		cfgTagItems, getterErr := GetConfigEnvItems(*config)
		if getterErr != nil {
			return getterErr
		}

		var promptErr error
		for _, cti := range cfgTagItems {
			var result string
			if cti.Kind == reflect.Bool {
				prompt := promptui.Select{
					Label:     cti.Name,
					Items:     []string{"False", "True"},
					CursorPos: map[bool]int{false: 0, true: 1}[cti.Val == true],
				}
				_, result, promptErr = seam.getSelector(&prompt).Run()
			} else {
				prompt := promptui.Prompt{
					Label:     cti.Name,
					Default:   fmt.Sprintf("%v", cti.Val),
					AllowEdit: true,
				}
				if cti.Secret != "" {
					prompt.HideEntered = true
					prompt.AllowEdit = false
					if cti.Secret == "mask" {
						prompt.Mask = '*'
					}
				}
				result, promptErr = seam.getPrompter(&prompt).Run()
			}
			if promptErr != nil {
				return promptErr
			}

			if setErr := SetConfigEnvItem(config, cti.Name, result); setErr != nil {
				log.Printf("NOTE: while setting item: %v\n", setErr)
			}
		}

		prompt := promptui.Prompt{
			Label:     "Done",
			Default:   "n",
			IsConfirm: true,
		}
		var isDone string
		isDone, promptErr = seam.getPrompter(&prompt).Run()

		if strings.ToLower(strings.Trim(isDone, " \t")) == "y" {
			break
		}

		loopCounter++
		if loopCounter >= maxTimes {
			return fmt.Errorf("too many edit attempts(%d)", loopCounter)
		}
	}

	return nil
}

type promptRunner interface {
	Run() (string, error)
}

type selectRunner interface {
	Run() (int, string, error)
}

type promptUiSeam interface {
	getPrompter(pr promptRunner) promptRunner
	getSelector(sr selectRunner) selectRunner
}

type promptUiSeamNoop struct{}

func (*promptUiSeamNoop) getPrompter(pr promptRunner) promptRunner {
	return pr
}

func (*promptUiSeamNoop) getSelector(sr selectRunner) selectRunner {
	return sr
}
