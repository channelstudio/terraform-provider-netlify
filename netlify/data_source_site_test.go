package netlify

import (
	"errors"
	"fmt"
	"math/rand"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/netlify/open-api/v2/go/models"
	"github.com/netlify/open-api/v2/go/plumbing/operations"
)

func TestAccDSSite(t *testing.T) {
	var site models.Site
	resourceName := "netlify_site.test"
	randomString := RandStringBytes(6)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckSiteDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDSSiteConfig, randomString, randomString),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSiteExists(resourceName, &site),
					testAccCheckSiteMatches("data.netlify_site.test", site),
				),
			},
		},
	})
}

func testAccCheckSiteMatches(n string, other models.Site) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		meta := testAccProvider.Meta().(*Meta)
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No membership ID is set")
		}

		// get the site on the backend
		params := operations.NewGetSiteParams()
		params.SiteID = rs.Primary.ID
		resp, err := meta.Netlify.Operations.GetSite(params, meta.AuthInfo)
		if err != nil {
			return err
		}
		site := resp.Payload

		// check values are the same
		if rs.Primary.Attributes["name"] != site.Name {
			return errors.New("name not matching")
		} else if rs.Primary.Attributes["account_slug"] != site.AccountSlug {
			return errors.New("account slug not matching")
		} else if rs.Primary.Attributes["account_name"] != site.AccountName {
			return errors.New("account name not matching")
		} else {
			return nil
		}
	}
}

var testAccDSSiteConfig = `
resource "netlify_site" "test" {
	name = "testing-site-%s"
}

data "netlify_site" "test" {
	name = "testing-site-%s"
	depends_on = [netlify_site.test]
}
`

const letterBytes = "abcdefghijklmnopqrstuvwxyz0123456789"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
