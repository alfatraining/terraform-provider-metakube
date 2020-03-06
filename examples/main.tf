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

  name          = "my-cluster"     // has in-place update
  version       = "1.17.3"         // k8s version, change forces new
  dc            = "syseleven-dbl1" // openstack datacenter, change forces new
  audit_logging = true             // has in-place update

  // openstack 
  tenant            = "" // change forces new
  provider_username = "" // sensitive, not persisted in tfstate, change forces new
  provider_password = "" // sensitive, not persisted in tfstate, change forces new

  // clusters node deployment
  nodedepl {
    name     = "my-cluster-nodedepl-one" // change forces new
    replicas = 2                         // has in-place update

    flavor          = "l1.small"                  // has in-place update
    image           = "Rescue Ubuntu 18.04 sys11" // has in-place update
    use_floating_ip = false                       // has in-place update
  }
}
