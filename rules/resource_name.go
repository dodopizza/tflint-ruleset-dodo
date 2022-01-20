package rules

import (
	"fmt"
	"regexp"

	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

const (
	resourceNameMessageTemplate = "name should include only lowercase letters or '_' symbol, address: \"%s\""
)

func NewResourceNameRule() *BaseRule {
	allowedSymbolsRegex := regexp.MustCompile("^[a-z0-9_]*$")

	return NewRule(
		"resource_name_lowercased",
		func(runner tflint.Runner, rule tflint.Rule) error {
			cfg, _ := runner.Config()
			for address, res := range cfg.Module.ManagedResources {
				if !allowedSymbolsRegex.MatchString(res.Name) {
					if err := runner.EmitIssue(
						rule,
						fmt.Sprintf(
							resourceNameMessageTemplate,
							address,
						),
						res.DeclRange,
					); err != nil {
						return err
					}
				}
			}

			return nil
		},
	)
}
