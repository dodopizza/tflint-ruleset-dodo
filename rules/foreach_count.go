package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

const (
	foreachCountFirstArgumentMessage = "for_each/count should go as first argument on first line in resource"
	foreachCountDelimitedMessage     = "for_each/count should be delimited by empty line after it"
)

func NewForeachCountRule() *Rule {
	return NewRule(
		"foreach_count",
		func(runner tflint.Runner, rule tflint.Rule) error {

			resource, err := runner.GetModuleContent(&hclext.BodySchema{
				Blocks: []hclext.BlockSchema{
					{
						Type:       "resource",
						LabelNames: []string{"type", "name"},
						Body: &hclext.BodySchema{
							Mode: hclext.SchemaJustAttributesMode,
						},
					},
				},
			}, nil)
			if err != nil {
				return err
			}

			for _, block := range resource.Blocks {
				resource_block_start := block.DefRange.Start.Line
				var foreachCountAttribute *hclext.Attribute

				for _, attr := range block.Body.Attributes {
					if attr.Name == "for_each" || attr.Name == "count" {

						// Check if for_each or count is first argument
						if attr.Range.Start.Line != resource_block_start+1 {
							return runner.EmitIssue(
								rule,
								foreachCountFirstArgumentMessage,
								attr.Range,
							)
						}

						foreachCountAttribute = attr
					}
				}

				// Check if for_each or count is delimited by empty line
				if foreachCountAttribute == nil {
					continue
				}

				var firstAttr *hclext.Attribute
				for _, attr := range block.Body.Attributes {
					if attr.Name == "for_each" ||
						attr.Name == "count" {
						continue
					}

					if firstAttr == nil ||
						attr.Range.Start.Line < firstAttr.Range.Start.Line {
						firstAttr = attr
					}
				}

				if firstAttr.Range.Start.Line-foreachCountAttribute.Range.End.Line != 2 {
					return runner.EmitIssue(
						rule,
						foreachCountDelimitedMessage,
						firstAttr.Range,
					)
				}
			}

			return nil
		},
	)
}
