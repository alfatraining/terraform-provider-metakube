package metakube

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"gitlab.com/furkhat/terraform-provider-metakube/gometakube"
)

const (
	accProviderDCEnvname       = "ACC_PROVIDER_DC"
	accProviderUsernameEnvname = "ACC_PROVIDER_USERNAME"
	accProviderPasswordEnvname = "ACC_PROVIDER_PASSWORD"
)

func testAccCheckMetakubeClusterConfig(project, dc, username, password string) string {
	return fmt.Sprintf(`
provider "metakube" {

}

resource "metakube_project" "cluster-project" {
	name = "%s"

	labels = {}
}

resource "metakube_cluster" "bar" {
	project_id = metakube_project.cluster-project.id
	name = "bar"
	version = "1.17.3"
	dc = "%s"
	tenant = "syseleveneigenbedarf-syseleven-%s"
	provider_username = "%s"
	provider_password = "%s"

	nodepool {
		name = "my-nodepool"
		replicas = 2

		flavor_type = "Local Storage"
		flavor = "l1.small"
		image = "Rescue Ubuntu 16.04 sys11"
		use_floating_ip = false
	}

}
`, project, dc, username, username, password)
}

func TestAccMetakubeCluster_Basic(t *testing.T) {
	testDC := os.Getenv(accProviderDCEnvname)
	testProviderUsername := os.Getenv(accProviderUsernameEnvname)
	testProviderPassword := os.Getenv(accProviderPasswordEnvname)
	config := testAccCheckMetakubeClusterConfig(
		acctest.RandString(8),
		testDC,
		testProviderUsername,
		testProviderPassword)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testEnvSet(t, "METAKUBE_API_TOKEN")
			testEnvSet(t, accProviderDCEnvname)
			testEnvSet(t, accProviderUsernameEnvname)
			testEnvSet(t, accProviderPasswordEnvname)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetakubeClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterResourceCreated("metakube_cluster.bar"),
					resource.TestCheckResourceAttr("metakube_cluster.bar", "name", "bar"),
					resource.TestCheckResourceAttr("metakube_cluster.bar", "version", "1.17.3"),
					resource.TestCheckResourceAttr("metakube_cluster.bar", "dc", testDC),
					resource.TestCheckResourceAttr("metakube_cluster.bar", "tenant", "syseleveneigenbedarf-syseleven-"+testProviderUsername),
					resource.TestCheckResourceAttr("metakube_cluster.bar", "provider_username", testProviderUsername),
					resource.TestCheckResourceAttr("metakube_cluster.bar", "provider_password", testProviderPassword),
					resource.TestCheckResourceAttr("metakube_cluster.bar", "nodepool.#", "1"),
					resource.TestCheckResourceAttr("metakube_cluster.bar", "nodepool.0.name", "my-nodepool"),
					resource.TestCheckResourceAttr("metakube_cluster.bar", "nodepool.0.replicas", "2"),
					resource.TestCheckResourceAttr("metakube_cluster.bar", "nodepool.0.flavor_type", "Local Storage"),
					resource.TestCheckResourceAttr("metakube_cluster.bar", "nodepool.0.flavor", "l1.small"),
					resource.TestCheckResourceAttr("metakube_cluster.bar", "nodepool.0.use_floating_ip", "false"),
				),
			},
		},
	})
}

func testAccCheckMetakubeClusterDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*gometakube.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "metakube_cluster" {
			continue
		}
		projectID := rs.Primary.Attributes["project_id"]
		dc := rs.Primary.Attributes["dc"]
		obj, err := client.Clusters.Get(context.Background(), projectID, dc, rs.Primary.ID)
		if err == nil && obj == nil {
			return nil
		}

		if obj.DeletionTimestamp == nil {
			return fmt.Errorf("found not deleted cluster, project: %s, dc: %s, id: %s",
				projectID, dc, rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckClusterResourceCreated(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("Not found: %s", r)
		}

		client := testAccProvider.Meta().(*gometakube.Client)
		projectID := rs.Primary.Attributes["project_id"]
		dcName := rs.Primary.Attributes["dc"]
		dc, err := client.Datacenters.Get(context.Background(), dcName)
		if err != nil {
			return fmt.Errorf("failed to get datacenter details: %v", err)
		}
		obj, err := client.Clusters.Get(context.Background(), projectID, dc.Spec.Seed, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("failed to get cluster details: %v", err)
		}
		if obj == nil {
			return fmt.Errorf("cluster not created")
		}

		return nil
	}
}
