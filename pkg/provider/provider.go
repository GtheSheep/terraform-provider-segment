package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	//     "github.com/gthesheep/terraform-provider-segment/pkg/data_sources"
	"github.com/gthesheep/terraform-provider-segment/pkg/resources"
	"github.com/gthesheep/terraform-provider-segment/pkg/segment"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SEGMENT_API_TOKEN", nil),
				Description: "Public API token for your Segment account",
			},
			"api_url": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SEGMENT_API_URL", "https://api.segmentapis.com"),
				Description: "Base Api URL to use, i.e. https://eu1.api.segmentapis.com if your Segment account is hosted in the EU",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{},
		ResourcesMap: map[string]*schema.Resource{
			"segment_destination": resources.ResourceDestination(),
			"segment_source":      resources.ResourceSource(),
			"segment_warehouse":   resources.ResourceWarehouse(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {

	token := d.Get("token").(string)
	apiURL := d.Get("api_url").(string)

	var diags diag.Diagnostics

	if (token != "") && (apiURL != "") {
		c, err := segment.NewClient(apiURL, &token)

		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to login to Segment",
				Detail:   err.Error(),
			})
			return nil, diags
		}

		return c, diags
	}

	c, err := segment.NewClient("", nil)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create Segment client",
			Detail:   err.Error(),
		})
		return nil, diags
	}

	return c, diags
}
