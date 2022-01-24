package rules

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

const (
	requiredBackendType        = "azurerm"
	backendTypeMessageTemplate = "backend type should be \"%s\" but defined: \"%s\""
)

func NewBackendTypeRule() *Rule {
	return NewRule(
		"backend_type",
		func(runner tflint.Runner, rule tflint.Rule) error {
			backend, err := runner.Backend()
			if err != nil {
				return err
			}
			if backend == nil {
				return nil
			}

			if backend.Type != requiredBackendType {
				return runner.EmitIssue(
					rule,
					fmt.Sprintf(
						backendTypeMessageTemplate,
						requiredBackendType,
						backend.Type,
					),
					backend.DeclRange,
				)
			}

			return nil
		},
	)
}
