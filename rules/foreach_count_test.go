package rules

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_NewForeachCountRule(t *testing.T) {
	tests := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{

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
			Name: "no for_each or count",
			Content: `
resource "null_resource" "test" {
  name = "test"
}
`,
			Expected: helper.Issues{},
		},
	}

	rule := NewForeachCountRule()

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{"resource.tf": test.Content})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			helper.AssertIssues(t, test.Expected, runner.Issues)
		})
	}
}
