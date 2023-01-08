package netlify

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/netlify/open-api/v2/go/models"
	"github.com/netlify/open-api/v2/go/plumbing/operations"
)

func resourceDnsZone() *schema.Resource {
	return &schema.Resource{
		Create: resourceDnsZoneCreate,
		Read:   resourceDnsZoneRead,
		Delete: resourceDnsZoneDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"site_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDnsZoneCreate(d *schema.ResourceData, metaRaw interface{}) error {
	params := operations.NewCreateDNSZoneParams()
	params.DNSZoneParams = &models.DNSZoneSetup{
		SiteID: d.Get("site_id").(string),
		Name:   d.Get("name").(string),
	}

	meta := metaRaw.(*Meta)
	resp, err := meta.Netlify.Operations.CreateDNSZone(params, meta.AuthInfo)
	if err != nil {
		return err
	}

	d.SetId(resp.Payload.ID)
	return resourceDnsZoneRead(d, metaRaw)
}

func resourceDnsZoneRead(d *schema.ResourceData, metaRaw interface{}) error {
	meta := metaRaw.(*Meta)
	params := operations.NewGetDNSZoneParams()
	params.ZoneID = d.Id()
	resp, err := meta.Netlify.Operations.GetDNSZone(params, meta.AuthInfo)
	if err != nil {
		// If it is a 404 it was removed remotely
		if v, ok := err.(*operations.GetDNSZoneDefault); ok && v.Code() == 404 {
			d.SetId("")
			return nil
		}

		return err
	}

	zone := resp.Payload
	d.Set("site_id", zone.SiteID)
	d.Set("name", zone.Name)
	d.Set("domain", zone.Domain)

	return nil
}

func resourceDnsZoneDelete(d *schema.ResourceData, metaRaw interface{}) error {
	meta := metaRaw.(*Meta)
	params := operations.NewDeleteDNSZoneParams()
	params.ZoneID = d.Id()
	_, err := meta.Netlify.Operations.DeleteDNSZone(params, meta.AuthInfo)
	return err
}
