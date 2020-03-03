package metakube

import (
	"context"
	"fmt"

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
	create := &gometakube.ProjectCreateRequest{}
	create.Name = d.Get("name").(string)
	if attr, ok := d.GetOk("labels"); ok {
		create.Labels = make(map[string]string)
		for k, v := range attr.(map[string]interface{}) {
			create.Labels[k] = v.(string)
		}
	}
	c := meta.(*gometakube.Client)
	project, err := c.Projects.Create(context.Background(), create)
	if err != nil {
		return fmt.Errorf("could not create project: %v", err)
	}
	d.SetId(project.ID)
	return nil
}

func resourceProjectRead(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*gometakube.Client)
	obj, err := c.Projects.Get(context.Background(), d.Id())
	if err != nil {
		return err
	}
	if obj == nil || obj.DeletionTimestamp != nil {
		d.SetId("")
		return nil
	}
	d.Set("name", obj.Name)
	d.Set("labels", obj.Labels)
	return nil
}

func resourceProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	project := &gometakube.ProjectCreateRequest{}
	project.Name = d.Get("name").(string)
	if attr, ok := d.GetOk("labels"); ok {
		project.Labels = make(map[string]string)
		for k, v := range attr.(map[string]interface{}) {
			project.Labels[k] = v.(string)
		}
	}
	c := meta.(*gometakube.Client)
	updated, err := c.Projects.Update(context.Background(), d.Id(), project)
	if err != nil {
		return err
	}
	if updated == nil || updated.DeletionTimestamp != nil {
		d.SetId("")
		return nil
	}
	d.Set("name", updated.Name)
	d.Set("labels", updated.Name)
	return nil
}

func resourceProjectDelete(d *schema.ResourceData, meta interface{}) error {
	c := meta.(*gometakube.Client)
	err := c.Projects.Delete(context.Background(), d.Id())
	// HACK: gometakube returns ErrForbidden even for non-existing resources -
	// handling this as resource absence.
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
