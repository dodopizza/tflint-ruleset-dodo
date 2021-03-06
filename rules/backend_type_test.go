package rules

import (
	"fmt"
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/require"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_BackendType(t *testing.T) {
	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "no issues",
			Content: `
terraform {
	backend "azurerm" {
		resource_group_name  = "rg"
		storage_account_name = "sa"
		container_name       = "tfstate"
		key    							 = "path/to/my/key"
	}
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "issue found",
			Content: `
terraform {
  backend "s3" {
		bucket = "mybucket"
		key    = "path/to/my/key"
    region = "us-east-1"
  }
}`,
			Expected: helper.Issues{
				{
					Rule: NewBackendTypeRule(),
					Message: fmt.Sprintf(
						backendTypeMessageTemplate,
						requiredBackendType,
						"s3",
					),
					Range: hcl.Range{
						Filename: filename,
						Start:    hcl.Pos{Line: 3, Column: 3},
						End:      hcl.Pos{Line: 3, Column: 15},
					},
				},
			},
		},
	}

	rule := NewBackendTypeRule()

	for _, tc := range cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{filename: tc.Content})

			require.NoError(t, rule.Check(runner))
			helper.AssertIssues(t, tc.Expected, runner.Issues)
		})
	}
}
