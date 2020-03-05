package metakube

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"gitlab.com/furkhat/terraform-provider-metakube/gometakube"
)

func resourceCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceClusterCreate,
		Read:   resourceClusterRead,
		Update: resourceClusterUpdate,
		Delete: resourceClusterDelete,

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
				ForceNew:     true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"version": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"1.17.3", "1.15.10", "1.16.7"}, false),
			},
			"dc": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"provider_username": {
				Type:         schema.TypeString,
				Required:     true,
				Sensitive:    true,
				ValidateFunc: validation.NoZeroValues,
			},
			"provider_password": {
				Type:         schema.TypeString,
				Required:     true,
				Sensitive:    true,
				ValidateFunc: validation.NoZeroValues,
			},
			"nodepool": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true, // remove when in-place update implemented
							ValidateFunc: validation.NoZeroValues,
						},
						"replicas": {
							Type:         schema.TypeInt,
							Required:     true,
							ForceNew:     true, // remove when in-place update implemented
							ValidateFunc: validation.IntAtLeast(1),
						},
						"flavor_type": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringInSlice([]string{"Local Storage", "Network Storage"}, false),
						},
						"flavor": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.NoZeroValues,
						},
						"image": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.NoZeroValues,
						},
						"use_floating_ip": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
							Default:  true,
						},
					},
				},
			},
		},
	}
}

func resourceClusterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gometakube.Client)
	if dc, err := getClusterDatacenter(client, d.Get("dc").(string)); err != nil {
		return err
	} else if err := checkClusterNodepoolImage(client, dc, d); err != nil {
		return err
	} else {
		pool := d.Get("nodepool").([]interface{})[0].(map[string]interface{})
		create := &gometakube.CreateClusterRequest{
			Cluster: gometakube.Cluster{
				Name: d.Get("name").(string),
				Spec: &gometakube.ClusterSpec{
					Version: d.Get("version").(string),
					Cloud: &gometakube.ClusterSpecCloud{
						OpenStack: &gometakube.ClusterSpecCloudOpenstack{
							Domain:         "Default",
							Username:       d.Get("provider_username").(string),
							Password:       d.Get("provider_password").(string),
							FloatingIPPool: "ext-net",
						},
						DataCenter: d.Get("dc").(string),
					},
					MachineNetworks: []gometakube.ClusterSpecMachineNetwork{},
				},
				Type:    "kubernetes",
				SSHKeys: []string{},
			},
			NodeDeployment: gometakube.NodeDeployment{
				Name: pool["name"].(string),
				Spec: gometakube.NodeDeploymentSpec{
					Template: gometakube.NodeDeploymentSpecTemplate{
						Cloud: gometakube.NodeDeploymentSpecTemplateCloud{
							Openstack: gometakube.NodeDeploymentSpecTemplateCloudOpenstack{
								FlavorType:    pool["flavor_type"].(string),
								Flavor:        pool["flavor"].(string),
								Image:         pool["image"].(string),
								UseFloatingIP: pool["use_floating_ip"].(bool),
							},
						},
						OperatingSystem: gometakube.NodeDeploymentSpecTemplateOS{
							Ubuntu: &gometakube.NodeDeploymentSpecTemplateOSOptions{
								DistUpgradeOnBoot: new(bool),
							},
						},
					},
					Replicas: uint(pool["replicas"].(int)),
				},
			},
		}
		projectID := d.Get("project_id").(string)
		cluster, err := client.Clusters.Create(context.Background(), projectID, dc.Spec.Seed, create)
		if err != nil {
			return fmt.Errorf("could not create cluster: %v", err)
		}
		d.SetId(cluster.ID)
		return waitForClusterHealthy(client, projectID, dc.Spec.Seed, cluster.ID)
	}
}

func resourceClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gometakube.Client)
	id := d.Id()
	project := d.Get("project_id").(string)
	if dc, err := getClusterDatacenter(client, d.Get("dc").(string)); err != nil {
		return err
	} else if obj, err := getCluster(client, project, dc.Spec.Seed, id); err != nil {
		return err
	} else if obj == nil || obj.DeletionTimestamp != nil {
		// Cluster was deleted
		d.SetId("")
		return nil
	} else {
		d.Set("name", obj.Name)
		d.Set("version", obj.Spec.Version)
		d.Set("dc", obj.Spec.Cloud.DataCenter)
		// TODO: update nodepool vals
		return nil
	}
}

func resourceClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(true)
	if d.HasChange("name") {
		client := meta.(*gometakube.Client)
		projectID := d.Get("project_id").(string)
		if dc, err := getClusterDatacenter(client, d.Get("dc").(string)); err != nil {
			return err
		} else if cluster, err := getCluster(client, projectID, dc.Spec.Seed, d.Id()); err != nil {
			return err
		} else if cluster == nil {
			// Cluster was deleted
			return nil
		} else {
			patch := &gometakube.PatchClusterRequest{
				Name: d.Get("name").(string),
			}
			_, err = client.Clusters.Patch(context.Background(), projectID, dc.Spec.Seed, d.Id(), patch)
			if err != nil {
				return fmt.Errorf("could not patch cluster (is cluster provisioning compete?). error: %v", err)
			}
			d.SetPartial("name")
		}
	}
	d.Partial(false)
	return nil
	// TODO: patch nodepool if nodepool's name changes
}

func resourceClusterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gometakube.Client)
	id := d.Id()
	project := d.Get("project_id").(string)
	if dc, err := getClusterDatacenter(client, d.Get("dc").(string)); err != nil {
		return err
	} else if err := client.Clusters.Delete(context.Background(), project, dc.Spec.Seed, id); err != nil {
		return fmt.Errorf("could not delete cluster: %v", err)
	} else {
		return nil
	}
}

func waitForClusterHealthy(client *gometakube.Client, prj, dc, id string) error {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	timeout := 10 * 60
	n := 0
	for range ticker.C {
		h, _ := client.Clusters.Health(context.Background(), prj, dc, id)
		if h != nil && h.Healthy() {
			return nil
		}
		if n > timeout {
			return fmt.Errorf("Timeout waiting to create cluster")
		}
		n++
	}
	return nil
}

func checkClusterNodepoolImage(client *gometakube.Client, dc *gometakube.Datacenter, d *schema.ResourceData) error {
	providerUsername := d.Get("provider_username").(string)
	providerPassword := d.Get("provider_password").(string)
	images, err := client.Images.List(context.Background(), dc.Metadata.Name, "Default", providerUsername, providerPassword)
	if err != nil {
		return fmt.Errorf("could not get list of images: %v", err)
	}
	pool := d.Get("nodepool").([]interface{})[0].(map[string]interface{})
	imageName := pool["image"].(string)
	for _, image := range images {
		if image.Name == imageName {
			return nil
		}
	}
	availableImages := make([]string, 0)
	for _, image := range images {
		availableImages = append(availableImages, image.Name)
	}
	return fmt.Errorf("image `%s` is not avaialable in datacenter `%s`. Consider changing to one of:\n%s",
		dc.Metadata.Name,
		imageName,
		strings.Join(availableImages, "\n"))
}

func getClusterDatacenter(c *gometakube.Client, n string) (*gometakube.Datacenter, error) {
	dc, err := c.Datacenters.Get(context.Background(), n)
	if err != nil {
		return nil, fmt.Errorf("could not get details on datacenter: %v", err)
	}
	return dc, nil
}

func getCluster(c *gometakube.Client, prj, dc, id string) (*gometakube.Cluster, error) {
	obj, err := c.Clusters.Get(context.Background(), prj, dc, id)
	if err != nil {
		return nil, fmt.Errorf("could not get cluster details: %v", err)
	}
	return obj, nil
}
