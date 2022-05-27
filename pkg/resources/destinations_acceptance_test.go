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

func TestAccSegmentDestinationResource(t *testing.T) {

	name := strings.ToLower(acctest.RandStringFromCharSet(3, acctest.CharSetAlpha))
	sourceSlug := strings.ToLower(acctest.RandStringFromCharSet(3, acctest.CharSetAlpha))
	sourceName := strings.ToLower(acctest.RandStringFromCharSet(3, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSegmentDestinationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSegmentDestinationResourceBasicConfig(sourceSlug, sourceName, "google-tag-manager", name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSegmentDestinationExists("segment_destination.test_destination"),
					resource.TestCheckResourceAttr("segment_destination.test_destination", "destination_slug", "google-tag-manager"),
					resource.TestCheckResourceAttr("segment_destination.test_destination", "name", name),
				),
			},
			// IMPORT
			{
				ResourceName:            "segment_destination.test_destination",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func testAccSegmentDestinationResourceBasicConfig(sourceSlug, sourceName, destinationSlug, name string) string {
	return fmt.Sprintf(`
resource "segment_source" "test_source" {
  slug        = "%s"
  name        = "%s"
  source_slug = "facebook-ads"
  enabled     = false
  settings {
    track {
    }
    identify {
    }
    group {
    }
  }
}

resource "segment_destination" "test_destination" {
  name        = "%s"
  destination_slug = "%s"
  enabled     = false
  source_id   = segment_source.test_source.id
  settings = {
    containerId = "xxxx"
    environment = "gtm_auth=xxxx"
    trackAllPages = "false"
    trackCategorizedPages = "false"
    trackNamedPages = "false"
  }
}
`, sourceSlug, sourceName, name, destinationSlug)
}

//
// func testAccSegmentSourceResourceFullConfig(slug, name, sourceSlug string) string {
// 	return fmt.Sprintf(`
// resource "segment_source" "test_source" {
//   slug        = "%s"
//   name        = "%s"
//   source_slug = "%s"
//   enabled     = false
//   settings {
//     track {
//       allow_unplanned_events = false
//     }
//     identify {
//       allow_unplanned_traits = false
//     }
//     group {
//       allow_unplanned_traits = true
//       allow_traits_on_violations = false
//     }
//   }
// }
// `, slug, name, sourceSlug)
// }

func testAccCheckSegmentDestinationExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}
		apiClient := testAccProvider.Meta().(*segment.Client)
		destinationID := rs.Primary.ID

		_, err := apiClient.GetDestination(destinationID)
		if err != nil {
			return fmt.Errorf("error fetching item with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckSegmentDestinationDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*segment.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "segment_destination" {
			continue
		}
		destinationID := rs.Primary.ID

		_, err := apiClient.GetDestination(destinationID)
		if err == nil {
			return fmt.Errorf("Destination still exists")
		}
		notFoundErr := "not found"
		expectedErr := regexp.MustCompile(notFoundErr)
		if !expectedErr.Match([]byte(err.Error())) {
			return fmt.Errorf("expected %s, got %s", notFoundErr, err)
		}
	}

	return nil
}
