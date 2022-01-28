package rules

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/terraform/configs"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

const (
	variablesFile = "variables.tf"
	outputsFile   = "outputs.tf"
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
	file, err := runner.File(variablesFile)
	if err != nil {
		return err
	}
	if file == nil {
		return nil
	}

	body, ok := file.Body.(*hclsyntax.Body)
	if !ok {
		return nil
	}
	for _, block := range body.Blocks {
		if block.Type != "variable" {
			if err := runner.EmitIssue(
				rule,
				fmt.Sprintf(
					wrongResourceTypeMessageTemplate,
					variablesFile,
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
		if variable.DeclRange.Filename != variablesFile {
			if err := runner.EmitIssue(
				rule,
				fmt.Sprintf(
					wrongFileMessageTemplate,
					"variable",
					variable.Name,
					variable.DeclRange.Filename,
					variablesFile,
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
	file, err := runner.File(outputsFile)
	if err != nil {
		return err
	}
	if file == nil {
		return nil
	}

	body, ok := file.Body.(*hclsyntax.Body)
	if !ok {
		return nil
	}
	for _, block := range body.Blocks {
		if block.Type != "output" {
			if err := runner.EmitIssue(
				rule,
				fmt.Sprintf(
					wrongResourceTypeMessageTemplate,
					outputsFile,
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
		if output.DeclRange.Filename != outputsFile {
			if err := runner.EmitIssue(
				rule,
				fmt.Sprintf(
					wrongFileMessageTemplate,
					"output",
					output.Name,
					output.DeclRange.Filename,
					outputsFile,
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
