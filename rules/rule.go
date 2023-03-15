package rules

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

const rulePrefix = "dodo"

type Rule struct {
	tflint.Rule
	name      string
	checkFunc func(tflint.Runner, tflint.Rule) error
}

func NewRule(
	name string,
	checkFunc func(tflint.Runner, tflint.Rule) error,
) *Rule {
	return &Rule{
		name:      fmt.Sprintf("%s_%s", rulePrefix, name),
		checkFunc: checkFunc,
	}
}

func (rule *Rule) Name() string {
	return rule.name
}

func (rule *Rule) Enabled() bool {
	return true
}

func (rule *Rule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (rule *Rule) Link() string {
	return ""
}

func (rule *Rule) Metadata() interface{} {
	return nil
}

func (rule *Rule) Check(runner tflint.Runner) error {
	// config, err := runner.Config()
	// if err != nil {
	// 	return err
	// }

	// // Check if it is child module and do not evaluate them.
	// if len(config.Path) != 0 {
	// 	return nil
	// }

	return rule.checkFunc(runner, rule)
}
