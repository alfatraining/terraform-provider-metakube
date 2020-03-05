package metakube

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var (
	testAccProviders map[string]terraform.ResourceProvider
	testAccProvider  *schema.Provider
)

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]terraform.ResourceProvider{
		"metakube": testAccProvider,
	}
}

func testEnvSet(t *testing.T, e string) {
	if os.Getenv(e) == "" {
		t.Fatalf("%s must be set for acceptance tests", e)
	}
}
