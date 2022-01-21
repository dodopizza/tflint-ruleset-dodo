package rules

import (
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

const commentsMessage = "Single line comments should begin with \"//\""

func NewCommentsRule() *BaseRule {
	return NewRule(
		"comments",
		func(runner tflint.Runner, rule tflint.Rule) error {
			files, err := runner.Files()
			if err != nil {
				return err
			}

			for filename, file := range files {
				if err := checkComments(runner, rule, filename, file); err != nil {
					return err
				}
			}

			return nil
		},
	)
}

func checkComments(
	runner tflint.Runner,
	rule tflint.Rule,
	filename string,
	file *hcl.File,
) error {
	if strings.HasSuffix(filename, ".json") {
		return nil
	}

	tokens, diags := hclsyntax.LexConfig(file.Bytes, filename, hcl.InitialPos)
	if diags.HasErrors() {
		return diags
	}

	for _, token := range tokens {
		if token.Type != hclsyntax.TokenComment {
			continue
		}

		if strings.HasPrefix(string(token.Bytes), "#") {
			if err := runner.EmitIssue(
				rule,
				commentsMessage,
				token.Range,
			); err != nil {
				return err
			}
		}
	}

	return nil
}
