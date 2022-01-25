package rules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

const (
	emptyFirstLineMessageTemplate  = "File \"%s\" is starting with empty string"
	spaceBetweenObjectsMessage     = "Objects should be delimited with one empty line"
	noNewLineAtTheEndOfFileMessage = "There is no empty line at the end of file"
)

func NewFileContentRule() *Rule {
	return NewRule(
		"file_content",
		func(runner tflint.Runner, rule tflint.Rule) error {
			files, err := runner.Files()
			if err != nil {
				return err
			}

			for filename, file := range files {
				if err := checkFirstLine(runner, rule, filename, file); err != nil {
					return err
				}
				if err := checkLastLine(runner, rule, filename, file); err != nil {
					return err
				}
				if err := checkSpaceBetweenObjects(runner, rule, filename, file); err != nil {
					return err
				}
			}

			return nil
		},
	)
}

func checkFirstLine(
	runner tflint.Runner,
	rule tflint.Rule,
	filename string,
	file *hcl.File,
) error {
	lines := strings.Split(string(file.Bytes), "\n")
	if strings.Trim(lines[0], " ") != "" || len(lines) == 1 {
		return nil
	}

	if err := runner.EmitIssue(
		rule,
		fmt.Sprintf(emptyFirstLineMessageTemplate, filename),
		hcl.Range{
			Filename: filename,
			Start:    hcl.InitialPos,
			End:      hcl.InitialPos,
		},
	); err != nil {
		return err
	}

	return nil
}

func checkLastLine(
	runner tflint.Runner,
	rule tflint.Rule,
	filename string,
	file *hcl.File,
) error {
	lines := strings.Split(string(file.Bytes), "\n")
	if lines[len(lines)-1] == "" {
		return nil
	}

	if err := runner.EmitIssue(
		rule,
		noNewLineAtTheEndOfFileMessage,
		hcl.Range{
			Filename: filename,
			Start: hcl.Pos{
				Line: len(lines),
			},
			End: hcl.Pos{
				Line: len(lines),
			},
		},
	); err != nil {
		return err
	}

	return nil
}

func checkSpaceBetweenObjects(
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

	var endOfObjectFound bool
	var depth int
	for i := 0; i < len(tokens); i++ {
		if tokens[i].Type == hclsyntax.TokenOBrace {
			depth++
			continue
		}
		if tokens[i].Type == hclsyntax.TokenCBrace {
			depth--
			if depth == 0 {
				endOfObjectFound = true
			}

			continue
		}
		if !endOfObjectFound {
			continue
		}

		var newlineCounter int
		for {
			if tokens[i].Type == hclsyntax.TokenNewline {
				newlineCounter++
				i++
				continue
			}
			if tokens[i].Type == hclsyntax.TokenCBrace {
				newlineCounter = 0
				i++
				continue
			}
			if tokens[i].Type == hclsyntax.TokenComment {
				i++
				continue
			}

			break
		}
		if newlineCounter != 2 &&
			tokens[i].Type != hclsyntax.TokenEOF {
			if err := runner.EmitIssue(
				rule,
				spaceBetweenObjectsMessage,
				tokens[i].Range,
			); err != nil {
				return err
			}
		}
		endOfObjectFound = false
	}

	return nil
}
