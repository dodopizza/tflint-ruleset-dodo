package rules

import (
	"fmt"
	"path/filepath"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/terraform/internal/configs"
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
	files, err := runner.Files()
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
	cfg, _ := runner.Config()

	for _, variable := range cfg.Module.Variables {
		if filepath.Base(variable.DeclRange.Filename) != variablesFilename {
			if err := runner.EmitIssue(
				rule,
				fmt.Sprintf(
					wrongFileMessageTemplate,
					"variable",
					variable.Name,
					variable.DeclRange.Filename,
					filepath.Join(
						filepath.Dir(variable.DeclRange.Filename),
						variablesFilename,
					),
				),
				variable.DeclRange,
			); err != nil {
				return err
			}
		}
	}

	return nil
}

func checkOutputsFile(runner tflint.Runner, rule tflint.Rule) error {
	files, err := runner.Files()
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
	outputs := getOutputs(runner)
	for _, output := range outputs {
		if filepath.Base(output.DeclRange.Filename) != outputsFilename {
			if err := runner.EmitIssue(
				rule,
				fmt.Sprintf(
					wrongFileMessageTemplate,
					"output",
					output.Name,
					output.DeclRange.Filename,
					filepath.Join(
						filepath.Dir(output.DeclRange.Filename),
						outputsFilename,
					),
				),
				output.DeclRange,
			); err != nil {
				return err
			}
		}
	}

	return nil
}

func getOutputs(runner tflint.Runner) []*configs.Output {
	files, _ := runner.Files()

	outputs := []*configs.Output{}
	for _, file := range files {
		body, ok := file.Body.(*hclsyntax.Body)
		if !ok {
			return nil
		}

		for _, block := range body.Blocks {
			if block.Type != "output" {
				continue
			}
			outputs = append(outputs, decodeOutputBlock(block))
		}
	}

	return outputs
}

func decodeOutputBlock(block *hclsyntax.Block) *configs.Output {
	return &configs.Output{
		Name:      block.Labels[0],
		DeclRange: block.Range(),
	}
}
