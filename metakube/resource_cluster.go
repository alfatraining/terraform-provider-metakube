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
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"sshkeys": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.NoZeroValues,
				},
			},
			"version": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
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
			"provider_username": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Sensitive:    true,
				ValidateFunc: validation.NoZeroValues,
			},
			"provider_password": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Sensitive:    true,
				ValidateFunc: validation.NoZeroValues,
			},
			"audit_logging": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"nodedepl": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
							ForceNew: true,
						},
						"replicas": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntAtLeast(1),
						},
						"autoscale": {
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"min_replicas": {
										Type:     schema.TypeInt,
										Optional: true,
										Default:  0,
									},
									"max_replicas": {
										Type:     schema.TypeInt,
										Optional: true,
										Default:  0,
									},
								},
							},
						},
						"flavor": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.NoZeroValues,
						},
						"image": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.NoZeroValues,
						},
						"use_floating_ip": {
							Type:     schema.TypeBool,
							Optional: true,
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
	if minReplicas, maxReplicas, err := checkClusterAutoscaleValid(d); err != nil {
		return err
	} else if dc, err := getClusterDatacenter(client, d.Get("dc").(string)); err != nil {
		return err
	} else if err := checkClusterTenantValid(client, dc, d); err != nil {
		return err
	} else if err := checkClusterNodedeplImage(client, dc, d); err != nil {
		return err
	} else if project, err := client.Projects.Get(context.Background(), d.Get("project_id").(string)); err != nil {
		return err
	} else if err := checkClusterDoesNotRedefineProjectLabels(project, d); err != nil {
		return err
	} else if clusterVersion, err := getClusterVersionToUse(client, d.Get("version").(string)); err != nil {
		return err
	} else {
		nodedepl := d.Get("nodedepl").([]interface{})[0].(map[string]interface{})
		create := &gometakube.CreateClusterRequest{
			Cluster: gometakube.Cluster{
				Name:   d.Get("name").(string),
				Labels: clusterLabelsMap(d),
				Spec: &gometakube.ClusterSpec{
					Version: clusterVersion,
					AuditLogging: gometakube.ClusterSpecAuditLogging{
						Enabled: d.Get("audit_logging").(bool),
					},
					Cloud: &gometakube.ClusterSpecCloud{
						OpenStack: &gometakube.ClusterSpecCloudOpenstack{
							Domain:         "Default",
							Tenant:         d.Get("tenant").(string),
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
				Name: nodedepl["name"].(string),
				Spec: gometakube.NodeDeploymentSpec{
					Template: gometakube.NodeDeploymentSpecTemplate{
						Cloud: gometakube.NodeDeploymentSpecTemplateCloud{
							Openstack: gometakube.NodeDeploymentSpecTemplateCloudOpenstack{
								Flavor:        nodedepl["flavor"].(string),
								Image:         nodedepl["image"].(string),
								UseFloatingIP: nodedepl["use_floating_ip"].(bool),
							},
						},
						OperatingSystem: gometakube.NodeDeploymentSpecTemplateOS{
							Ubuntu: &gometakube.NodeDeploymentSpecTemplateOSOptions{
								DistUpgradeOnBoot: new(bool),
							},
						},
					},
					Replicas:    uint(nodedepl["replicas"].(int)),
					MinReplicas: uint(minReplicas),
					MaxReplicas: uint(maxReplicas),
				},
			},
		}
		prj := d.Get("project_id").(string)
		obj, err := client.Clusters.Create(context.Background(), prj, dc.Spec.Seed, create)
		if err != nil {
			return fmt.Errorf("could not create cluster: %v", err)
		}
		d.SetId(obj.ID)
		if err := manageSSHKeysInCluster(client, nil, d.Get("sshkeys"), prj, dc.Spec.Seed, d.Id()); err != nil {
			return err
		}
		return waitForClusterRunningAndNodeDeploymentCreate(client, prj, dc.Spec.Seed, obj.ID, d.Get("nodedepl.0.name").(string))
	}
}

func resourceClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gometakube.Client)
	id := d.Id()
	projectID := d.Get("project_id").(string)
	if dc, err := getClusterDatacenter(client, d.Get("dc").(string)); err != nil {
		return err
	} else if obj, err := getCluster(client, projectID, dc.Spec.Seed, id); err != nil {
		return err
	} else if obj == nil || obj.DeletionTimestamp != nil {
		// Cluster was deleted
		d.SetId("")
		return nil
	} else if nodeDeployment, err := getClusterNodeDeployment(client, projectID, dc.Spec.Seed, id, d.Get("nodedepl.0.name").(string)); err != nil {
		return err
	} else if project, err := client.Projects.Get(context.Background(), projectID); err != nil {
		return err
	} else if sshkeys, err := client.SSHKeys.ListAssigned(context.Background(), projectID, dc.Spec.Seed, id); err != nil {
		return fmt.Errorf("could not list sshkeys assinged to cluster: %v", err)
	} else {
		d.Set("name", obj.Name)
		labelsToSet := obj.Labels
		for k := range project.Labels {
			delete(labelsToSet, k)
		}
		d.Set("labels", labelsToSet)
		version := d.Get("version").(string)
		if obj.Spec.Version[:len(version)] != version {
			d.Set("version", obj.Spec.Version)
		}
		d.Set("dc", obj.Spec.Cloud.DataCenter)
		d.Set("audit_logging", obj.Spec.AuditLogging.Enabled)

		d.Set("nodedepl", nodeDeploymentUpdatesMap(nodeDeployment))

		keynames := make([]string, 0)
		for _, key := range sshkeys {
			keynames = append(keynames, key.Name)
		}
		d.Set("sshkeys", keynames)
		return nil
	}
}

func resourceClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(true)
	defer d.Partial(false)
	client := meta.(*gometakube.Client)
	projectID := d.Get("project_id").(string)
	dc, err := getClusterDatacenter(client, d.Get("dc").(string))
	if err != nil {
		return err
	}
	if d.HasChanges("name", "labels", "audit_logging") {
		if cluster, err := getCluster(client, projectID, dc.Spec.Seed, d.Id()); err != nil {
			return err
		} else if cluster == nil {
			// Cluster was deleted
			d.SetId("")
			return nil
		} else if project, err := client.Projects.Get(context.Background(), d.Get("project_id").(string)); err != nil {
			return err
		} else if err := checkClusterDoesNotRedefineProjectLabels(project, d); err != nil {
			return err
		} else {
			patch := &gometakube.PatchClusterRequest{
				Name:   d.Get("name").(string),
				Labels: clusterLabelsMap(d),
				Spec: &gometakube.PatchClusterRequestSpec{
					AuditLogging: &gometakube.ClusterSpecAuditLogging{
						Enabled: d.Get("audit_logging").(bool),
					},
				},
			}
			_, err = client.Clusters.Patch(context.Background(), projectID, dc.Spec.Seed, d.Id(), patch)
			if err != nil {
				return fmt.Errorf("could not patch cluster (is cluster provisioning compete?). error: %v", err)
			}
			d.SetPartial("name")
			d.SetPartial("labels")
			d.SetPartial("audit_logging")
		}
	}
	if d.HasChange("nodedepl") {
		if minReplicas, maxReplicas, err := checkClusterAutoscaleValid(d); err != nil {
			return err
		} else if nodedepl, err := getClusterNodeDeployment(client, projectID, dc.Spec.Seed, d.Id(), d.Get("nodedepl.0.name").(string)); err != nil {
			return err
		} else if err := checkClusterNodedeplImage(client, dc, d); err != nil {
			return err
		} else {
			patch := &gometakube.NodeDeploymentsPatchRequest{Spec: nodedepl.Spec}
			patch.Spec.Replicas = uint(d.Get("nodedepl.0.replicas").(int))
			patch.Spec.MinReplicas = uint(minReplicas)
			patch.Spec.MaxReplicas = uint(maxReplicas)
			patch.Spec.Template.Cloud.Openstack.Flavor = d.Get("nodedepl.0.flavor").(string)
			patch.Spec.Template.Cloud.Openstack.Image = d.Get("nodedepl.0.image").(string)
			patch.Spec.Template.Cloud.Openstack.UseFloatingIP = d.Get("nodedepl.0.use_floating_ip").(bool)
			_, err = client.NodeDeployments.Patch(context.Background(), projectID, dc.Spec.Seed, d.Id(), nodedepl.ID, patch)
			if err != nil {
				return fmt.Errorf("could not patch node deployment: %v", err)
			}
			d.SetPartial("nodedepl")
		}
	}
	if d.HasChange("version") {
		versionPrefix := d.Get("version").(string)
		if cluster, err := client.Clusters.Get(context.Background(), projectID, dc.Spec.Seed, d.Id()); err != nil {
			return err
		} else if clusterVersionsHasPrefix(cluster.Spec.Version, versionPrefix) {
			return nil
		} else if versionToUse, err := getClusterVersionToUse(client, versionPrefix); err != nil {
			return err
		} else if invalidUpgrade, err := clusterVersionBigger(cluster.Spec.Version, versionToUse); err != nil {
			return nil
		} else if invalidUpgrade {
			return fmt.Errorf("cluster version cannot be downgraded")
		} else {
			// Upgrade cluster continuously to desired version.
			for {
				version, err := getClusterVersionToUpgradeInto(client, projectID, dc.Spec.Seed, d.Id())
				if version == "" {
					if clusterVersionsHasPrefix(cluster.Spec.Version, versionPrefix) {
						break
					}
					return fmt.Errorf("cluster has no further upgrades, stuck at %s", cluster.Spec.Version)
				}
				patch := &gometakube.PatchClusterRequest{
					Spec: &gometakube.PatchClusterRequestSpec{
						Version: version,
					},
				}
				cluster, err = client.Clusters.Patch(context.Background(), projectID, dc.Spec.Seed, d.Id(), patch)
				if err != nil {
					return fmt.Errorf("could not patch cluster (is cluster provisioning compete?). error: %v", err)
				}
				if err := waitForClusterHealthy(client, projectID, dc.Spec.Seed, d.Id()); err != nil {
					return err
				}
				if cluster.Spec.Version == versionToUse {
					break
				}
			}
			err = client.NodeDeployments.Upgrade(context.Background(), projectID, dc.Spec.Seed, d.Id(), &gometakube.UpgradeNodesRequest{
				Version: cluster.Spec.Version,
			})
			if err != nil {
				return fmt.Errorf("could not upgrade cluster nodes: %v", err)
			}
		}
	}
	if d.HasChange("sshkeys") {
		old, new := d.GetChange("sshkeys")
		if err := manageSSHKeysInCluster(client, old, new, projectID, dc.Spec.Seed, d.Id()); err != nil {
			return err
		}
	}
	return nil
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
		return waitForClusterDelete(client, project, dc.Spec.Seed, id)
	}
}

