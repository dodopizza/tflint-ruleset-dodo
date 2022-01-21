package rules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

var (
	emptyFirstLineMessageTemplate  = "File \"%s\" is starting with empty string"
	noSpaceBetweenObjectsMessage   = "There is no empty line after object"
	noNewLineAtTheEndOfFileMessage = "There is no empty line at the end of file"
)

func NewFileContentRule() *BaseRule {
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
	if strings.Trim(lines[0], " ") != "" {
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

func checkSpaceBetweenObjects(
	runner tflint.Runner,
	rule tflint.Rule,
	filename string,
	file *hcl.File,
) error {
	lines := strings.Split(string(file.Bytes), "\n")

	var endOfObjectFound bool
	for i := 0; i < len(lines); i++ {
		if lines[i] == "}" {
			endOfObjectFound = true

			continue
		}
		if !endOfObjectFound {
			continue
		}
		if strings.HasPrefix(lines[i], "#") ||
			strings.HasPrefix(lines[i], "//") {
			continue
		}
		if lines[i] == "" {
			endOfObjectFound = false

			continue
		}

		if err := runner.EmitIssue(
			rule,
			noSpaceBetweenObjectsMessage,
			hcl.Range{
				Filename: filename,
				Start: hcl.Pos{
					Line: i + 1,
				},
				End: hcl.Pos{
					Line: i + 1,
				},
			},
		); err != nil {
			return err
		}

		endOfObjectFound = false
	}
	if endOfObjectFound {
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
	}

	return nil
}
