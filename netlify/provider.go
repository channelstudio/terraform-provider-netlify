package netlify

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	// schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
	// 	desc := s.Description
	// 	if s.Default != nil {
	// 		desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
	// 	}
	// 	return strings.TrimSpace(desc)
	// }
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"token": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("NETLIFY_TOKEN", nil),
					Description: "The OAuth token used to connect to Netlify.",
				},

				"base_url": {
					Type:        schema.TypeString,
					Optional:    true,
					DefaultFunc: schema.EnvDefaultFunc("NETLIFY_BASE_URL", defaultBaseUrl),
					Description: "The Netlify Base API URL",
				},
			},
			DataSourcesMap: map[string]*schema.Resource{
				"netlify_site": dataSourceSite(),
			},
			ResourcesMap: map[string]*schema.Resource{
				"netlify_build_hook":                 resourceBuildHook(),
				"netlify_branch_deploy":              resourceBranchDeploy(),
				"netlify_deploy_key":                 resourceDeployKey(),
				"netlify_hook":                       resourceHook(),
				"netlify_site":                       resourceSite(),
				"netlify_environment_variable":       resourceEnvVar(),
				"netlify_environment_variable_value": resourceEnvVarValue(),
				"netlify_dns_zone":                   resourceDnsZone(),
				"netlify_dns_record":                 resourceDnsRecord(),
			},
		}
		p.ConfigureContextFunc = configure(version, p)
		return p
	}
}

// The default Netlify base URL.
const defaultBaseUrl = "https://api.netlify.com/api/v1"

// configures the Netlify context to use with the provider
func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (any, diag.Diagnostics) {
	return func(c context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {
		config := Config{
			Token:   d.Get("token").(string),
			BaseURL: d.Get("base_url").(string),
		}
		client, err := config.Client()
		return client, diag.FromErr(err)
	}
}
