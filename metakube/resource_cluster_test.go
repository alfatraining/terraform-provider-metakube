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
	accTenantEnvname           = "ACC_TENANT"
	accProviderUsernameEnvname = "ACC_PROVIDER_USERNAME"
	accProviderPasswordEnvname = "ACC_PROVIDER_PASSWORD"
)

func testAccMetakubeClusterConfig(project, dc, tenant, username, password string) string {
	return fmt.Sprintf(`
provider "metakube" {

}

resource "metakube_project" "cluster-project" {
	name = "%s"

	labels = {}
}

resource "metakube_cluster" "bar" {
	project_id = metakube_project.cluster-project.id
	name = "my-cluster"
	labels = {
		"version" = "alpha"
	}
	version = "1.17.3"
	dc = "%s"
	tenant = "%s"
	provider_username = "%s"
	provider_password = "%s"
	audit_logging = true

	nodedepl {
		name = "my-nodedepl"
		replicas = 2

		flavor = "l1.small"
		image = "Rescue Ubuntu 16.04 sys11"
		use_floating_ip = false
	}
}
`, project, dc, tenant, username, password)
}

func testAccMetakubeClusterConfigUpdate(project, dc, tenant, username, password string) string {
	return fmt.Sprintf(`
provider "metakube" {

}

resource "metakube_project" "cluster-project" {
	name = "%s"

	labels = {}
}

resource "metakube_cluster" "bar" {
	project_id = metakube_project.cluster-project.id
	name = "my-cluster-edit"
	labels = {
		"version" = "beta"
	}
	version = "1.17.3"
	dc = "%s"
	tenant = "%s"
	provider_username = "%s"
	provider_password = "%s"
	audit_logging = false

	nodedepl {
		name = "my-nodedepl"
		replicas = 1

		flavor = "m1c.medium"
		image = "Rescue Ubuntu 18.04 sys11"
		use_floating_ip = true
	}
}
`, project, dc, tenant, username, password)
}

func TestAccMetakubeCluster_CreateAndInPlaceUpdates(t *testing.T) {
	testDC := os.Getenv(accProviderDCEnvname)
	testTenant := os.Getenv(accTenantEnvname)
	testProviderUsername := os.Getenv(accProviderUsernameEnvname)
	testProviderPassword := os.Getenv(accProviderPasswordEnvname)
	projectName := acctest.RandString(8)
	config := testAccMetakubeClusterConfig(
		projectName,
		testDC,
		testTenant,
		testProviderUsername,
		testProviderPassword)
	configUpdated := testAccMetakubeClusterConfigUpdate(
		projectName,
		testDC,
		testTenant,
		testProviderUsername,
		testProviderPassword,
	)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testEnvSet(t, "METAKUBE_API_TOKEN")
			testEnvSet(t, accProviderDCEnvname)
			testEnvSet(t, accTenantEnvname)
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
					testAccCheckClustersNodeDeployment("metakube_cluster.bar", "my-nodedepl", "l1.small", "Rescue Ubuntu 16.04 sys11", false, 2),
					resource.TestCheckResourceAttr("metakube_cluster.bar", "name", "my-cluster"),
					resource.TestCheckResourceAttr("metakube_cluster.bar", "labels.version", "alpha"),
					resource.TestCheckResourceAttr("metakube_cluster.bar", "version", "1.17.3"),
					resource.TestCheckResourceAttr("metakube_cluster.bar", "dc", testDC),
					resource.TestCheckResourceAttr("metakube_cluster.bar", "audit_logging", "true"),
					resource.TestCheckResourceAttr("metakube_cluster.bar", "provider_username", testProviderUsername),
					resource.TestCheckResourceAttr("metakube_cluster.bar", "provider_password", testProviderPassword),
					resource.TestCheckResourceAttr("metakube_cluster.bar", "nodedepl.#", "1"),
					resource.TestCheckResourceAttr("metakube_cluster.bar", "nodedepl.0.name", "my-nodedepl"),
					resource.TestCheckResourceAttr("metakube_cluster.bar", "nodedepl.0.replicas", "2"),
					resource.TestCheckResourceAttr("metakube_cluster.bar", "nodedepl.0.flavor", "l1.small"),
					resource.TestCheckResourceAttr("metakube_cluster.bar", "nodedepl.0.image", "Rescue Ubuntu 16.04 sys11"),
					resource.TestCheckResourceAttr("metakube_cluster.bar", "nodedepl.0.use_floating_ip", "false"),
				),
			},
			{
				Config: configUpdated,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterResourceCreated("metakube_cluster.bar"),
					testAccCheckClustersNodeDeployment("metakube_cluster.bar", "my-nodedepl", "m1c.medium", "Rescue Ubuntu 18.04 sys11", true, 1),
					resource.TestCheckResourceAttr("metakube_cluster.bar", "name", "my-cluster-edit"),
					resource.TestCheckResourceAttr("metakube_cluster.bar", "labels.version", "beta"),
					resource.TestCheckResourceAttr("metakube_cluster.bar", "audit_logging", "false"),
					resource.TestCheckResourceAttr("metakube_cluster.bar", "nodedepl.0.replicas", "1"),
					resource.TestCheckResourceAttr("metakube_cluster.bar", "nodedepl.0.flavor", "m1c.medium"),
					resource.TestCheckResourceAttr("metakube_cluster.bar", "nodedepl.0.image", "Rescue Ubuntu 18.04 sys11"),
					resource.TestCheckResourceAttr("metakube_cluster.bar", "nodedepl.0.use_floating_ip", "true"),
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

func testAccCheckClustersNodeDeployment(r, name, flavor, image string, floatingIP bool, replicas uint) resource.TestCheckFunc {
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
		var nodedepl *gometakube.NodeDeployment
		items, err := client.NodeDeployments.List(context.Background(), projectID, dc.Spec.Seed, rs.Primary.ID)
		for _, item := range items {
			if item.Name == name {
				nodedepl = &item
				break
			}
		}
		if nodedepl == nil {
			return fmt.Errorf("Not found node deployment with name: %s", name)
		}
		if nodedepl.Spec.Replicas != replicas {
			return fmt.Errorf("want replicas: %v, got: %v", replicas, nodedepl.Spec.Replicas)
		}
		if want, got := flavor, nodedepl.Spec.Template.Cloud.Openstack.Flavor; want != got {
			return fmt.Errorf("want flavor: %v, got: %v", want, got)
		}
		if want, got := image, nodedepl.Spec.Template.Cloud.Openstack.Image; want != got {
			return fmt.Errorf("want image: %v, got: %v", want, got)
		}
		if want, got := floatingIP, nodedepl.Spec.Template.Cloud.Openstack.UseFloatingIP; want != got {
			return fmt.Errorf("want use_floating_ip: %v, got: %v", want, got)
		}
		return nil
	}
}
