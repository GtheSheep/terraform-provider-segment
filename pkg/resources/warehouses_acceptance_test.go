package resources_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/gthesheep/terraform-provider-segment/pkg/segment"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSegmentWarehouseResource(t *testing.T) {

	name := strings.ToLower(acctest.RandStringFromCharSet(3, acctest.CharSetAlpha))
	name2 := strings.ToLower(acctest.RandStringFromCharSet(3, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSegmentWarehouseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSegmentWarehouseResourceBasicConfig(name, "snowflake"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSegmentWarehouseExists("segment_warehouse.test_warehouse"),
					resource.TestCheckResourceAttr("segment_warehouse.test_warehouse", "name", name),
					resource.TestCheckResourceAttr("segment_warehouse.test_warehouse", "warehouse_slug", "snowflake"),
					resource.TestCheckResourceAttr("segment_warehouse.test_warehouse", "settings.0.username", "Moo"),
					resource.TestCheckResourceAttr("segment_warehouse.test_warehouse", "settings.0.password", "Password"),
					resource.TestCheckResourceAttr("segment_warehouse.test_warehouse", "settings.0.port", "22"),
					resource.TestCheckResourceAttr("segment_warehouse.test_warehouse", "settings.0.hostname", "www.snowflake.com"),
				),
			},
			// RENAME
			{
				Config: testAccSegmentWarehouseResourceBasicConfig(name2, "snowflake"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSegmentWarehouseExists("segment_warehouse.test_warehouse"),
					resource.TestCheckResourceAttr("segment_warehouse.test_warehouse", "name", name2),
					resource.TestCheckResourceAttr("segment_warehouse.test_warehouse", "warehouse_slug", "snowflake"),
					resource.TestCheckResourceAttr("segment_warehouse.test_warehouse", "settings.0.username", "Moo"),
					resource.TestCheckResourceAttr("segment_warehouse.test_warehouse", "settings.0.password", "Password"),
					resource.TestCheckResourceAttr("segment_warehouse.test_warehouse", "settings.0.port", "22"),
					resource.TestCheckResourceAttr("segment_warehouse.test_warehouse", "settings.0.hostname", "www.snowflake.com"),
				),
			},
			// IMPORT
			{
				ResourceName:            "segment_warehouse.test_warehouse",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"settings.0.password"},
			},
		},
	})
}

func testAccSegmentWarehouseResourceBasicConfig(name, warehouseSlug string) string {
	return fmt.Sprintf(`
resource "segment_warehouse" "test_warehouse" {
  name        = "%s"
  warehouse_slug = "%s"
  enabled     = false
  settings {
    username = "Moo"
    password = "Password"
    port = 22
    hostname = "www.snowflake.com"
  }
}
`, name, warehouseSlug)
}

func testAccCheckSegmentWarehouseExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}
		apiClient := testAccProvider.Meta().(*segment.Client)
		warehouseID := rs.Primary.ID

		_, err := apiClient.GetWarehouse(warehouseID)
		if err != nil {
			return fmt.Errorf("error fetching item with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckSegmentWarehouseDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*segment.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "segment_warehouse" {
			continue
		}
		warehouseID := rs.Primary.ID

		_, err := apiClient.GetWarehouse(warehouseID)
		if err == nil {
			return fmt.Errorf("Warehouse still exists")
		}
		notFoundErr := "not found"
		expectedErr := regexp.MustCompile(notFoundErr)
		if !expectedErr.Match([]byte(err.Error())) {
			return fmt.Errorf("expected %s, got %s", notFoundErr, err)
		}
	}

	return nil
}
