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
			body, err := runner.GetModuleContent(
				&hclext.BodySchema{
					Blocks: []hclext.BlockSchema{
						{Type: "backend", LabelNames: []string{"type"}, Body: &hclext.BodySchema{
							Attributes: []hclext.AttributeSchema{{Name: "key"}},
						}},
					},
				},
				&tflint.GetModuleContentOption{},
			)
			if err != nil {
				return err
			}

			for _, backend := range body.Blocks {
				if backend.Type != requiredBackendType || backend.Labels[0] != requiredYCBackendType {
					return runner.EmitIssue(
						rule,
						fmt.Sprintf(
							backendTypeMessageTemplate,
							requiredBackendType,
							backend.Type,
						),
						backend.DefRange,
					)
				}
			}

			return nil
		},
	)
}