func manageSSHKeysInCluster(client *gometakube.Client, old, new interface{}, prj, dc, cls string) error {
	allKeys, err := client.SSHKeys.List(context.Background(), prj)
	if err != nil {
		return fmt.Errorf("could not get list of available sshkeys to assign: %v", err)
	}
	actions := make(map[string]bool)
	if old != nil {
		for _, v := range old.(*schema.Set).List() {
			actions[v.(string)] = false
		}
	}
	if new != nil {
		for _, v := range new.(*schema.Set).List() {
			actions[v.(string)] = true
		}
	}
	for v, action := range actions {
		// action is true -> assign, action is false -> unassign.
		var id string
		for _, key := range allKeys {
			if key.Name == v {
				id = key.ID
			}
		}
		if action {
			if id == "" {
				return fmt.Errorf("no ssh key with name `%s` found the the project", v)
			}
			_, err = client.SSHKeys.AssignToCluster(context.Background(), prj, dc, cls, id)
			if err != nil {
				return fmt.Errorf("could not assign sshkey to cluster: %v", err)
			}
		} else if id != "" {
			err = client.SSHKeys.RemoveFromCluster(context.Background(), prj, dc, cls, id)
			if err != nil {
				return fmt.Errorf("could not unassign sshkey from cluster: %v", err)
			}
		}
	}
	return nil
}

