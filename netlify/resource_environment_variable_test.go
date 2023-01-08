package netlify

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/netlify/open-api/v2/go/models"
	"github.com/netlify/open-api/v2/go/plumbing/operations"
)

func TestAccEnvVar_basic(t *testing.T) {
	var site models.Site
	var envVar1 models.EnvVar
	var envVar2 models.EnvVar
	resourceName := "netlify_site.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSiteAndEnvVarsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEnvVarSiteSpecificConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSiteExists(resourceName, &site),
					testAccCheckEnvVarExists("var1", "var1", &envVar1),
					testAccCheckEnvVarExists("var2", "var2", &envVar2),
				),
			},
		},
	})
}

func TestAccEnvVar_disappears(t *testing.T) {
	var site models.Site
	var envVar models.EnvVar
	var envVar2 models.EnvVar

	destroy := func(*terraform.State) error {
		meta := testAccProvider.Meta().(*Meta)
		params := operations.NewDeleteEnvVarParams()
		params.AccountID = site.AccountSlug
		params.SiteID = &site.ID
		params.Key = envVar.Key
		_, err := meta.Netlify.Operations.DeleteEnvVar(params, meta.AuthInfo)
		return err
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSiteAndEnvVarsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEnvVarSiteSpecificConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckSiteExists("netlify_site.test", &site),
					testAccCheckEnvVarExists("var1", "var1", &envVar),
					testAccCheckEnvVarExists("var2", "var2", &envVar2),
					destroy,
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccEnvVar_updatesKey(t *testing.T) {
	var envVar models.EnvVar

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSiteAndEnvVarsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEnvVarSiteSpecificConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckEnvVarExists("var1", "var1", &envVar),
				),
			},
			{
				Config: testAccEnvVarSiteSpecificConfigUpdateKey,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckEnvVarExists("var1", "var3", &envVar),
				),
			},
		},
	})
}

func TestAccEnvVar_updatesValue(t *testing.T) {
	var envVar models.EnvVar
	var envVarUpdated models.EnvVar

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSiteAndEnvVarsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEnvVarSiteSpecificConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckEnvVarExists("var1", "var1", &envVar),
				),
			},
			{
				Config: testAccEnvVarSiteSpecificConfigUpdateValue,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckEnvVarExists("var1", "var1", &envVarUpdated),
					testAccAssert("environment variable changed contexts and values", func() bool {
						for _, v := range envVar.Values {
							for _, v2 := range envVarUpdated.Values {
								if v2.ID == v.ID && (v2.Context == v.Context || v2.Value == v.Value) {
									return false
								}
							}
						}
						return true
					}),
				),
			},
		},
	})
}

func testAccCheckEnvVarExists(resource_name string, key string, envvar *models.EnvVar) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources["netlify_environment_variable."+resource_name]
		if !ok {
			return fmt.Errorf("Not Found: var %s", key)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No environment variable ID is set")
		}
		meta := testAccProvider.Meta().(*Meta)
		params := operations.NewGetEnvVarParams()
		params.AccountID = rs.Primary.Attributes["account_id"]
		site_id := rs.Primary.Attributes["site_id"]
		params.SiteID = &site_id
		params.Key = key
		resp, err := meta.Netlify.Operations.GetEnvVar(params, meta.AuthInfo)
		if err != nil {
			return err
		}

		*envvar = *resp.Payload
		return nil
	}
}

func testAccCheckSiteAndEnvVarsDestroy(s *terraform.State) error {
	err := testAccCheckSiteDestroy(s)
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "netlify_environment_variable" {
			continue
		}

		meta := testAccProvider.Meta().(*Meta)
		params := operations.NewGetEnvVarParams()
		params.AccountID = rs.Primary.Attributes["account_id"]
		site_id := rs.Primary.Attributes["site_id"]
		params.SiteID = &site_id
		params.Key = rs.Primary.Attributes["key"]
		resp, err := meta.Netlify.Operations.GetEnvVar(params, meta.AuthInfo)
		if err == nil && resp.Payload != nil {
			return fmt.Errorf("Environment variable still exists: %s", rs.Primary.ID)
		}

		if err != nil {
			if v, ok := err.(*operations.GetEnvVarDefault); ok && v.Code() == 404 {
				return nil
			}
		}
	}
	return nil
}

var testAccEnvVarSiteSpecificConfig = `
resource "netlify_site" "test" {}

resource "netlify_environment_variable" "var1" {
	account_id = netlify_site.test.account_slug
	site_id = netlify_site.test.id
	key	= "var1"
	value {
		context = "dev"
		value = "SOME VALUE"
	}
	value {
		context = "production"
		value = "SOME VALUE"
	}
}

resource "netlify_environment_variable" "var2" {
	account_id = netlify_site.test.account_slug
	site_id = netlify_site.test.id
	key	= "var2"
	value {
		context = "dev"
		value = "SOME VALUE 1"
	}
	value {
		context = "production"
		value = "SOME VALUE 2"
	}
}
`

var testAccEnvVarSiteSpecificConfigUpdateKey = `
resource "netlify_site" "test" {}

resource "netlify_environment_variable" "var1" {
	account_id = netlify_site.test.account_slug
	site_id = netlify_site.test.id
	key	= "var3"
	value {
		context = "dev"
		value = "SOME VALUE"
	}
	value {
		context = "production"
		value = "SOME VALUE"
	}
}

resource "netlify_environment_variable" "var2" {
	account_id = netlify_site.test.account_slug
	site_id = netlify_site.test.id
	key	= "var2"
	value {
		context = "dev"
		value = "SOME VALUE 1"
	}
	value {
		context = "production"
		value = "SOME VALUE 2"
	}
}
`

var testAccEnvVarSiteSpecificConfigUpdateValue = `
resource "netlify_site" "test" {}

resource "netlify_environment_variable" "var1" {
	account_id = netlify_site.test.account_slug
	site_id = netlify_site.test.id
	key	= "var1"
	value {
		context = "all"
		value = "changeddev"
	}
	value {
		context = "deploy-preview"
		value = "changedprod"
	}
}

resource "netlify_environment_variable" "var2" {
	account_id = netlify_site.test.account_slug
	site_id = netlify_site.test.id
	key	= "var2"
	value {
		context = "dev"
		value = "SOME VALUE 1"
	}
	value {
		context = "production"
		value = "SOME VALUE 2"
	}
}
`
