package rules

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

const rulePrefix = "dodo"

type BaseRule struct {
	name      string
	checkFunc func(tflint.Runner, tflint.Rule) error
}

var _ tflint.Rule = &BaseRule{}

func NewRule(
	name string,
	checkFunc func(tflint.Runner, tflint.Rule) error,
) *BaseRule {
	return &BaseRule{
		name:      fmt.Sprintf("%s_%s", rulePrefix, name),
		checkFunc: checkFunc,
	}
}

func (rule *BaseRule) Name() string {
	return rule.name
}

func (rule *BaseRule) Enabled() bool {
	return true
}

func (rule *BaseRule) Severity() string {
	return tflint.ERROR
}

func (rule *BaseRule) Link() string {
	return ""
}

func (rule *BaseRule) Check(runner tflint.Runner) error {
	config, err := runner.Config()
	if err != nil {
		return err
	}

	// Check if it is child module and do not evaluate them.
	if len(config.Path) != 0 {
		return nil
	}

	return rule.checkFunc(runner, rule)
}
