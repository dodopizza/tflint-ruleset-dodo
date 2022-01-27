package rules

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/terraform/configs"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

const (
	variablesFile = "variables.tf"
	outputsFile   = "outputs.tf"
)

const wrongFileMessageTemplate = "%s \"%s\" should be moved from %s to %s file"

var outputsSchema = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type:       "output",
			LabelNames: []string{"name"},
		},
	},
}

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

			return nil
		},
	)
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
		content, _ := file.Body.Content(outputsSchema)

		for _, block := range content.Blocks {
			if block.Type != "output" {
				continue
			}
			outputs = append(outputs, decodeOutputBlock(block))
		}
	}

	return outputs
}

func decodeOutputBlock(block *hcl.Block) *configs.Output {
	return &configs.Output{
		Name:      block.Labels[0],
		DeclRange: block.DefRange,
	}
}
