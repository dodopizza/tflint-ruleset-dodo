package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/stretchr/testify/require"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_ResourceNameLowercased(t *testing.T) {
	t.Parallel()

	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "no issues",
			Content: `
resource "null_resource" "test" {
  name = "test"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "no issues again",
			Content: `
resource "null_resource" "example" {
  name = "test"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "issue with first letter",
			Content: `
resource "null_resource" "Test" {
  name = "test"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewResourceNameLowercasedRule(),
					Message: "name is not lowercased: Test",
					Range: hcl.Range{
						Filename: filename,
					},
				},
			},
		},
		{
			Name: "issue with last letter",
			Content: `
resource "akamai_dns_record" "example_com_A" {
	zone       = akamai_dns_zone.example_com.zone
	target     = ["169.254.169.254"]
	name       = "example.com"
	recordtype = "A"
	ttl        = 60
}`,
			Expected: helper.Issues{
				{
					Rule:    NewResourceNameLowercasedRule(),
					Message: "name is not lowercased: example_com_A",
					Range: hcl.Range{
						Filename: filename,
					},
				},
			},
		},
	}
	rule := NewResourceNameLowercasedRule()

	for _, tc := range cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			runner := helper.TestRunner(t, map[string]string{filename: tc.Content})

			require.NoError(t, rule.Check(runner))
			helper.AssertIssuesWithoutRange(t, tc.Expected, runner.Issues)
		})
	}
}
