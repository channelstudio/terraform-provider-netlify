package netlify

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/netlify/open-api/v2/go/models"
)

func TestAccEnvVarValue_basic(t *testing.T) {
	var site models.Site
	var envVar1 models.EnvVar
	var envVar2 models.EnvVar
	resourceName := "netlify_site.test"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckSiteAndEnvVarsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEnvVarValueConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSiteExists(resourceName, &site),
					testAccCheckEnvVarExists("var1", "var1", &envVar1),
					testAccCheckEnvVarExists("var2", "var2", &envVar2),
				),
			},
		},
	})
}

// func testAccCheckEnvVarHasValue(key string, context string, value string) resource.TestCheckFunc {
// 	return func(s *terraform.State) error {
// 		rs, ok := s.RootModule().Resources["netlify_environment_variable."+resource_name]
// 		// initialize read parameters for top-level key
// 		params := operations.NewGetEnvVarParams()
// 		account_id, site_id, key := getEnvVarInfoFromResourceId(d.Get("environment_variable_id").(string))
// 		params.AccountID = account_id
// 		params.SiteID = site_id
// 		params.Key = key

// 		// perform operation
// 		resp, err := meta.Netlify.Operations.GetEnvVar(params, meta.AuthInfo)
// 		if err != nil {
// 			// If it is a 404 it was removed remotely
// 			if v, ok := err.(*operations.GetSiteDefault); ok && v.Code() == 404 {
// 				d.SetId("")
// 				return nil
// 			}
// 			return err
// 		}

// 		// find the environment variable value
// 		value_found := false
// 		for _, value := range resp.Payload.Values {
// 			if value.Context == d.Get("context").(string) {
// 				value_found = true
// 				d.SetId(value.ID)
// 				d.Set("context", value.Context)
// 				d.Set("value", value.Value)
// 			}
// 		}
// 	}
// }

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
