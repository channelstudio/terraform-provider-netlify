package netlify

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/netlify/open-api/v2/go/models"
	"github.com/netlify/open-api/v2/go/plumbing/operations"
)

func dataSourceSite() *schema.Resource {
	return &schema.Resource{
		Description: "Queries a site within the Netlify account by name or ID.",
		ReadContext: dataSourceSiteRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Description:  "The name of the site. Required if ID is not specified.",
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"name", "site_id"},
			},
			"site_id": {
				Description:  "The ID of the site. Required if name is not specified.",
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"name", "site_id"},
			},
			"custom_domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"deploy_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"account_slug": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"account_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"repo": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"command": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"deploy_key_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"dir": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"provider": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"repo_path": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"repo_branch": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceSiteRead(ctx context.Context, d *schema.ResourceData, metaRaw interface{}) diag.Diagnostics {
	meta := metaRaw.(*Meta)
	var site *models.Site
	// if using ID, it's easy
	if id, ok := d.GetOk("site_id"); ok {
		params := operations.NewGetSiteParams()
		params.SiteID = id.(string)
		resp, err := meta.Netlify.Operations.GetSite(params, meta.AuthInfo)
		if err != nil {
			// If it is a 404 it was removed remotely
			if v, ok := err.(*operations.GetSiteDefault); ok && v.Code() == 404 {
				d.SetId("")
				return nil
			}
			return diag.FromErr(err)
		}
		site = resp.Payload
		bs, _ := json.Marshal(site)
		fmt.Print(string(bs))
		// otherwise, query all sites and look for ones that match
	} else {
		params := operations.NewListSitesParams()
		name := d.Get("name").(string)
		params.Name = &name
		resp, err := meta.Netlify.Operations.ListSites(params, meta.AuthInfo)
		if err != nil {
			return diag.FromErr(err)
		}
		// name match should be first one. if it's not, then name is not specific enough.
		sites := resp.Payload
		has_match := false
		for _, siteI := range sites {
			if siteI.Name == name {
				site = siteI
				has_match = true
				break
			}
		}
		if !has_match {
			d.SetId("")
			return nil
		}
	}

	// at this point, the correct site has been gotten, if it exists. so
	// we populate it exactly as we do for the site resource.
	d.SetId(site.ID)
	d.Set("site_id", site.ID)
	d.Set("name", site.Name)
	d.Set("custom_domain", site.CustomDomain)
	d.Set("deploy_url", site.DeployURL)
	d.Set("account_slug", site.AccountSlug)
	d.Set("account_name", site.AccountName)
	d.Set("repo", nil)

	if site.BuildSettings != nil && site.BuildSettings.RepoPath != "" {
		d.Set("repo", []interface{}{
			map[string]interface{}{
				"command":       site.BuildSettings.Cmd,
				"deploy_key_id": site.BuildSettings.DeployKeyID,
				"dir":           site.BuildSettings.Dir,
				"provider":      site.BuildSettings.Provider,
				"repo_path":     site.BuildSettings.RepoPath,
				"repo_branch":   site.BuildSettings.RepoBranch,
			},
		})
	}

	return nil
}
