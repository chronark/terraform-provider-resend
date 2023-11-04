package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDomainResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "resend_domain" "test" {
  name = "resend.chronark.com"
  region = "us-east-1"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("resend_domain.test", "name", "resend.chronark.com"),
					resource.TestCheckResourceAttrSet("resend_domain.test", "id"),
					resource.TestCheckResourceAttr("resend_domain.test", "region", "us-east-1"),
					// resource.TestCheckResourceAttr("resend_domain.test", "records#", "2"),
					// resource.TestCheckResourceAttrSet("resend_domain.test", "records.0.record"),
					// resource.TestCheckResourceAttrSet("resend_domain.test", "records.0.name"),
					// resource.TestCheckResourceAttrSet("resend_domain.test", "records.0.type"),
					// resource.TestCheckResourceAttrSet("resend_domain.test", "records.0.ttl"),
					// resource.TestCheckResourceAttrSet("resend_domain.test", "records.0.status"),
					// resource.TestCheckResourceAttrSet("resend_domain.test", "records.0.value"),
					// resource.TestCheckResourceAttrSet("resend_domain.test", "records.0.priority"),
				),
			},
			// ImportState testing
			{
				ResourceName: "resend_domain.test",
				ImportState:  true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
