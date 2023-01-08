package netlify

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/netlify/open-api/v2/go/models"
	"github.com/netlify/open-api/v2/go/plumbing/operations"
)

func resourceDnsRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceDnsRecordCreate,
		Read:   resourceDnsRecordRead,
		Delete: resourceDnsRecordDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"hostname": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"value": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceDnsRecordCreate(d *schema.ResourceData, metaRaw interface{}) error {
	params := operations.NewCreateDNSRecordParams()
	params.ZoneID = d.Get("zone_id").(string)
	params.DNSRecord = &models.DNSRecordCreate{
		Hostname: d.Get("hostname").(string),
		Type:     d.Get("type").(string),
		Value:    d.Get("value").(string),
	}

	meta := metaRaw.(*Meta)
	resp, err := meta.Netlify.Operations.CreateDNSRecord(params, meta.AuthInfo)
	if err != nil {
		return err
	}

	d.SetId(resp.Payload.ID)
	return resourceDnsRecordRead(d, metaRaw)
}

func resourceDnsRecordRead(d *schema.ResourceData, metaRaw interface{}) error {
	meta := metaRaw.(*Meta)
	params := operations.NewGetIndividualDNSRecordParams()
	params.ZoneID = d.Get("zone_id").(string)
	params.DNSRecordID = d.Id()
	resp, err := meta.Netlify.Operations.GetIndividualDNSRecord(params, meta.AuthInfo)
	if err != nil {
		// If it is a 404 it was removed remotely
		if v, ok := err.(*operations.GetIndividualDNSRecordDefault); ok && v.Code() == 404 {
			d.SetId("")
			return nil
		}

		return err
	}

	zone := resp.Payload
	d.Set("site_id", zone.SiteID)
	d.Set("hostname", zone.Hostname)
	d.Set("type", zone.Type)
	d.Set("value", zone.Value)

	return nil
}

func resourceDnsRecordDelete(d *schema.ResourceData, metaRaw interface{}) error {
	meta := metaRaw.(*Meta)
	params := operations.NewDeleteDNSRecordParams()
	params.ZoneID = d.Get("zone_id").(string)
	params.DNSRecordID = d.Id()
	_, err := meta.Netlify.Operations.DeleteDNSRecord(params, meta.AuthInfo)
	return err
}
