package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gthesheep/terraform-provider-segment/pkg/segment"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceWarehouse() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWarehouseCreate,
		ReadContext:   resourceWarehouseRead,
		UpdateContext: resourceWarehouseUpdate,
		DeleteContext: resourceWarehouseDelete,

		Schema: map[string]*schema.Schema{
			"enabled": &schema.Schema{
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Flag for whether or not the warehouse is enabled",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Descriptive name for the warehouse",
			},
			"warehouse_slug": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Slug for the warehouse, from a list",
				ValidateFunc: validation.StringInSlice(segment.WarehouseSlugs, false),
			},
			"settings": &schema.Schema{
				Type:        schema.TypeList,
				Required:    true,
				Description: "Map containing settings for the warehouse",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hostname": {
							Description: "Hostname for the warehouse",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"database": {
							Description: "Database for the warehouse",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"port": {
							Description: "Port for the warehouse",
							Type:        schema.TypeInt,
							Optional:    true,
						},
						"username": {
							Description: "Username for the warehouse",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"password": {
							Description: "Password for the warehouse, can't be imported, and beware of storing in the state file!",
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
						},
						"ciphertext": {
							Description: "Ciphertext for the warehouse",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"name": {
							Description: "Warehouse name gets stored in the settings after creation",
							Type:        schema.TypeString,
							Computed:    true,
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

func resourceWarehouseCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*segment.Client)

	var diags diag.Diagnostics

	enabled := d.Get("enabled").(bool)
	name := d.Get("name").(string)
	warehouseSlug := d.Get("warehouse_slug").(string)
	settings := d.Get("settings").([]interface{})
	warehouseSettings := mapToWarehouseSettings(settings)

	warehouse, err := c.CreateWarehouse(enabled, name, warehouseSlug, warehouseSettings)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprintf("%s", *warehouse.ID))

	resourceWarehouseRead(ctx, d, m)

	return diags
}

func mapToWarehouseSettings(settings []interface{}) segment.WarehouseSettings {
	if len(settings) == 0 {
		return segment.WarehouseSettings{}
	}
	actualSettings := settings[0].(map[string]interface{})
	port := actualSettings["port"].(int)
	portString := strconv.Itoa(port)
	return segment.WarehouseSettings{
		Hostname: actualSettings["hostname"].(string),
		Database: actualSettings["database"].(string),
		Port:     portString,
		Username: actualSettings["username"].(string),
		Password: actualSettings["password"].(string),
	}
	return segment.WarehouseSettings{}
}

func flattenWarehouseSettings(warehouseSettings segment.WarehouseSettings) []interface{} {
	var settings map[string]interface{}
	settingsJson, _ := json.Marshal(warehouseSettings)
	json.Unmarshal(settingsJson, &settings)
	portInt, _ := strconv.Atoi(warehouseSettings.Port)
	settings["port"] = portInt
	flatSettings := make([]interface{}, 1)
	flatSettings[0] = settings
	return flatSettings
}

func resourceWarehouseRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*segment.Client)

	var diags diag.Diagnostics

	warehouseID := d.Id()

	warehouse, err := c.GetWarehouse(warehouseID)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("enabled", warehouse.Enabled); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", warehouse.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("warehouse_slug", warehouse.Metadata.Slug); err != nil {
		return diag.FromErr(err)
	}

	if warehouse.Settings.Password == "" {
		existingSettings := d.Get("settings").([]interface{})
		existingWarehouseSettings := mapToWarehouseSettings(existingSettings)
		warehouse.Settings.Password = existingWarehouseSettings.Password
	}
	s := flattenWarehouseSettings(warehouse.Settings)
	if err := d.Set("settings", s); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceWarehouseUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*segment.Client)

	warehouseID := d.Id()

	if d.HasChange("name") || d.HasChange("enabled") || d.HasChange("settings") {
		warehouse, err := c.GetWarehouse(warehouseID)
		if err != nil {
			return diag.FromErr(err)
		}

		if d.HasChange("name") {
			name := d.Get("name").(string)
			warehouse.Name = name
		}
		if d.HasChange("enabled") {
			enabled := d.Get("enabled").(bool)
			warehouse.Enabled = enabled
		}
		if d.HasChange("settings") {
			settings := d.Get("settings").([]interface{})
			warehouseSettings := mapToWarehouseSettings(settings)
			warehouse.Settings = warehouseSettings
		}

		_, err = c.UpdateWarehouse(*warehouse.ID, warehouse.Enabled, warehouse.Name, warehouse.Settings)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceWarehouseRead(ctx, d, m)
}

func resourceWarehouseDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*segment.Client)

	var diags diag.Diagnostics

	warehouseID := d.Id()

	_, err := c.DeleteWarehouse(warehouseID)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
