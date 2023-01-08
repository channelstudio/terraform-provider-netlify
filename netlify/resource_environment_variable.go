package netlify

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/netlify/open-api/v2/go/models"
	"github.com/netlify/open-api/v2/go/plumbing/operations"
)

func resourceEnvVar() *schema.Resource {
	return &schema.Resource{
		Create: resourceEnvVarCreate,
		Read:   resourceEnvVarRead,
		Update: resourceEnvVarUpdate,
		Delete: resourceEnvVarDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeString,
				Description: "The account ID / slug to create the environment variable for.",
				Required:    true,
				ForceNew:    true,
			},

			"site_id": {
				Type:        schema.TypeString,
				Description: "If provided, creates the environment variable on the site level, not the account level",
				Optional:    true,
				ForceNew:    true,
			},

			"key": {
				Type:        schema.TypeString,
				Description: "The name of the environment variable (case-sensitive).",
				Required:    true,
			},

			"scopes": {
				Type:        schema.TypeSet,
				Description: "The scopes that this environment variable is set to (Pro plans and above)",
				Optional:    true,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"value": {
				Type:        schema.TypeList,
				Description: "Values of this environment variable for a specific deploy context.",
				Required:    true,
				MinItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"context": {
							Type:        schema.TypeString,
							Description: "The deploy context in which this value will be used. `dev` refers to local development when running `netlify dev`. Enum: [`all` `dev` `branch-deploy` `deploy-preview` `production`]",
							Optional:    true,
						},

						"id": {
							Type:        schema.TypeString,
							Description: "The environment variable value's universally unique ID",
							Computed:    true,
						},

						"value": {
							Type:        schema.TypeString,
							Description: "The environment variable's unencrypted value",
							Required:    true,
						},
					},
				},
			},
		},
	}
}

func resourceEnvVarCreate(d *schema.ResourceData, metaRaw interface{}) error {
	meta := metaRaw.(*Meta)

	// initialize creation parameters with default account ID, or supplied.
	params := operations.NewCreateEnvVarsParams()
	site_id := d.Get("site_id").(string)
	params.AccountID = d.Get("account_id").(string)
	params.SiteID = &site_id
	key := d.Get("key").(string)

	// parse environment variables from struct
	environment_variables := []*models.CreateEnvVarsParamsBodyItems{}
	// build env vars create object
	env_vars := models.CreateEnvVarsParamsBodyItems{}
	env_vars.Key = key

	// set the scope of the environment variable
	if scopesI, ok := d.GetOk("scopes"); ok {
		scopesList := scopesI.(*schema.Set).List()
		if len(scopesList) > 0 {
			scopes := []string{}
			for _, scope := range scopesList {
				scopes = append(scopes, scope.(string))
			}
			env_vars.Scopes = scopes
		}
	}
	if env_vars.Scopes == nil {
		env_vars.Scopes = []string{"builds", "functions", "post_processing", "runtime"}
	}

	// set the values of the environment variable
	values := []*models.EnvVarValue{}
	for _, value := range d.Get("value").([]interface{}) {
		valueM := value.(map[string]interface{})
		values = append(values, &models.EnvVarValue{
			Context: valueM["context"].(string),
			ID:      valueM["id"].(string),
			Value:   valueM["value"].(string),
		})
	}
	env_vars.Values = values
	environment_variables = append(environment_variables, &env_vars)
	params.EnvVars = environment_variables

	// perform the operation
	_, err := meta.Netlify.Operations.CreateEnvVars(params, meta.AuthInfo)
	if err != nil {
		return err
	}

	// set the resource id from account ID, site ID, and key
	d.SetId(getResourceIdFromEnvVarInfo(params.AccountID, params.SiteID, key))
	return resourceEnvVarRead(d, metaRaw)
}

