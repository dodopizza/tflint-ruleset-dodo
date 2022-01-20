package rules

import (
	"fmt"
	"strings"

	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type ResourceNameLowercasedRule struct{}

var _ tflint.Rule = &ResourceNameLowercasedRule{}

func NewResourceNameLowercasedRule() *ResourceNameLowercasedRule {
	return &ResourceNameLowercasedRule{}
}

func (r *ResourceNameLowercasedRule) Name() string {
	return "resource_name_lowercased"
}

func (r *ResourceNameLowercasedRule) Enabled() bool {
	return true
}

func (r *ResourceNameLowercasedRule) Severity() string {
	return tflint.ERROR
}

// Link returns the rule reference link.
func (r *ResourceNameLowercasedRule) Link() string {
	return ""
}

// Check checks whether ...
func (r *ResourceNameLowercasedRule) Check(runner tflint.Runner) error {
	cfg, _ := runner.Config()
	for _, res := range cfg.Module.ManagedResources {
		if res.Name != strings.ToLower(res.Name) {
			if err := runner.EmitIssue(
				r,
				fmt.Sprintf("name is not lowercased: %s", res.Name),
				res.DeclRange,
			); err != nil {
				return err
			}
		}
	}

	return nil
}
