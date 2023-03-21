package rules

import (
	"fmt"
	"path/filepath"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

const (
	variablesFilename = "variables.tf"
	outputsFilename   = "outputs.tf"
)

const (
	wrongFileMessageTemplate         = "%s \"%s\" should be moved from %s to %s file"
	wrongResourceTypeMessageTemplate = "There should be no other resources in %s file except %s"
)

func NewModuleStructureRule() *Rule {
	return NewRule(
		"module_structure",
		func(runner tflint.Runner, rule tflint.Rule) error {
			if err := checkVariables(runner, rule); err != nil {
				return err
			}
			if err := checkOutputs(runner, rule); err != nil {
				return err
			}
			if err := checkVariablesFile(runner, rule); err != nil {
				return err
			}
			if err := checkOutputsFile(runner, rule); err != nil {
				return err
			}

			return nil
		},
	)
}

func checkVariablesFile(runner tflint.Runner, rule tflint.Rule) error {
	files, err := runner.GetFiles()
	if err != nil {
		return err
	}
	var variablesFile *hcl.File
	for filename, file := range files {
		if filepath.Base(filename) == variablesFilename {
			variablesFile = file
		}
	}
	if variablesFile == nil {
		return nil
	}

	body, ok := variablesFile.Body.(*hclsyntax.Body)
	if !ok {
		return nil
	}
	for _, block := range body.Blocks {
		if block.Type != "variable" {
			if err := runner.EmitIssue(
				rule,
				fmt.Sprintf(
					wrongResourceTypeMessageTemplate,
					block.Range().Filename,
					"variables",
				),
				block.Range(),
			); err != nil {
				return err
			}
		}
	}

	return nil
}

func checkVariables(runner tflint.Runner, rule tflint.Rule) error {
	content, err := runner.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type:       "variable",
				LabelNames: []string{"name"},
			},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, variable := range content.Blocks {
		if filepath.Base(variable.DefRange.Filename) != variablesFilename {
			if err := runner.EmitIssue(
				rule,
				fmt.Sprintf(
					wrongFileMessageTemplate,
					"variable",
					variable.Labels[0],
					variable.DefRange.Filename,
					filepath.Join(
						filepath.Dir(variable.DefRange.Filename),
						variablesFilename,
					),
				),
				variable.DefRange,
			); err != nil {
				return err
			}
		}
	}
	return nil
}

func checkOutputsFile(runner tflint.Runner, rule tflint.Rule) error {
	files, err := runner.GetFiles()
	if err != nil {
		return err
	}
	var outputsFile *hcl.File
	for filename, file := range files {
		if filepath.Base(filename) == outputsFilename {
			outputsFile = file
		}
	}
	if outputsFile == nil {
		return nil
	}

	body, ok := outputsFile.Body.(*hclsyntax.Body)
	if !ok {
		return nil
	}
	for _, block := range body.Blocks {
		if block.Type != "output" {
			if err := runner.EmitIssue(
				rule,
				fmt.Sprintf(
					wrongResourceTypeMessageTemplate,
					block.Range().Filename,
					"outputs",
				),
				block.Range(),
			); err != nil {
				return err
			}
		}
	}

	return nil
}

func checkOutputs(runner tflint.Runner, rule tflint.Rule) error {
	content, err := runner.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type:       "output",
				LabelNames: []string{"name"},
			},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, output := range content.Blocks {
		if filepath.Base(output.DefRange.Filename) != outputsFilename {
			if err := runner.EmitIssue(
				rule,
				fmt.Sprintf(
					wrongFileMessageTemplate,
					"output",
					output.Labels[0],
					output.DefRange.Filename,
					filepath.Join(
						filepath.Dir(output.DefRange.Filename),
						outputsFilename,
					),
				),
				output.DefRange,
			); err != nil {
				return err
			}
		}
	}

	return nil
}
