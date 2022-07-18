// package resources
//
// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
//
// 	"github.com/gthesheep/terraform-provider-segment/pkg/segment"
// 	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
// 	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
// 	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
// )
//
//
// func ResourceSource() *schema.Resource {
// 	return &schema.Resource{
// 		CreateContext: resourceUserCreate,
// 		ReadContext:   resourceUserRead,
// 		UpdateContext: resourceUserUpdate,
// 		DeleteContext: resourceUserDelete,
//
// 		Schema: map[string]*schema.Schema{
// 			"email": &schema.Schema{
// 				Type:        schema.TypeString,
// 				Required:    true,
// 				Description: "Email address of the user to invite",
// 			},
// 			"role_name": &schema.Schema{
// 				Type:        schema.TypeString,
// 				Required:    true,
// 				Description: "User's assigned role",
// 			},
// 		},
//
// 		Importer: &schema.ResourceImporter{
// 			StateContext: schema.ImportStatePassthroughContext,
// 		},
// 	}
// }
//
// func resourceSourceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
// 	c := m.(*segment.Client)
//
// 	var diags diag.Diagnostics
//
// 	email := d.Get("email").(string)
// 	roleName := d.Get("role_name").(string)
//
//     role, err := c.GetRole(roleName)
//     if err != nil {
// 		return diag.FromErr(err)
// 	}
//
// 	invitedEmail, err := c.CreateInviteLink(email, role.ID)
// 	if err != nil {
// 		return diag.FromErr(err)
// 	}
// 	d.SetId(invitedEmail)
//
// 	resourceUserRead(ctx, d, m)
//
// 	return diags
// }
//
//
// func resourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
// 	c := m.(*segment.Client)
//
// 	var diags diag.Diagnostics
//
// 	userID := d.Id()
//
// 	source, err := c.GetSource(sourceID)
// 	if err != nil {
// 		return diag.FromErr(err)
// 	}
//
// 	if err := d.Set("slug", source.Slug); err != nil {
// 		return diag.FromErr(err)
// 	}
// 	if err := d.Set("enabled", source.Enabled); err != nil {
// 		return diag.FromErr(err)
// 	}
// 	if err := d.Set("name", source.Name); err != nil {
// 		return diag.FromErr(err)
// 	}
// 	if err := d.Set("source_slug", source.Metadata.Slug); err != nil {
// 		return diag.FromErr(err)
// 	}
//
// 	s := flattenSourceSettings(source.Settings)
// 	if err := d.Set("settings", s); err != nil {
// 		return diag.FromErr(err)
// 	}
//
// 	return diags
// }
//
// func resourceSourceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
// 	c := m.(*segment.Client)
//
// 	sourceID := d.Id()
//
// 	if d.HasChange("name") || d.HasChange("enabled") || d.HasChange("settings") {
// 		source, err := c.GetSource(sourceID)
// 		if err != nil {
// 			return diag.FromErr(err)
// 		}
//
// 		if d.HasChange("name") {
// 			name := d.Get("name").(string)
// 			source.Name = name
// 		}
// 		if d.HasChange("enabled") {
// 			enabled := d.Get("enabled").(bool)
// 			source.Enabled = enabled
// 		}
// 		if d.HasChange("settings") {
// 			settings := d.Get("settings").([]interface{})
// 			sourceSettings := mapToSourceSettings(settings)
// 			source.Settings = sourceSettings
// 		}
//
// 		_, err = c.UpdateSource(*source.ID, source.Slug, source.Enabled, source.Name, source.Settings)
// 		if err != nil {
// 			return diag.FromErr(err)
// 		}
// 	}
//
// 	return resourceSourceRead(ctx, d, m)
// }
//
// func resourceSourceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
// 	c := m.(*segment.Client)
//
// 	var diags diag.Diagnostics
//
// 	sourceID := d.Id()
//
// 	_, err := c.DeleteSource(sourceID)
// 	if err != nil {
// 		return diag.FromErr(err)
// 	}
//
// 	return diags
// }
