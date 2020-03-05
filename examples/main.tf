provider "metakube" {
  // Do not forget to set METAKUBE_API_TOKEN environment variable.
}

resource "metakube_project" "my-project" {
  name = "my-project"

  labels = {
    "component" = "main"
  }
}

resource "metakube_cluster" "my-cluster" {
  project_id = metakube_project.my-project.id

  name    = "my-cluster"
  version = "1.17.3"
  dc      = "syseleven-cbk1"

  // openstack 
  tenant   = "syseleveneigenbedarf-syseleven-ext-spearce"
  username = "" // sensitive
  password = "" // sensitive

  node_pool {
    name     = "my-cluster-nodepool-one"
    replicas = 2

    flavor_type     = "Local Storage"
    flavor          = "l1.small"
    image           = "Ubuntu Bionic 18.04 (2020-03-05)"
    use_floating_ip = false
  }
}
