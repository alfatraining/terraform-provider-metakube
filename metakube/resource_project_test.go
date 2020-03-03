package metakube

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"gitlab.com/furkhat/terraform-provider-metakube/gometakube"
)

const (
	testAccCheckMetakubeProjectConfig = `
provider "metakube" {

}

resource "metakube_project" "foo" {
	name = "foo name"

	labels = {
		additionalProp1 = "string"
		additionalProp2 = "string"
		additionalProp3 = "string"
	}
}
`
)

func TestAccMetakubeProject_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetakubeProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckMetakubeProjectConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckProjectResourceExist("metakube_project.foo", "foo name", map[string]string{
						"additionalProp1": "string",
						"additionalProp2": "string",
						"additionalProp3": "string",
					}),
					resource.TestCheckResourceAttr("metakube_project.foo", "name", "foo name"),
				),
			},
		},
	})
}

func testAccCheckMetakubeProjectDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*gometakube.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "metakube_project" {
			continue
		}

		obj, err := client.Projects.Get(context.Background(), rs.Primary.ID)
		if err == nil && obj == nil {
			return nil
		}

		if obj.Status != "Terminating" {
			return fmt.Errorf("found not deleted project in `%s` status", obj.Status)
		}
	}
	return nil
}

func testAccCheckProjectResourceExist(r, name string, labels map[string]string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("Not found: %s", r)
		}

		client := testAccProvider.Meta().(*gometakube.Client)

		if obj, err := client.Projects.Get(context.Background(), rs.Primary.ID); err != nil {
			return fmt.Errorf("Couldnt retrieve project: %v", err)
		} else if want, got := name, obj.Name; want != got {
			return fmt.Errorf("Unexpected project name, want: %s, got: %s", want, got)
		} else if want, got := labels, obj.Labels; !reflect.DeepEqual(want, got) {
			return fmt.Errorf("Unexpected project labels, want: %v, got: %v", want, got)
		}
		return nil
	}
}
