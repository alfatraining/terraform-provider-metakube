package metakube

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"gitlab.com/furkhat/terraform-provider-metakube/gometakube"
)

func resourceProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectCreate,
		Read:   resourceProjectRead,
		Update: resourceProjectUpdate,
		Delete: resourceProjectDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"labels": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     schema.TypeString,
			},
		},
	}
}

func resourceProjectCreate(d *schema.ResourceData, meta interface{}) error {
	create := &gometakube.ProjectCreateAndUpdateRequest{
		Name:   d.Get("name").(string),
		Labels: projectLabelsMap(d),
	}
	client := meta.(*gometakube.Client)
	project, _, err := client.Projects.Create(context.Background(), create)
	if err != nil {
		return fmt.Errorf("could not create project: %v", err)
	}
	d.SetId(project.ID)
	return waitProjectCreatedAndActive(client, project.ID)
}

func resourceProjectRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*gometakube.Client)
	obj, _, err := c.Projects.Get(context.Background(), d.Id())
	if err != nil {
		return err
	}
	if obj == nil || obj.DeletionTimestamp != nil {
		// Project was deleted.
		d.SetId("")
		return nil
	}
	d.Set("name", obj.Name)
	d.Set("labels", obj.Labels)
	return nil
}

func resourceProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	d.Partial(true)
	defer d.Partial(false)

	update := &gometakube.ProjectCreateAndUpdateRequest{
		Name:   d.Get("name").(string),
		Labels: projectLabelsMap(d),
	}
	client := meta.(*gometakube.Client)
	updated, _, err := client.Projects.Update(context.Background(), d.Id(), update)
	if err != nil {
		return err
	}
	if updated == nil || updated.DeletionTimestamp != nil {
		// Project was deleted.
		d.SetId("")
		return nil
	}
	d.SetPartial("name")
	d.SetPartial("labels")

	return nil
}

func resourceProjectDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*gometakube.Client)
	_, err := c.Projects.Delete(context.Background(), d.Id())
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func waitProjectCreatedAndActive(client *gometakube.Client, id string) error {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	timeout := 120
	n := 0
	for range ticker.C {
		project, _, err := client.Projects.Get(context.Background(), id)
		if err == nil && project.Status == "Active" {
			return nil
		}
		if n > timeout {
			return fmt.Errorf("Timeout waiting for project activation")
		}
		n++
	}
	return nil
}

func projectLabelsMap(d *schema.ResourceData) (ret map[string]string) {
	return labelsMap(d)
}
