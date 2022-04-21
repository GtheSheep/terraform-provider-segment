package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

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
				DefaultFunc: schema.EnvDefaultFunc("SEGMENT_API_TOKEN", "https://api.segmentapis.com"),
				Description: "Base Api URL to use, i.e. https://eu1.api.segmentapis.com if your Segment account is hosted in the EU",
			},
		},
		DataSourcesMap:       map[string]*schema.Resource{},
		ResourcesMap:         map[string]*schema.Resource{},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {

	token := d.Get("token").(string)
	account_id := d.Get("account_id").(int)

	var diags diag.Diagnostics

	if (token != "") && (account_id != 0) {
		c, err := dbt_cloud.NewClient(&account_id, &token)

		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to login to DBT Cloud",
				Detail:   err.Error(),
			})
			return nil, diags
		}

		return c, diags
	}

	c, err := dbt_cloud.NewClient(nil, nil)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create DBT Cloud client",
			Detail:   err.Error(),
		})
		return nil, diags
	}

	return c, diags
}
