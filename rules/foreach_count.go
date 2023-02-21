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
							Attributes: []hclext.AttributeSchema{
								{Name: "for_each"},
								{Name: "count"},
							},
						},
					},
				},
			}, nil)
			if err != nil {
				return err
			}

			for _, block := range resource.Blocks {
				resource_block_start := block.DefRange.Start.Line
				for _, attr := range block.Body.Attributes {
					if attr.Name == "for_each" || attr.Name == "count" {

						// Check if for_each or count is first argument
						if attr.Range.Start.Line != resource_block_start+1 {
							runner.EmitIssue(
								rule,
								foreachCountFirstArgumentMessage,
								attr.Range,
							)
						}

						// Check if for_each or count is delimited by empty line
					}
				}
			}

			_ = resource

			// resources, err := runner.GetResourceContent(
			// 	"",
			// 	&hclext.BodySchema{
			// 		Attributes: []hclext.AttributeSchema{
			// 			{Name: "resource"},
			// 		},
			// 	},
			// 	nil,
			// )
			// if err != nil {
			// 	return err
			// }

			// logger.Debug(fmt.Sprintf("Get %d instances", len(resources.Blocks)))

			// for _, resource := range resources.Blocks {
			// 	if resource.Body.Attributes["resource"].Range.Start.Line != 1 {
			// 		runner.EmitIssue(
			// 			rule,
			// 			foreachCountFirstArgumentMessage,
			// 			resource.Body.Attributes["resource"].Range,
			// 		)
			// 	}
			// }

			return nil
		},
	)
}
