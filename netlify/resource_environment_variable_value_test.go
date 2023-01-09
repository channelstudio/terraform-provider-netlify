package netlify

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/netlify/open-api/v2/go/models"
	"github.com/netlify/open-api/v2/go/plumbing/operations"
)

func TestAccEnvVarValue_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckSiteAndEnvVarsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEnvVarValueConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEnvVarHasValue("var1_dev", "dev", "test"),
				),
			},
		},
	})
}

func TestAccEnvVarValue_disappears(t *testing.T) {
	var site models.Site
	destroy := func(*terraform.State) error {
		meta := testAccProvider.Meta().(*Meta)
		params := operations.NewDeleteSiteParams()
		params.SiteID = site.ID
		_, err := meta.Netlify.Operations.DeleteSite(params, meta.AuthInfo)
		return err
	}

	getEnvVars := func(*terraform.State) error {
		params := operations.NewGetEnvVarParams()
		params.AccountID = site.AccountSlug
		params.SiteID = &site.ID
		params.Key = "var1"
		meta := testAccProvider.Meta().(*Meta)
		resp, err := meta.Netlify.Operations.GetEnvVar(params, meta.AuthInfo)
		if err != nil {
			// If it is a 404 it was removed remotely
			if v, ok := err.(*operations.GetEnvVarDefault); ok && v.Code() == 404 {
				return nil
			}
			return err
		}
		return nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckSiteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSiteConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckSiteExists("netlify_site.test", &site),
					destroy,
					getEnvVars,
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccEnvVarValue_changesValue(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckSiteAndEnvVarsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEnvVarValueConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEnvVarHasValue("var1_dev", "dev", "test"),
				),
			},
			{
				Config: testAccEnvVarValueConfig_changeVal,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEnvVarHasValue("var1_dev", "dev", "test5"),
				),
			},
		},
	})

}

func testAccCheckEnvVarHasValue(name string, context string, val string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources["netlify_environment_variable_value."+name]
		if !ok {
			return fmt.Errorf("Not Found: var value %s", name)
		}

		// initialize read parameters for top-level key
		params := operations.NewGetEnvVarParams()
		account_id, site_id, key := getEnvVarInfoFromResourceId(rs.Primary.Attributes["environment_variable_id"])
		params.AccountID = account_id
		params.SiteID = site_id
		params.Key = key

		// perform operation
		meta := testAccProvider.Meta().(*Meta)
		resp, err := meta.Netlify.Operations.GetEnvVar(params, meta.AuthInfo)
		if err != nil {
			return err
		}

		// find the environment variable value
		for _, value := range resp.Payload.Values {
			if value.Context == context && value.Value == val {
				return nil
			}
		}

		return errors.New("Couldn't find environment variable value")
	}
}

var testAccEnvVarValueConfig = `
resource "netlify_site" "test" {}

resource "netlify_environment_variable" "var1" {
	account_id = netlify_site.test.account_slug
	site_id = netlify_site.test.id
	key	= "var1"
}

resource "netlify_environment_variable_value" "var1_dev" {
	environment_variable_id = netlify_environment_variable.var1.id
	context = "dev"
	value = "test"
}
`

var testAccEnvVarValueConfig_changeVal = `
resource "netlify_site" "test" {}

resource "netlify_environment_variable" "var1" {
	account_id = netlify_site.test.account_slug
	site_id = netlify_site.test.id
	key	= "var1"
}

resource "netlify_environment_variable_value" "var1_dev" {
	environment_variable_id = netlify_environment_variable.var1.id
	context = "dev"
	value = "test5"
}
`
