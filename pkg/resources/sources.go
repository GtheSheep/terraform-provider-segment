package resources

import (
	"context"
	"fmt"

	"github.com/gthesheep/terraform-provider-segment/pkg/segment"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var (
	ViolationEvents = []string{
		"ALLOW",
		"BLOCK",
		"OMIT_PROPERTIES",
	}
)

func ResourceSource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSourceCreate,
		ReadContext:   resourceSourceRead,
		UpdateContext: resourceSourceUpdate,
		DeleteContext: resourceSourceDelete,

		Schema: map[string]*schema.Schema{
			"slug": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Slug for the source, lower case",
			},
			"enabled": &schema.Schema{
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Flag for whether or not the source is enabled",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Descriptive name for the source",
			},
			"source_slug": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Slug for the source, from a list",
				ValidateFunc: validation.StringInSlice(segment.SourceSlugs, false),
			},
			"settings": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Map containing settings for the source",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"forwarding_violations_to": {
							Description: "SourceId to forward violations to.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"forwarding_blocked_events_to": {
							Description: "SourceId to forward blocked events to.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"track": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"allow_unplanned_events": {
										Description: "Enable to allow unplanned track events.",
										Type:        schema.TypeBool,
										Required:    true,
									},
									"allow_unplanned_event_properties": {
										Description: "Enable to allow unplanned track event properties.",
										Type:        schema.TypeBool,
										Required:    true,
									},
									"allow_event_on_violations": {
										Description: "Allow track event on violations.",
										Type:        schema.TypeBool,
										Required:    true,
									},
									"allow_properties_on_violations": {
										Description: "Enable to allow track properties on violations.",
										Type:        schema.TypeBool,
										Required:    true,
									},
									"common_event_on_violations": {
										Description:  "The common track event on violations.",
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice(ViolationEvents, false),
									},
								},
							},
						},
						"identify": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"allow_unplanned_traits": {
										Description: "Enable to allow unplanned identify traits.",
										Type:        schema.TypeBool,
										Required:    true,
									},
									"allow_traits_on_violations": {
										Description: "Enable to allow identify traits on violations.",
										Type:        schema.TypeBool,
										Required:    true,
									},
									"common_event_on_violations": {
										Description:  "The common track event on violations.",
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice(ViolationEvents, false),
									},
								},
							},
						},
						"group": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"allow_unplanned_traits": {
										Description: "Enable to allow unplanned identify traits.",
										Type:        schema.TypeBool,
										Required:    true,
									},
									"allow_traits_on_violations": {
										Description: "Enable to allow identify traits on violations.",
										Type:        schema.TypeBool,
										Required:    true,
									},
									"common_event_on_violations": {
										Description:  "The common track event on violations.",
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice(ViolationEvents, false),
									},
								},
							},
						},
					},
				},
			},
		},

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceSourceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*segment.Client)

	var diags diag.Diagnostics

	slug := d.Get("slug").(string)
	enabled := d.Get("enabled").(bool)
	name := d.Get("name").(string)
	sourceSlug := d.Get("source_slug").(string)

	source, err := c.CreateSource(slug, enabled, name, sourceSlug)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprintf("%s", *source.ID))

	resourceSourceRead(ctx, d, m)

	return diags
}

func resourceSourceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*segment.Client)

	var diags diag.Diagnostics

	sourceID := d.Id()

	source, err := c.GetSource(sourceID)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("slug", source.Slug); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("enabled", source.Enabled); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", source.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("source_slug", source.Metadata.Slug); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceSourceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*segment.Client)

	sourceID := d.Id()

	if d.HasChange("name") || d.HasChange("enabled") {
		source, err := c.GetSource(sourceID)
		if err != nil {
			return diag.FromErr(err)
		}

		if d.HasChange("name") {
			name := d.Get("name").(string)
			source.Name = name
		}
		if d.HasChange("enabled") {
			enabled := d.Get("enabled").(bool)
			source.Enabled = enabled
		}

		_, err = c.UpdateSource(*source.ID, source.Slug, source.Enabled, source.Name)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceSourceRead(ctx, d, m)
}

func resourceSourceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*segment.Client)

	var diags diag.Diagnostics

	sourceID := d.Id()

	_, err := c.DeleteSource(sourceID)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
