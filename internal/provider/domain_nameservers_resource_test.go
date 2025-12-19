package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDomainNameServersResource(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("Acceptance tests skipped unless env 'TF_ACC' set")
	}

	domain := os.Getenv("PORKBUN_TEST_DOMAIN")
	if domain == "" {
		t.Fatal("PORKBUN_TEST_DOMAIN environment variable must be set")
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing with custom nameservers
			{
				Config: testAccDomainNameServersResourceConfig(domain, []string{
					"curitiba.ns.porkbun.com",
					"fortaleza.ns.porkbun.com",
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("porkbun_domain_nameservers.test", "domain", domain),
					resource.TestCheckResourceAttr("porkbun_domain_nameservers.test", "nameservers.#", "2"),
				),
			},
			// Update with more nameservers
			{
				Config: testAccDomainNameServersResourceConfig(domain, []string{
					"curitiba.ns.porkbun.com",
					"fortaleza.ns.porkbun.com",
					"maceio.ns.porkbun.com",
					"salvador.ns.porkbun.com",
				}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("porkbun_domain_nameservers.test", "domain", domain),
					resource.TestCheckResourceAttr("porkbun_domain_nameservers.test", "nameservers.#", "4"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "porkbun_domain_nameservers.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDomainNameServersResourceConfig(domain string, nameservers []string) string {
	nsStr := ""
	for _, ns := range nameservers {
		nsStr += fmt.Sprintf(`    "%s",`+"\n", ns)
	}

	return fmt.Sprintf(`
resource "porkbun_domain_nameservers" "test" {
  domain      = %q
  nameservers = [
%s  ]
}
`, domain, nsStr)
}
