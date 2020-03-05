package metakube

import (
	"context"
	"fmt"
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
			"tenant": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"username": {
				Type:         schema.TypeString,
				Required:     true,
				Sensitive:    true,
				ValidateFunc: validation.NoZeroValues,
			},
			"password": {
				Type:         schema.TypeString,
				Required:     true,
				Sensitive:    true,
				ValidateFunc: validation.NoZeroValues,
			},
			"node_pool": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.NoZeroValues,
						},
						"replicas": {
							Type:         schema.TypeInt,
							Required:     true,
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
	c := meta.(*gometakube.Client)
	dc, err := c.Datacenters.Get(context.Background(), d.Get("dc").(string))
	if err != nil {
		return fmt.Errorf("could not get details on datacenter: %v", err)
	}
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	images, err := c.Images.List(context.Background(), dc.Metadata.Name, "Default", username, password)
	if err != nil {
		return fmt.Errorf("could not get list of images: %v", err)
	}
	pool := d.Get("node_pool").([]interface{})[0].(map[string]interface{})
	// TODO: extract to func
	imageName := pool["image"].(string)
	imageValid := false
	for _, image := range images {
		if image.Name == imageName {
			imageValid = true
			break
		}
	}
	if !imageValid {
		return fmt.Errorf("image `%s` is not avaialable", imageName)
	}
	create := &gometakube.CreateClusterRequest{
		Cluster: gometakube.Cluster{
			Name: d.Get("name").(string),
			Spec: &gometakube.ClusterSpec{
				Version: d.Get("version").(string),
				Cloud: &gometakube.ClusterSpecCloud{
					OpenStack: &gometakube.ClusterSpecCloudOpenstack{
						Tenant:         d.Get("tenant").(string),
						Domain:         "Default",
						Username:       username,
						Password:       password,
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
							Image:         imageName,
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
	// TODO: proper cancellation
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	project := new(gometakube.Project)
	for project == nil || project.Status != "Active" {
		select {
		case <-ctx.Done():
			break
		default:
			project, err = c.Projects.Get(context.Background(), projectID)
		}
	}
	if err != nil {
		return fmt.Errorf("could not get project: %v", err)
	}
	if project == nil {
		return fmt.Errorf("project with id: `%s` does not exist", projectID)
	}
	cluster, err := c.Clusters.Create(context.Background(), projectID, dc.Spec.Seed, create)
	if err != nil {
		return fmt.Errorf("could not create cluster: %v", err)
	}
	d.SetId(cluster.ID)
	return nil
}

func resourceClusterRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*gometakube.Client)
	id := d.Id()
	project := d.Get("project_id").(string)
	dc, err := c.Datacenters.Get(context.Background(), d.Get("dc").(string))
	if err != nil {
		return fmt.Errorf("could not get details on datacenter: %v", err)
	}
	cluster, err := c.Clusters.Get(context.Background(), project, dc.Spec.Seed, id)
	if err != nil {
		return fmt.Errorf("could not get cluster details: %v", err)
	}
	if cluster == nil || cluster.DeletionTimestamp != nil {
		d.SetId("")
		return nil
	}
	return nil
}

func resourceClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*gometakube.Client)
	projectID := d.Get("project_id").(string)
	dc, err := c.Datacenters.Get(context.Background(), d.Get("dc").(string))
	if err != nil {
		return fmt.Errorf("could not get details on datacenter: %v", err)
	}
	cluster, err := c.Clusters.Get(context.Background(), projectID, dc.Spec.Seed, d.Id())
	if err != nil {
		return fmt.Errorf("could not retrieve cluster: %v", err)
	}
	if cluster == nil {
		// Cluster was deleted
		d.SetId("")
		return nil
	}
	patch := &gometakube.PatchClusterRequest{
		Name: d.Get("name").(string),
	}
	_, err = c.Clusters.Patch(context.Background(), projectID, dc.Spec.Seed, d.Id(), patch)
	if err != nil {
		return fmt.Errorf("could not patch cluster: %v", err)
	}
	return nil
	// TODO: patch nodepool if nodepool's name changes
}

func resourceClusterDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*gometakube.Client)
	id := d.Id()
	project := d.Get("project_id").(string)
	dc, err := c.Datacenters.Get(context.Background(), d.Get("dc").(string))
	if err != nil {
		return fmt.Errorf("could not get details on datacenter: %v", err)
	}
	if err := c.Clusters.Delete(context.Background(), project, dc.Spec.Seed, id); err != nil {
		return fmt.Errorf("could not delete cluster: %v", err)
	}
	return nil
}
