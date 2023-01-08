package netlify

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
		},
	}
}

func resourceEnvVarCreate(d *schema.ResourceData, metaRaw interface{}) error {
	meta := metaRaw.(*Meta)

	// initialize creation parameters with default account ID, or supplied.
	params := operations.NewCreateEnvVarsParams()
	key := d.Get("key").(string)
	site_id := d.Get("site_id").(string)
	params.AccountID = d.Get("account_id").(string)
	params.SiteID = &site_id

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
	// if no scope found, set to all scopes
	if env_vars.Scopes == nil {
		env_vars.Scopes = []string{"builds", "functions", "post_processing", "runtime"}
	}

	// we need a placeholder value to be able to create the key
	env_vars.Values = []*models.EnvVarValue{
		{
			Context: "all",
			Value:   "",
		},
	}
	params.EnvVars = []*models.CreateEnvVarsParamsBodyItems{&env_vars}

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
	account_id, site_id, key := getEnvVarInfoFromResourceId(d.Id())
	params.AccountID = account_id
	params.SiteID = site_id
	params.Key = key

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

	// query for previous values, which we need to preserve
	params_get := operations.NewGetEnvVarParams()
	params_get.AccountID = account_id
	params_get.SiteID = site_id
	params_get.Key = key
	resp_get, err_get := meta.Netlify.Operations.GetEnvVar(params_get, meta.AuthInfo)
	if err_get != nil {
		return err_get
	}
	env_vars.Values = resp_get.Payload.Values
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
