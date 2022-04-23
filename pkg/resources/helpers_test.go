package resources_test

import (
	"os"
	"testing"

	"github.com/gthesheep/terraform-provider-segment/pkg/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func providers() map[string]*schema.Provider {
	p := provider.Provider()
	return map[string]*schema.Provider{
		"segment": p,
	}
}

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = provider.Provider()
	testAccProviders = map[string]*schema.Provider{
		"segment": testAccProvider,
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("SEGMENT_API_TOKEN"); v == "" {
		t.Fatal("SEGMENT_API_TOKEN must be set for acceptance tests")
	}
}
