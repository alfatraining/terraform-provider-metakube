provider "metakube" {
  // Do not forget to set METAKUBE_API_TOKEN environment variable.
}

resource "metakube_project" "my-project" {
  // project name, in-place updatable
  name = "my-project"

  // project labels, in-place updatable
  labels = {
    "component" = "main"
  }
}

resource "metakube_cluster" "my-cluster" {
  project_id = metakube_project.my-project.id // change forces new

  labels = { // has in-place update.
    "environment" = "staging"
  }

  sshkeys = [
  ]

  name          = "my-cluster"     // has in-place update
  version       = "1.17"           // k8s version. will use the version that has such prefix
  dc            = "syseleven-dbl1" // openstack datacenter, change forces new
  audit_logging = true             // has in-place update

  // openstack 
  tenant            = "" // change forces new
  provider_username = "" // sensitive, not persisted in tfstate, change forces new
  provider_password = "" // sensitive, not persisted in tfstate, change forces new

  // clusters node deployment
  nodedepl {
    name     = "my-cluster-nodedepl-one" // change forces new
    replicas = 1                         // has in-place update

    autoscale {
      min_replicas = 1 // optional, not setting and setting to zero have the same effect.
      max_replicas = 2 // optional, not setting and setting to zero have the same effect.
    }

    flavor          = "l1.small"                  // has in-place update
    image           = "Rescue Ubuntu 18.04 sys11" // has in-place update
    use_floating_ip = false                       // has in-place update
  }
}

resource "metakube_sshkey" "my-key" {
  project_id = metakube_project.my-project.id // change foreces new

  name = "some"

  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCut5oRyqeqYci3E9m6Z6mtxfqkiyb+xNFJM6+/sllhnMDX0vzrNj8PuIFfGkgtowKY//QWLgoB+RpvXqcD4bb4zPkLdXdJPtUf1eAoMh/qgyThUjBs3n7BXvXMDg1Wdj0gq/sTnPLvXsfrSVPjiZvWN4h0JdID2NLnwYuKIiltIn+IbUa6OnyFfOEpqb5XJ7H7LK1mUKTlQ/9CFROxSQf3YQrR9UdtASIeyIZL53WgYgU31Yqy7MQaY1y0fGmHsFwpCK6qFZj1DNruKl/IR1lLx/Bg3z9sDcoBnHKnzSzVels9EVlDOG6bW738ho269QAIrWQYBtznsvWKu5xZPuuj user@machine"
}
