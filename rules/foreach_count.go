package rules

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/terraform/internal/configs"
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
			cfg, _ := runner.Config()
			for _, res := range cfg.Module.ManagedResources {
				r, ok := findForeachCount(res)
				if !ok {
					continue
				}

				if r.Start.Line-res.DeclRange.End.Line != 1 {
					if err := runner.EmitIssue(
						rule,
						foreachCountFirstArgumentMessage,
						r,
					); err != nil {
						return err
					}

					return nil
				}

				attrs, _ := res.Config.JustAttributes()
				var firstAttr *hcl.Attribute
				for _, attr := range attrs {
					if attr.Name == "for_each" ||
						attr.Name == "count" {
						continue
					}

					if firstAttr == nil ||
						attr.Range.Start.Line < firstAttr.Range.Start.Line {
						firstAttr = attr
					}
				}

				if firstAttr.Range.Start.Line-r.End.Line != 2 {
					if err := runner.EmitIssue(
						rule,
						foreachCountDelimitedMessage,
						r,
					); err != nil {
						return err
					}
				}
			}

			return nil
		},
	)
}

func findForeachCount(res *configs.Resource) (hcl.Range, bool) {
	var r hcl.Range
	if res.ForEach != nil {
		r = res.ForEach.Range()
	}
	if res.Count != nil {
		r = res.Count.Range()
	}
	if r.Start.Line == 0 {
		return r, false
	}

	return r, true
}