func waitForClusterDelete(client *gometakube.Client, prj, dc, id string) error {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	timeout := 10 * 60
	n := 0
	for range ticker.C {
		v, err := client.Clusters.Get(context.Background(), prj, dc, id)
		if err != nil {
			return fmt.Errorf("couldn't check cluster delete: %v", err)
		}
		if v == nil {
			return nil
		}
		if n > timeout {
			return fmt.Errorf("Timeout waiting to create delete")
		}
		n++
	}
	return nil
}

func waitForClusterRunningAndNodeDeploymentCreate(client *gometakube.Client, prj, dc, cls, nodedepl string) error {
	if err := waitForClusterHealthy(client, prj, dc, cls); err != nil {
		return err
	}
	return waitNodeDeploymentCreate(client, prj, dc, cls, nodedepl)
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

func waitNodeDeploymentCreate(client *gometakube.Client, prj, dc, cls, name string) (err error) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	timeout := 5 * 60
	n := 0
	for range ticker.C {
		if _, err = getClusterNodeDeployment(client, prj, dc, cls, name); err == nil {
			return nil
		}
		if n > timeout {
			return fmt.Errorf("Timeout waiting to create cluster node deployment: %v", err)
		}
		n++
	}
	return err
}

func checkClusterNodedeplImage(client *gometakube.Client, dc *gometakube.Datacenter, d *schema.ResourceData) error {
	providerUsername := d.Get("provider_username").(string)
	providerPassword := d.Get("provider_password").(string)
	images, err := client.Openstack.Images(context.Background(), dc.Metadata.Name, "Default", providerUsername, providerPassword)
	if err != nil {
		return fmt.Errorf("could not get list of images: %v", err)
	}
	nodedepl := d.Get("nodedepl").([]interface{})[0].(map[string]interface{})
	imageName := nodedepl["image"].(string)
	for _, image := range images {
		if image.Name == imageName {
			return nil
		}
	}
	availableImages := make([]string, 0)
	for _, image := range images {
		availableImages = append(availableImages, "* "+image.Name)
	}
	return fmt.Errorf("image `%s` is not avaialable in datacenter `%s`. Consider changing to one of:\n%s",
		imageName,
		dc.Metadata.Name,
		strings.Join(availableImages, "\n"))
}

