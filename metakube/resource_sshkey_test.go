package metakube

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	testSSHPubKey        = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCzoO6BIidD4Us9a9Kh0GzaUUxosl61GNUZzqcIdmf4EYZDdRtLa+nu88dHPHPQ2dj52BeVV9XVN9EufqdAZCaKpPLj5XxEwMpGcmdrOAl38kk2KKbiswjXkrdhYSBw3w0KkoCPKG/+yNpAUI9z+RJZ9lukeYBvxdDe8nuvUWX7mGRaPaumCpQaBHwYKNn6jMVns2RrumgE9w+Z6jlaKHk1V7T5rCBDcjXwcy6waOX6hKdPPBk84FpUfcfN/SdpwSVGFrcykazrpmzD2nYr71EcOm9T6/yuhBOiIa3H/TOji4G9wr02qtSWuGUpULkqWMFD+BQcYQQA71GSAa+rTZuf user@machine.local"
	testAccSSHKeyConfig1 = `
provider "metakube" {
}

resource "metakube_project" "sshkey-project" {
	name = "foo"
	labels = {}
}

resource "metakube_sshkey" "test-sshkey" {
	project_id = metakube_project.sshkey-project.id

	name = "my-key"
	public_key = "` + testSSHPubKey + `"
}
`
)

func TestAccMetakubeSSHKey_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testEnvSet(t, APITokenEnvName)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMetakubeProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSSHKeyConfig1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("metakube_sshkey.test-sshkey", "name", "my-key"),
					resource.TestCheckResourceAttr("metakube_sshkey.test-sshkey", "public_key", testSSHPubKey),
				),
			},
		},
	})
}
