package metakube

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"gitlab.com/furkhat/terraform-provider-metakube/gometakube"
)

// provider env.
const (
	APITokenEnvName = "METAKUBE_API_TOKEN"
)

// Provider returns MetaKube Provider.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc(APITokenEnvName, nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"metakube_project": resourceProject(),
			"metakube_cluster": resourceCluster(),
			"metakube_sshkey":  resourceSSHKey(),
		},
		ConfigureFunc: func(d *schema.ResourceData) (interface{}, error) {
			token := d.Get("token").(string)
			return gometakube.NewClient(gometakube.WithBearerToken(token)), nil
		},
	}
}
