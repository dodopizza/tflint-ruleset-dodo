package rules

import (
	"fmt"
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/require"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_ModuleStructure(t *testing.T) {
	t.Parallel()

	cases := []struct {
		Name     string
		Content  map[string]string
		Expected helper.Issues
	}{
		{
			Name: "no issues",
			Content: map[string]string{
				filename: `resource "null_resource" "test" {}`,
			},
			Expected: helper.Issues{},
		},
		{
			Name: "no issues with values",
			Content: map[string]string{
				variablesFile: `variable "test" {}`,
				outputsFile:   `output "o" { value = null }`,
			},
			Expected: helper.Issues{},
		},
		{
			Name: "variable in wrong file",
			Content: map[string]string{
				filename: `variable "test" {}`,
			},
			Expected: helper.Issues{
				{
					Rule: NewModuleStructureRule(),
					Message: fmt.Sprintf(
						wrongFileMessageTemplate,
						"variable",
						"test",
						filename,
						variablesFile,
					),
					Range: hcl.Range{
						Filename: filename,
						Start: hcl.Pos{
							Line:   1,
							Column: 1,
						},
						End: hcl.Pos{
							Line:   1,
							Column: 16,
						},
					},
				},
			},
		},
		{
			Name: "output in wrong file",
			Content: map[string]string{
				filename: `output "test" { value = null }`,
			},
			Expected: helper.Issues{
				{
					Rule: NewModuleStructureRule(),
					Message: fmt.Sprintf(
						wrongFileMessageTemplate,
						"output",
						"test",
						filename,
						outputsFile,
					),
					Range: hcl.Range{
						Filename: filename,
						Start: hcl.Pos{
							Line:   1,
							Column: 1,
						},
						End: hcl.Pos{
							Line:   1,
							Column: 14,
						},
					},
				},
			},
		},
	}
	rule := NewModuleStructureRule()

	for _, tc := range cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			runner := helper.TestRunner(t, tc.Content)

			require.NoError(t, rule.Check(runner))
			helper.AssertIssues(t, tc.Expected, runner.Issues)
		})
	}
}
