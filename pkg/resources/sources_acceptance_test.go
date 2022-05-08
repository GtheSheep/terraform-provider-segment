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

func TestAccSegmentSourceResource(t *testing.T) {

	slug := strings.ToLower(acctest.RandStringFromCharSet(4, acctest.CharSetAlpha))
	name := strings.ToLower(acctest.RandStringFromCharSet(3, acctest.CharSetAlpha))
	name2 := strings.ToLower(acctest.RandStringFromCharSet(3, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSegmentSourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSegmentSourceResourceBasicConfig(slug, name, "facebook-ads"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSegmentSourceExists("segment_source.test_source"),
					resource.TestCheckResourceAttr("segment_source.test_source", "slug", slug),
					resource.TestCheckResourceAttr("segment_source.test_source", "name", name),
					resource.TestCheckResourceAttr("segment_source.test_source", "source_slug", "facebook-ads"),
				),
			},
			// RENAME
			{
				Config: testAccSegmentSourceResourceBasicConfig(slug, name2, "facebook-ads"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSegmentSourceExists("segment_source.test_source"),
					resource.TestCheckResourceAttr("segment_source.test_source", "slug", slug),
					resource.TestCheckResourceAttr("segment_source.test_source", "name", name2),
					resource.TestCheckResourceAttr("segment_source.test_source", "source_slug", "facebook-ads"),
				),
			},
			// FULL CONFIG
			{
				Config: testAccSegmentSourceResourceFullConfig(slug, name, "facebook-ads"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSegmentSourceExists("segment_source.test_source"),
					resource.TestCheckResourceAttr("segment_source.test_source", "slug", slug),
					resource.TestCheckResourceAttr("segment_source.test_source", "name", name),
					resource.TestCheckResourceAttr("segment_source.test_source", "source_slug", "facebook-ads"),
				),
			},
			// IMPORT
			{
				ResourceName:            "segment_source.test_source",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func testAccSegmentSourceResourceBasicConfig(slug, name, sourceSlug string) string {
	return fmt.Sprintf(`
resource "segment_source" "test_source" {
  slug        = "%s"
  name        = "%s"
  source_slug = "%s"
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
`, slug, name, sourceSlug)
}

func testAccSegmentSourceResourceFullConfig(slug, name, sourceSlug string) string {
	return fmt.Sprintf(`
resource "segment_source" "test_source" {
  slug        = "%s"
  name        = "%s"
  source_slug = "%s"
  enabled     = false
  settings {
    track {
      allow_unplanned_events = false
    }
    identify {
      allow_unplanned_traits = false
    }
    group {
      allow_unplanned_traits = true
      allow_traits_on_violations = false
    }
  }
}
`, slug, name, sourceSlug)
}

func testAccCheckSegmentSourceExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}
		apiClient := testAccProvider.Meta().(*segment.Client)
		sourceID := rs.Primary.ID

		_, err := apiClient.GetSource(sourceID)
		if err != nil {
			return fmt.Errorf("error fetching item with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckSegmentSourceDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*segment.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "segment_source" {
			continue
		}
		sourceID := rs.Primary.ID

		_, err := apiClient.GetSource(sourceID)
		if err == nil {
			return fmt.Errorf("Source still exists")
		}
		notFoundErr := "not found"
		expectedErr := regexp.MustCompile(notFoundErr)
		if !expectedErr.Match([]byte(err.Error())) {
			return fmt.Errorf("expected %s, got %s", notFoundErr, err)
		}
	}

	return nil
}
