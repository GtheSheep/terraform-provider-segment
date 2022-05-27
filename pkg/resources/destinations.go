package resources

import (
	"context"
	"fmt"
	"strconv"

	"github.com/gthesheep/terraform-provider-segment/pkg/segment"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceDestination() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDestinationCreate,
		ReadContext:   resourceDestinationRead,
		UpdateContext: resourceDestinationUpdate,
		DeleteContext: resourceDestinationDelete,

		Schema: map[string]*schema.Schema{
			"enabled": &schema.Schema{
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Flag for whether or not the destination is enabled",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Descriptive name for the destination",
			},
			"destination_slug": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Slug for the destination, from a list",
				ValidateFunc: validation.StringInSlice(segment.DestinationSlugs, false),
			},
			"source_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Identifier of the source to connect this destination to",
			},
			"settings": &schema.Schema{
				Type:        schema.TypeMap,
				Required:    true,
				Description: "Map containing settings for the destination, currently all values must be provided as strings",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceDestinationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*segment.Client)

	var diags diag.Diagnostics

	enabled := d.Get("enabled").(bool)
	name := d.Get("name").(string)
	sourceID := d.Get("source_id").(string)
	destinationSlug := d.Get("destination_slug").(string)
	settings := d.Get("settings").(map[string]interface{})

	destination, err := c.CreateDestination(sourceID, enabled, name, destinationSlug, settings)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprintf("%s", *destination.ID))

	resourceDestinationRead(ctx, d, m)

	return diags
}

func resourceDestinationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*segment.Client)

	var diags diag.Diagnostics

	destinationID := d.Id()

	destination, err := c.GetDestination(destinationID)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("enabled", destination.Enabled); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", destination.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("destination_slug", destination.Metadata.Slug); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("source_id", destination.SourceID); err != nil {
		return diag.FromErr(err)
	}
	for k, v := range destination.Settings {
		if _, ok := v.(bool); ok {
			destination.Settings[k] = strconv.FormatBool(v.(bool))
		}
	}
	if err := d.Set("settings", destination.Settings); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceDestinationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*segment.Client)

	destinationID := d.Id()

	if d.HasChange("name") || d.HasChange("enabled") || d.HasChange("settings") {
		destination, err := c.GetDestination(destinationID)
		if err != nil {
			return diag.FromErr(err)
		}

		if d.HasChange("name") {
			name := d.Get("name").(string)
			destination.Name = name
		}
		if d.HasChange("enabled") {
			enabled := d.Get("enabled").(bool)
			destination.Enabled = enabled
		}
		if d.HasChange("settings") {
			settings := d.Get("settings").(map[string]interface{})
			destination.Settings = settings
		}

		_, err = c.UpdateDestination(*destination.ID, destination.SourceID, destination.Enabled, destination.Name, destination.Settings)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceDestinationRead(ctx, d, m)
}

func resourceDestinationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*segment.Client)

	var diags diag.Diagnostics

	destinationID := d.Id()

	_, err := c.DeleteDestination(destinationID)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
