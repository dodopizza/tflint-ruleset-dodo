package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/require"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_ForeachCount(t *testing.T) {
	t.Parallel()

	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "no for_each or count",
			Content: `
resource "null_resource" "test" {
  name = "test"
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "no issues for_each",
			Content: `
resource "null_resource" "test" {
	for_each = to_set(["test", "test"])

	name = each.key
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "no issues for_each multiline map",
			Content: `
resource "null_resource" "test" {
	for_each = {
		test = "test"
		key  = "value"
	}

  name        = each.key
	description = each.value
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "no issues nested for_each",
			Content: `
resource "null_resource" "test" {
	name = "test"

	dynamic "config" {
		for_each = toset([1, 2, 3, 4, 5])

		content {
			key = value
		}
	}
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "no issues for_each for var",
			Content: `
locals {
	test = {
		name = "value"
	}
}

resource "null_resource" "test" {
  for_each = { for t in local.t : t.name => t }

  name = each.value
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "for_each not on first line",
			Content: `
resource "null_resource" "test" {

	for_each = to_set(["test", "test"])

	name = each.key
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewForeachCountRule(),
					Message: foreachCountFirstArgumentMessage,
					Range: hcl.Range{
						Filename: filename,
						Start: hcl.Pos{
							Line:   4,
							Column: 13,
						},
						End: hcl.Pos{
							Line:   4,
							Column: 37,
						},
					},
				},
			},
		},
		{
			Name: "for_each not as first argument",
			Content: `
resource "null_resource" "test" {
	name = each.key
	for_each = to_set(["test", "test"])
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewForeachCountRule(),
					Message: foreachCountFirstArgumentMessage,
					Range: hcl.Range{
						Filename: filename,
						Start: hcl.Pos{
							Line:   4,
							Column: 13,
						},
						End: hcl.Pos{
							Line:   4,
							Column: 37,
						},
					},
				},
			},
		},
		{
			Name: "count not as first argument",
			Content: `
resource "null_resource" "test" {
	name = "count.index"
	count = 2
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewForeachCountRule(),
					Message: foreachCountFirstArgumentMessage,
					Range: hcl.Range{
						Filename: filename,
						Start: hcl.Pos{
							Line:   4,
							Column: 10,
						},
						End: hcl.Pos{
							Line:   4,
							Column: 11,
						},
					},
				},
			},
		},
	}
	rule := NewForeachCountRule()

	for _, tc := range cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			runner := helper.TestRunner(t, map[string]string{filename: tc.Content})

			require.NoError(t, rule.Check(runner))
			helper.AssertIssues(t, tc.Expected, runner.Issues)
		})
	}
}
