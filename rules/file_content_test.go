package rules

import (
	"fmt"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/require"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_FileContent(t *testing.T) {
	t.Parallel()

	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "no issues",
			Content: `resource "null_resource" "test" {
  name = "test"
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "no issues with multiple resources",
			Content: `resource "null_resource" "test" {
  name = "test"
}

resource "null_resource" "another_test" {
  name = "another_test"
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "no issues with multiple resources and comments",
			Content: `// test resource
resource "null_resource" "test" {
  name = "test"
}
// - test resource

resource "null_resource" "another_test" {
  name = "another_test"
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "empty first line",
			Content: `
resource "null_resource" "test" {
  name = "test"
}
`,
			Expected: helper.Issues{
				{
					Rule: NewFileContentRule(),
					Message: fmt.Sprintf(
						emptyFirstLineMessageTemplate,
						filename,
					),
					Range: hcl.Range{
						Filename: filename,
						Start:    hcl.InitialPos,
						End:      hcl.InitialPos,
					},
				},
			},
		},
		{
			Name: "no new line at the end",
			Content: `resource "null_resource" "test" {
  name = "test"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewFileContentRule(),
					Message: noNewLineAtTheEndOfFileMessage,
					Range: hcl.Range{
						Filename: filename,
						Start: hcl.Pos{
							Line: 3,
						},
						End: hcl.Pos{
							Line: 3,
						},
					},
				},
			},
		},
		{
			Name: "no new line between resources",
			Content: `resource "null_resource" "test" {
  name = "test"
}
resource "null_resource" "another_test" {
  name = "another_test"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewFileContentRule(),
					Message: noSpaceBetweenObjectsMessage,
					Range: hcl.Range{
						Filename: filename,
						Start: hcl.Pos{
							Line: 4,
						},
						End: hcl.Pos{
							Line: 4,
						},
					},
				},
			},
		},
		{
			Name: "no new line between resources with comments",
			Content: `// null resource
resource "null_resource" "test" {
  name = "test"
}
// - null resource
resource "null_resource" "another_test" {
  name = "another_test"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewFileContentRule(),
					Message: noSpaceBetweenObjectsMessage,
					Range: hcl.Range{
						Filename: filename,
						Start: hcl.Pos{
							Line: 6,
						},
						End: hcl.Pos{
							Line: 6,
						},
					},
				},
			},
		},
	}
	rule := NewFileContentRule()

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
