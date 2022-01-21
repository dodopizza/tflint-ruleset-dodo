package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

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
		name:      name,
		checkFunc: checkFunc,
	}
}

func (r *BaseRule) Name() string {
	return r.name
}

func (r *BaseRule) Enabled() bool {
	return true
}

func (r *BaseRule) Severity() string {
	return tflint.ERROR
}

func (r *BaseRule) Link() string {
	return ""
}

func (r *BaseRule) Check(runner tflint.Runner) error {
	return r.checkFunc(runner, r)
}
