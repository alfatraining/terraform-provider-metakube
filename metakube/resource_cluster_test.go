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
	sshkeys = [
		"dev"
	]
	version = "1.15"
	dc = "%s"
	tenant = "%s"
	provider_username = "%s"
	provider_password = "%s"
	audit_logging = true

	nodedepl {
		name = "my-nodedepl"
		replicas = 2

		autoscale {
		  min_replicas = 1
		  max_replicas = 3
		}

		flavor = "l1.small"
		image = "Rescue Ubuntu 16.04 sys11"
		use_floating_ip = false
	}
}

resource "metakube_sshkey" "my-key" {
	project_id = metakube_project.cluster-project.id
	name = "dev"
	public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCut5oRyqeqYci3E9m6Z6mtxfqkiyb+xNFJM6+/sllhnMDX0vzrNj8PuIFfGkgtowKY//QWLgoB+RpvXqcD4bb4zPkLdXdJPtUf1eAoMh/qgyThUjBs3n7BXvXMDg1Wdj0gq/sTnPLvXsfrSVPjiZvWN4h0JdID2NLnwYuKIiltIn+IbUa6OnyFfOEpqb5XJ7H7LK1mUKTlQ/9CFROxSQf3YQrR9UdtASIeyIZL53WgYgU31Yqy7MQaY1y0fGmHsFwpCK6qFZj1DNruKl/IR1lLx/Bg3z9sDcoBnHKnzSzVels9EVlDOG6bW738ho269QAIrWQYBtznsvWKu5xZPuuj user@machine"
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
	sshkeys = []
	version = "1.17"
	dc = "%s"
	tenant = "%s"
	provider_username = "%s"
	provider_password = "%s"
	audit_logging = false

	nodedepl {
		name = "my-nodedepl"
		replicas = 1

		autoscale {
		  min_replicas = 1
		  max_replicas = 2
		}

		flavor = "m1c.medium"
		image = "Rescue Ubuntu 18.04 sys11"
		use_floating_ip = true
	}
}

resource "metakube_sshkey" "my-key" {
	project_id = metakube_project.cluster-project.id
	name = "dev"
	public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCut5oRyqeqYci3E9m6Z6mtxfqkiyb+xNFJM6+/sllhnMDX0vzrNj8PuIFfGkgtowKY//QWLgoB+RpvXqcD4bb4zPkLdXdJPtUf1eAoMh/qgyThUjBs3n7BXvXMDg1Wdj0gq/sTnPLvXsfrSVPjiZvWN4h0JdID2NLnwYuKIiltIn+IbUa6OnyFfOEpqb5XJ7H7LK1mUKTlQ/9CFROxSQf3YQrR9UdtASIeyIZL53WgYgU31Yqy7MQaY1y0fGmHsFwpCK6qFZj1DNruKl/IR1lLx/Bg3z9sDcoBnHKnzSzVels9EVlDOG6bW738ho269QAIrWQYBtznsvWKu5xZPuuj user@machine"
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
			testEnvSet(t, APITokenEnvName)
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
					testAccCheckClustersNodeDeployment("metakube_cluster.bar", "my-nodedepl", "l1.small", "Rescue Ubuntu 16.04 sys11", false, 2, 1, 3),
					resource.TestCheckResourceAttr("metakube_cluster.bar", "name", "my-cluster"),
					resource.TestCheckResourceAttr("metakube_cluster.bar", "labels.version", "alpha"),
					resource.TestCheckResourceAttr("metakube_cluster.bar", "version", "1.15"),
					resource.TestCheckResourceAttr("metakube_cluster.bar", "dc", testDC),
					resource.TestCheckResourceAttr("metakube_cluster.bar", "audit_logging", "true"),
					resource.TestCheckResourceAttr("metakube_cluster.bar", "provider_username", testProviderUsername),
					resource.TestCheckResourceAttr("metakube_cluster.bar", "provider_password", testProviderPassword),
					resource.TestCheckResourceAttr("metakube_cluster.bar", "nodedepl.#", "1"),
					resource.TestCheckResourceAttr("metakube_cluster.bar", "nodedepl.0.name", "my-nodedepl"),
					resource.TestCheckResourceAttr("metakube_cluster.bar", "nodedepl.0.replicas", "2"),
					resource.TestCheckResourceAttr("metakube_cluster.bar", "nodedepl.0.autoscale.0.min_replicas", "1"),
					resource.TestCheckResourceAttr("metakube_cluster.bar", "nodedepl.0.autoscale.0.max_replicas", "3"),
					resource.TestCheckResourceAttr("metakube_cluster.bar", "nodedepl.0.flavor", "l1.small"),
					resource.TestCheckResourceAttr("metakube_cluster.bar", "nodedepl.0.image", "Rescue Ubuntu 16.04 sys11"),
					resource.TestCheckResourceAttr("metakube_cluster.bar", "nodedepl.0.use_floating_ip", "false"),
				),
			},
			{
				Config: configUpdated,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckClusterResourceCreated("metakube_cluster.bar"),
					testAccCheckClustersNodeDeployment("metakube_cluster.bar", "my-nodedepl", "m1c.medium", "Rescue Ubuntu 18.04 sys11", true, 1, 1, 2),
					resource.TestCheckResourceAttr("metakube_cluster.bar", "name", "my-cluster-edit"),
					resource.TestCheckResourceAttr("metakube_cluster.bar", "labels.version", "beta"),
					resource.TestCheckResourceAttr("metakube_cluster.bar", "audit_logging", "false"),
					resource.TestCheckResourceAttr("metakube_cluster.bar", "nodedepl.0.replicas", "1"),
					resource.TestCheckResourceAttr("metakube_cluster.bar", "nodedepl.0.autoscale.0.min_replicas", "1"),
					resource.TestCheckResourceAttr("metakube_cluster.bar", "nodedepl.0.autoscale.0.max_replicas", "2"),
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

func testAccCheckClustersNodeDeployment(r, name, flavor, image string, floatingIP bool, replicas, minReplicas, maxReplicas uint) resource.TestCheckFunc {
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
		if nodedepl.Spec.MinReplicas != minReplicas {
			return fmt.Errorf("want min_replicas: %v, got: %v", minReplicas, nodedepl.Spec.MinReplicas)
		}
		if nodedepl.Spec.MaxReplicas != maxReplicas {
			return fmt.Errorf("want max_replicas: %v, got: %v", maxReplicas, nodedepl.Spec.MaxReplicas)
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
