package metakube

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/pkg/errors"
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

func TestAccMetakubeProject_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testEnvSet(t, APITokenEnvName) },
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

		obj, resp, err := client.Projects.Get(context.Background(), rs.Primary.ID)
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil
		}
		if err != nil {
			return err
		}

		if obj.Status != "Terminating" {
			return errors.Errorf("found not deleted project in `%s` status", obj.Status)
		}
	}
	return nil
}

func testAccCheckProjectResourceExist(r, name string, labels map[string]string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return errors.Errorf("not found %s", r)
		}

		client := testAccProvider.Meta().(*gometakube.Client)

		if obj, _, err := client.Projects.Get(context.Background(), rs.Primary.ID); err != nil {
			return errors.Wrap(err, "get project")
		} else if want, got := name, obj.Name; want != got {
			return errors.Errorf("want Name=%v, got %v", want, got)
		} else if want, got := labels, obj.Labels; !reflect.DeepEqual(want, got) {
			return errors.Errorf("want Labels=%v, got %v", want, got)
		}
		return nil
	}
}
