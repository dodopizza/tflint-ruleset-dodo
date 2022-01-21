package rules

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

const (
	backendTypeMessageTemplate = "backend type should be \"azurerm\" but defined: \"%s\""
)

func NewTerraformBackendTypeRule() *BaseRule {
	return NewRule(
		"terraform_backend_type",
		func(runner tflint.Runner, rule tflint.Rule) error {
			backend, err := runner.Backend()
			if err != nil {
				return err
			}
			if backend == nil {
				return nil
			}

			if backend.Type != "azurerm" {
				return runner.EmitIssue(
					rule,
					fmt.Sprintf(
						backendTypeMessageTemplate,
						backend.Type,
					),
					backend.DeclRange,
				)
			}

			return nil
		},
	)
}
