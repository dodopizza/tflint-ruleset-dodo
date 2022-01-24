package main

import (
	"github.com/terraform-linters/tflint-plugin-sdk/plugin"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"

	"github.com/dodopizza/tflint-ruleset-dodo/rules"
)

var version = "0.1.0"

func main() {
	plugin.Serve(&plugin.ServeOpts{
		RuleSet: &tflint.BuiltinRuleSet{
			Name:    "dodo",
			Version: version,
			Rules: []tflint.Rule{
				rules.NewBackendTypeRule(),
				rules.NewFileContentRule(),
				rules.NewCommentsRule(),
			},
		},
	})
}
