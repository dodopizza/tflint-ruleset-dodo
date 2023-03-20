package rules

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

const (
	requiredBackendType        = "azurerm"
	requiredYCBackendType      = "s3"
	backendTypeMessageTemplate = "backend type should be \"%s\" but defined: \"%s\""
)

func NewBackendTypeRule() *Rule {
	return NewRule(
		"backend_type",
		func(runner tflint.Runner, rule tflint.Rule) error {
			content, err := runner.GetModuleContent(&hclext.BodySchema{
				Blocks: []hclext.BlockSchema{
					{
						Type: "terraform",
						Body: &hclext.BodySchema{
							Blocks: []hclext.BlockSchema{
								{
									Type:       "backend",
									LabelNames: []string{"type"},
								},
							},
						},
					},
				},
			}, nil)
			if err != nil {
				return err
			}
			for _, terraform := range content.Blocks {
				for _, backend := range terraform.Body.Blocks {
					if backend.Labels[0] != requiredYCBackendType && backend.Labels[0] != requiredBackendType {
						return runner.EmitIssue(
							rule,
							fmt.Sprintf(
								backendTypeMessageTemplate,
								requiredBackendType,
								backend.Labels[0],
							),
							backend.DefRange,
						)
					}
				}
			}

			return nil
		},
	)
}
