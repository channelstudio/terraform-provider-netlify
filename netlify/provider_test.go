package netlify

import (
	"errors"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Use the same provider for all tests. This way it will be initialized so we
// can use its Netlify client to make our API calls for checking resources.
var testAccProvider *schema.Provider

// providerFactories are used to instantiate a provider during acceptance testing.
// The factory function will be invoked for every Terraform CLI command executed
// to create a provider server to which the CLI can reattach.
var providerFactories = map[string]func() (*schema.Provider, error){
	"netlify": func() (*schema.Provider, error) {
		return testAccProvider, nil
	},
}
var netlify *schema.Provider

func init() {
	testAccProvider = New("dev")()
}

func TestProvider(t *testing.T) {
	if err := New("dev")().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("NETLIFY_TOKEN"); v == "" {
		t.Fatal("NETLIFY_TOKEN must be set for acceptance tests")
	}
}

// common assertion test case
func testAccAssert(msg string, f func() bool) resource.TestCheckFunc {
	return func(*terraform.State) error {
		if f() {
			return nil
		}

		return errors.New("assertion failed: " + msg)
	}
}
