package metakube

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"gitlab.com/furkhat/terraform-provider-metakube/gometakube"
)

func resourceSSHKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceSSHKeyCreate,
		Read:   resourceSSHKeyRead,
		Delete: resourceSSHKeyDelete,

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
				ForceNew:     true,
			},
			"public_key": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
				DiffSuppressFunc: func(_, old, new string, _ *schema.ResourceData) bool {
					return strings.TrimSpace(old) == strings.TrimSpace(new)
				},
				ForceNew: true,
			},
		},
	}
}

func resourceSSHKeyCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*gometakube.Client)
	v, _, err := client.SSHKeys.Create(context.Background(), d.Get("project_id").(string), &gometakube.SSHKey{
		Name: d.Get("name").(string),
		Spec: gometakube.SSHKeySpec{
			PublicKey: d.Get("public_key").(string),
		},
	})
	if err != nil {
		return err
	}
	d.SetId(v.ID)
	return nil
}

func resourceSSHKeyRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*gometakube.Client)
	sshkeys, err := client.SSHKeys.List(context.Background(), d.Get("project_id").(string))
	if err != nil {
		return fmt.Errorf("could not get projects sshkeys: %v", err)
	}
	var v gometakube.SSHKey
	for _, sshkey := range sshkeys {
		if sshkey.ID == d.Id() {
			v = sshkey
		}
	}
	if v.ID != d.Id() || v.DeletionTimestamp != nil {
		// SSHKey not found in the project, it surely was deleted.
		d.SetId("")
		return nil
	}
	d.Set("name", v.Name)
	d.Set("public_key", v.Spec.PublicKey)
	return nil
}

func resourceSSHKeyDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*gometakube.Client)
	_, err := client.SSHKeys.Delete(context.Background(), d.Get("project_id").(string), d.Id())
	return err
}
