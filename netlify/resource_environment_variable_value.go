package netlify

import (
	"context"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/netlify/open-api/v2/go/models"
	"github.com/netlify/open-api/v2/go/plumbing/operations"
)

func resourceEnvVarValue() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEnvVarValueCreateOrUpdate,
		ReadContext:   resourceEnvVarValueRead,
		UpdateContext: resourceEnvVarValueCreateOrUpdate,
		DeleteContext: resourceEnvVarValueDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"environment_variable_id": {
				Type:        schema.TypeString,
				Description: "The ID of a netlify_environment_variable resource this is a value of.",
				Required:    true,
				ForceNew:    true,
			},

			"context": {
				Type:        schema.TypeString,
				Description: "The deploy context in which this value will be used. `dev` refers to local development when running `netlify dev`. Enum: [`all` `dev` `branch-deploy` `deploy-preview` `production`]",
				Required:    true,
				ValidateDiagFunc: func(value interface{}, path cty.Path) diag.Diagnostics {
					for _, v := range []string{"all", "dev", "branch-deploy", "deploy-preview", "production"} {
						if v == value.(string) {
							return nil
						}
					}
					return diag.Diagnostics{
						diag.Diagnostic{
							Severity: diag.Error,
							Summary:  "Context level invalid.",
							Detail:   "Must be one of [`all` `dev` `branch-deploy` `deploy-preview` `production`]",
						},
					}
				},
			},

			"value": {
				Type:        schema.TypeString,
				Description: "The environment variable's unencrypted value",
				Required:    true,
			},
		},
	}
}

func resourceEnvVarValueCreateOrUpdate(c context.Context, d *schema.ResourceData, metaRaw interface{}) diag.Diagnostics {
	meta := metaRaw.(*Meta)

	// initialize creation parameters
	params := operations.NewSetEnvVarValueParams()
	account_id, site_id, key := getEnvVarInfoFromResourceId(d.Get("environment_variable_id").(string))
	params.AccountID = account_id
	params.SiteID = site_id
	params.Key = key
	params.EnvVar = &models.SetEnvVarValueParamsBody{
		Context: d.Get("context").(string),
		Value:   d.Get("value").(string),
	}

	// perform operation
	_, err := meta.Netlify.Operations.SetEnvVarValue(params, meta.AuthInfo)
	if err != nil {
		// default response is OK if it's just the default
		if v, ok := err.(*operations.SetEnvVarValueDefault); !ok && v == nil {
			return diag.FromErr(err)
		}
	}

	// re-read the environment variable from the backend
	return resourceEnvVarValueRead(c, d, metaRaw)
}

func resourceEnvVarValueRead(c context.Context, d *schema.ResourceData, metaRaw interface{}) diag.Diagnostics {
	meta := metaRaw.(*Meta)

	// initialize read parameters for top-level key
	params := operations.NewGetEnvVarParams()
	account_id, site_id, key := getEnvVarInfoFromResourceId(d.Get("environment_variable_id").(string))
	params.AccountID = account_id
	params.SiteID = site_id
	params.Key = key

	// perform operation
	resp, err := meta.Netlify.Operations.GetEnvVar(params, meta.AuthInfo)
	if err != nil {
		// If it is a 404 it was removed remotely
		if v, ok := err.(*operations.GetSiteDefault); ok && v.Code() == 404 {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	// find the environment variable value
	value_found := false
	for _, value := range resp.Payload.Values {
		if value.Context == d.Get("context").(string) {
			value_found = true
			d.SetId(value.ID)
			d.Set("context", value.Context)
			d.Set("value", value.Value)
		}
	}

	// if context was not found, env. variable value is missing
	if !value_found {
		d.SetId("")
	}

	return nil
}

func resourceEnvVarValueDelete(c context.Context, d *schema.ResourceData, metaRaw interface{}) diag.Diagnostics {
	meta := metaRaw.(*Meta)

	// initialize creation parameters for setting it to no-value
	params := operations.NewSetEnvVarValueParams()
	account_id, site_id, key := getEnvVarInfoFromResourceId(d.Get("environment_variable_id").(string))
	params.AccountID = account_id
	params.SiteID = site_id
	params.Key = key
	params.EnvVar = &models.SetEnvVarValueParamsBody{
		Context: d.Get("context").(string),
		Value:   "",
	}

	// perform operation
	_, err := meta.Netlify.Operations.SetEnvVarValue(params, meta.AuthInfo)
	if err != nil {
		// default response is OK if it's just the default
		if v, ok := err.(*operations.SetEnvVarValueDefault); !ok && v == nil {
			return diag.FromErr(err)
		}
	}
	return nil
}