func checkClusterTenantValid(client *gometakube.Client, dc *gometakube.Datacenter, d *schema.ResourceData) error {
	providerUsername := d.Get("provider_username").(string)
	providerPassword := d.Get("provider_password").(string)
	tenants, err := client.Openstack.Tenants(context.Background(), dc.Metadata.Name, "Default", providerUsername, providerPassword)
	if err != nil {
		return fmt.Errorf("could not get list of tenants: %v", err)
	}
	specified := d.Get("tenant").(string)
	for _, t := range tenants {
		if t.Name == specified {
			return nil
		}
	}
	available := make([]string, 0)
	for _, t := range tenants {
		available = append(available, "* "+t.Name)
	}
	return fmt.Errorf("tenant `%s` is not avaialable in datacenter `%s`. Consider changing to one of:\n%s",
		specified,
		dc.Metadata.Name,
		strings.Join(available, "\n"))
}

func checkClusterDoesNotRedefineProjectLabels(project *gometakube.Project, d *schema.ResourceData) error {
	clusterLabels := d.Get("labels").(map[string]interface{})
	for k := range project.Labels {
		if v, ok := clusterLabels[k]; ok {
			return fmt.Errorf("cannot change labels inherited from project: %v=%v", k, v)
		}
	}
	return nil
}