func resourceEnvVarRead(d *schema.ResourceData, metaRaw interface{}) error {
	meta := metaRaw.(*Meta)
	params := operations.NewGetEnvVarParams()
	// get account ID, site ID, and Key from resource ID
	split := strings.Split(d.Id(), "/")
	params.Key = split[0]
	params.AccountID = split[1]
	if len(split[2]) > 0 {
		site_id := split[2]
		params.SiteID = &site_id
	}

	resp, err := meta.Netlify.Operations.GetEnvVar(params, meta.AuthInfo)
	if err != nil {
		// If it is a 404, it was removed remotely
		if v, ok := err.(*operations.GetEnvVarDefault); ok && v.Code() == 404 {
			d.SetId("")
			return nil
		}
	}
	envVar := resp.Payload
	d.Set("key", envVar.Key)
	d.Set("scopes", envVar.Scopes)

	// otherwise, populate each of the values with their corresponding info
	for i, value := range envVar.Values {
		d.Set(fmt.Sprintf("value.%d", i), map[string]interface{}{
			"context": value.Context,
			"id":      value.ID,
			"value":   value.Value,
		})
	}

	return err
}

func resourceEnvVarUpdate(d *schema.ResourceData, metaRaw interface{}) error {
	meta := metaRaw.(*Meta)
	params := operations.NewUpdateEnvVarParams()
	// get previous account ID, site ID, and Key from resource ID
	account_id, site_id, key := getEnvVarInfoFromResourceId(d.Id())
	params.AccountID = account_id
	params.SiteID = site_id
	params.Key = key

	// build env vars update object
	env_vars := models.UpdateEnvVarParamsBody{}
	env_vars.Key = d.Get("key").(string)

	// set the scope of the environment variable
	if scopesI, ok := d.GetOk("scopes"); ok {
		scopesList := scopesI.(*schema.Set).List()
		if len(scopesList) > 0 {
			scopes := []string{}
			for _, scope := range scopesList {
				scopes = append(scopes, scope.(string))
			}
			env_vars.Scopes = scopes
		}
	}
	if env_vars.Scopes == nil {
		env_vars.Scopes = []string{"builds", "functions", "post_processing", "runtime"}
	}

	// set the values for each environment variable
	values := []*models.EnvVarValue{}
	for _, value := range d.Get("value").([]interface{}) {
		valueM := value.(map[string]interface{})
		values = append(values, &models.EnvVarValue{
			Context: valueM["context"].(string),
			ID:      valueM["id"].(string),
			Value:   valueM["value"].(string),
		})
	}
	env_vars.Values = values
	params.EnvVar = &env_vars

	// perform the operation
	resp, err := meta.Netlify.Operations.UpdateEnvVar(params, meta.AuthInfo)
	if err != nil {
		return err
	}

	envVar := resp.Payload

	// set the resource id (which may have changed) from account id, site id, and key
	d.SetId(getResourceIdFromEnvVarInfo(params.AccountID, params.SiteID, envVar.Key))
	return resourceEnvVarRead(d, metaRaw)
}

func resourceEnvVarDelete(d *schema.ResourceData, metaRaw interface{}) error {
	meta := metaRaw.(*Meta)
	params := operations.NewDeleteEnvVarParams()
	params.AccountID = d.Get("account_id").(string)
	site_id := d.Get("site_id").(string)
	params.SiteID = &site_id
	params.Key = d.Get("key").(string)
	_, err := meta.Netlify.Operations.DeleteEnvVar(params, meta.AuthInfo)
	return err
}

func getEnvVarInfoFromResourceId(id string) (account_id string, site_id *string, key string) {
	split := strings.Split(id, "/")
	key = split[0]
	account_id = split[1]
	if len(split[2]) > 0 {
		siteid := split[2]
		site_id = &siteid
	}
	return account_id, site_id, key
}

func getResourceIdFromEnvVarInfo(account_id string, site_id *string, key string) string {
	if site_id != nil {
		return fmt.Sprintf("%s/%s/%s", key, account_id, *site_id)
	} else {
		return fmt.Sprintf("%s/%s//", key, account_id)
	}
}
