package rules

import (
	"fmt"
	"strings"

	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

func NewResourceNameLowercasedRule() *BaseRule {
	return NewRule(
		"resource_name_lowercased",
		func(runner tflint.Runner, rule tflint.Rule) error {
			cfg, _ := runner.Config()
			for address, res := range cfg.Module.ManagedResources {
				if res.Name != strings.ToLower(res.Name) {
					if err := runner.EmitIssue(
						rule,
						fmt.Sprintf(
							"name \"%s\" is not lowercased, address: \"%s\"",
							res.Name,
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
