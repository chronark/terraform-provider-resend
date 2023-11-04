package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccApiKeyResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "resend_api_key" "test" {
  name = "terraform"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("resend_api_key.test", "name", "terraform"),
					resource.TestCheckResourceAttrSet("resend_api_key.test", "id"),
					resource.TestCheckResourceAttrSet("resend_api_key.test", "token"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "resend_api_key.test",
				ImportState:       true,
				ImportStateVerify: false,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
