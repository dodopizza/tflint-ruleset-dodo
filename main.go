package main

import (
	"github.com/terraform-linters/tflint-plugin-sdk/plugin"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"

	"github.com/dodopizza/tflint-ruleset-dodo/rules"
)

// Version to override during build.
var Version = "0.1.0"

func main() {
	plugin.Serve(&plugin.ServeOpts{
		RuleSet: &tflint.BuiltinRuleSet{
			Name:    "dodo",
			Version: Version,
			Rules: []tflint.Rule{
				rules.NewResourceNameLowercasedRule(),
				rules.NewTerraformBackendTypeRule(),
			},
		},
	})
}