func getClusterVersionToUse(c *gometakube.Client, prefix string) (string, error) {
	versions, err := c.Clusters.Upgrades(context.Background())
	if err != nil {
		return "", fmt.Errorf("could not get available cluster versions: %v", err)
	}
	return maxVersionWithPrefix(versions, prefix)
}

func getClusterVersionToUpgradeInto(c *gometakube.Client, prj, dc, id string) (string, error) {
	versions, err := c.Clusters.ClusterUpgrades(context.Background(), prj, dc, id)
	if err != nil {
		return "", err
	}
	return maxVersionWithPrefix(versions, "")
}

func maxVersionWithPrefix(versions []gometakube.ClusterUpgrade, prefix string) (string, error) {
	if len(versions) == 0 {
		return "", fmt.Errorf("empty list of cluster versions returned from api")
	}
	ret := ""
	versionsStr := make([]string, 0)
	for _, item := range versions {
		if clusterVersionsHasPrefix(item.Version, prefix) {
			if ret == "" {
				ret = item.Version
			} else if bigger, err := clusterVersionBigger(item.Version, ret); err != nil {
				return "", err
			} else if bigger {
				ret = item.Version
			}
		}
		versionsStr = append(versionsStr, item.Version)
	}
	if ret == "" {
		return "", fmt.Errorf("could not find versions prefixed by: %s, available versions are: %s", prefix, strings.Join(versionsStr, ", "))
	}
	return ret, nil
}

func checkClusterAutoscaleValid(d *schema.ResourceData) (int, int, error) {
	minReplicas := d.Get("nodedepl.0.autoscale.0.min_replicas").(int)
	maxReplicas := d.Get("nodedepl.0.autoscale.0.max_replicas").(int)
	if minReplicas == 0 && maxReplicas == 0 {
		return 0, 0, nil
	}
	if minReplicas > maxReplicas {
		return 0, 0, fmt.Errorf("autoscale min_replicas(%d) must be less than or equal to max_replicas(%d)", minReplicas, maxReplicas)
	}
	replicas := d.Get("nodedepl.0.replicas").(int)
	if replicas > maxReplicas || replicas < minReplicas {
		return 0, 0, fmt.Errorf("got autoscale settings [%d; %d], but replicas: %d", minReplicas, maxReplicas, replicas)
	}
	return minReplicas, maxReplicas, nil
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

func clusterLabelsMap(d *schema.ResourceData) (ret map[string]string) {
	return labelsMap(d)
}

func getClusterNodeDeployment(c *gometakube.Client, prj, dc, cls, name string) (*gometakube.NodeDeployment, error) {
	items, err := c.NodeDeployments.List(context.Background(), prj, dc, cls)
	if err != nil {
		return nil, fmt.Errorf("could not get cluster node deployments: %v", err)
	}
	for _, item := range items {
		if item.Name == name {
			return &item, nil
		}
	}
	return nil, fmt.Errorf("could not find node deployment with given name: %s", name)
}

func nodeDeploymentUpdatesMap(nodedepl *gometakube.NodeDeployment) []interface{} {
	return []interface{}{map[string]interface{}{
		"name":     nodedepl.Name,
		"replicas": nodedepl.Spec.Replicas,
		"autoscale": []interface{}{map[string]interface{}{
			"min_replicas": nodedepl.Spec.MinReplicas,
			"max_replicas": nodedepl.Spec.MaxReplicas,
		}},
		"flavor":          nodedepl.Spec.Template.Cloud.Openstack.Flavor,
		"image":           nodedepl.Spec.Template.Cloud.Openstack.Image,
		"use_floating_ip": nodedepl.Spec.Template.Cloud.Openstack.UseFloatingIP,
	}}
}
