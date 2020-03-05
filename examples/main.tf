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
  # If you use API Account's token to create a project and cluster for it
  # it WILL FAIL
  # This is due to ownership issues with API Accounts and user account that 
  # it belongs to. So referencing as below wont work:
  # project_id = metakube_project.my-project.id
  # 
  # To create a cluster you should either use token that belongs to user account
  # itself and not to api account or create a project on UI and set explicit 
  # values for project_id:
  # project_id = "explicit-project-id"

  project_id = metakube_project.my-project.id

  name    = "my-cluster"
  version = "1.17.3"         // k8s version
  dc      = "syseleven-dbl1" // openstack datacenter

  // openstack 
  tenant            = "syseleveneigenbedarf-syseleven-ext-spearce"
  provider_username = "" // sensitive
  provider_password = "" // sensitive

  nodepool {
    name     = "my-cluster-nodepool-one"
    replicas = 2

    flavor_type     = "Local Storage"
    flavor          = "l1.small"
    image           = "Ubuntu Bionic 18.04 (2020-03-03)"
    use_floating_ip = false
  }
}
