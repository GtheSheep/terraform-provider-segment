package resources

import (
	"context"
	"encoding/json"
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
				Required:    true,
				Description: "Map containing settings for the source",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"forwarding_violations_to": {
							Description: "SourceId to forward violations to.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"forwarding_blocked_events_to": {
							Description: "SourceId to forward blocked events to.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"track": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"allow_unplanned_events": {
										Description: "Enable to allow unplanned track events.",
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
									},
									"allow_unplanned_event_properties": {
										Description: "Enable to allow unplanned track event properties.",
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
									},
									"allow_event_on_violations": {
										Description: "Allow track event on violations.",
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
									},
									"allow_properties_on_violations": {
										Description: "Enable to allow track properties on violations.",
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
									},
									"common_event_on_violations": {
										Description:  "The common track event on violations.",
										Type:         schema.TypeString,
										Optional:     true,
										Default:      "",
										ValidateFunc: validation.StringInSlice(ViolationEvents, false),
									},
								},
							},
						},
						"identify": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"allow_unplanned_traits": {
										Description: "Enable to allow unplanned identify traits.",
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
									},
									"allow_traits_on_violations": {
										Description: "Enable to allow identify traits on violations.",
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
									},
									"common_event_on_violations": {
										Description:  "The common track event on violations.",
										Type:         schema.TypeString,
										Optional:     true,
										Default:      "",
										ValidateFunc: validation.StringInSlice(ViolationEvents, false),
									},
								},
							},
						},
						"group": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"allow_unplanned_traits": {
										Description: "Enable to allow unplanned identify traits.",
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
									},
									"allow_traits_on_violations": {
										Description: "Enable to allow identify traits on violations.",
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     false,
									},
									"common_event_on_violations": {
										Description:  "The common track event on violations.",
										Type:         schema.TypeString,
										Default:      "",
										Optional:     true,
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
	settings := d.Get("settings").([]interface{})
	sourceSettings := mapToSourceSettings(settings)

	source, err := c.CreateSource(slug, enabled, name, sourceSlug, sourceSettings)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprintf("%s", *source.ID))

	resourceSourceRead(ctx, d, m)

	return diags
}

func mapToSourceSettings(settings []interface{}) segment.SourceSettings {
	if len(settings) == 0 {
		return segment.SourceSettings{}
	}
	actualSettings := settings[0].(map[string]interface{})

	trackSettingsList := actualSettings["track"].([]interface{})
	trackSettings := make(map[string]interface{})
	if len(trackSettingsList) > 0 {
		trackSettings = trackSettingsList[0].(map[string]interface{})
	}

	identifySettingsList := actualSettings["identify"].([]interface{})
	identifySettings := make(map[string]interface{})
	if len(identifySettingsList) > 0 {
		identifySettings = identifySettingsList[0].(map[string]interface{})
	}

	groupSettingsList := actualSettings["group"].([]interface{})
	groupSettings := make(map[string]interface{})
	if len(groupSettingsList) > 0 {
		groupSettings = groupSettingsList[0].(map[string]interface{})
	}
	return segment.SourceSettings{
		ForwardingViolationsTo:    actualSettings["forwarding_violations_to"].(string),
		ForwardingBlockedEventsTo: actualSettings["forwarding_blocked_events_to"].(string),
		Track: segment.TrackingSettings{
			AllowUnplannedEvents:          trackSettings["allow_unplanned_events"].(bool),
			AllowUnplannedEventProperties: trackSettings["allow_unplanned_event_properties"].(bool),
			AllowEventOnViolations:        trackSettings["allow_event_on_violations"].(bool),
			AllowPropertiesOnViolations:   trackSettings["allow_properties_on_violations"].(bool),
			CommonEventOnViolations:       trackSettings["common_event_on_violations"].(string),
		},
		Identify: segment.IdentifySettings{
			AllowUnplannedTraits:    identifySettings["allow_unplanned_traits"].(bool),
			AllowTraitsOnViolations: identifySettings["allow_traits_on_violations"].(bool),
			CommonEventOnViolations: identifySettings["common_event_on_violations"].(string),
		},
		Group: segment.GroupSettings{
			AllowUnplannedTraits:    groupSettings["allow_unplanned_traits"].(bool),
			AllowTraitsOnViolations: groupSettings["allow_traits_on_violations"].(bool),
			CommonEventOnViolations: groupSettings["common_event_on_violations"].(string),
		},
	}
	return segment.SourceSettings{}
}

func flattenSourceSettings(sourceSettings segment.SourceSettings) []interface{} {
	var settings map[string]interface{}
	settingsJson, _ := json.Marshal(sourceSettings)
	json.Unmarshal(settingsJson, &settings)

	flatSettings := make([]interface{}, 1)

	tracking := make([]interface{}, 1)
	tracking[0] = settings["track"]
	settings["track"] = tracking

	grouping := make([]interface{}, 1)
	grouping[0] = settings["group"]
	settings["group"] = grouping

	identify := make([]interface{}, 1)
	identify[0] = settings["identify"]
	settings["identify"] = identify

	flatSettings[0] = settings
	return flatSettings
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

	s := flattenSourceSettings(source.Settings)
	if err := d.Set("settings", s); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceSourceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*segment.Client)

	sourceID := d.Id()

	if d.HasChange("name") || d.HasChange("enabled") || d.HasChange("settings") {
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
		if d.HasChange("settings") {
			settings := d.Get("settings").([]interface{})
			sourceSettings := mapToSourceSettings(settings)
			source.Settings = sourceSettings
		}

		_, err = c.UpdateSource(*source.ID, source.Slug, source.Enabled, source.Name, source.Settings)
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
