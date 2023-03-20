package rules

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/require"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_Comments(t *testing.T) {
	t.Parallel()

	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "no comments",
			Content: `
resource "null_resource" "test" {
  name = "test"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "no issues with one-line comment",
			Content: `
// here is null resource
resource "null_resource" "test" {
  name = "test"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "no issues multiline comment",
			Content: `
/* here is
multiline
comment */
resource "null_resource" "test" {
  name = "test"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "bash style comment",
			Content: `
# null resource
resource "null_resource" "test" {
  name = "test"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewCommentsRule(),
					Message: commentsMessage,
					Range: hcl.Range{
						Filename: filename,
						Start: hcl.Pos{
							Line:   2,
							Column: 1,
						},
						End: hcl.Pos{
							Line:   3,
							Column: 1,
						},
					},
				},
			},
		},
	}
	rule := NewCommentsRule()

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
