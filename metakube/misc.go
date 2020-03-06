package metakube

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

func labelsMap(d *schema.ResourceData) (ret map[string]string) {
	if attr, ok := d.GetOk("labels"); ok {
		ret = make(map[string]string)
		for k, v := range attr.(map[string]interface{}) {
			ret[k] = v.(string)
		}
	}
	return ret
}
