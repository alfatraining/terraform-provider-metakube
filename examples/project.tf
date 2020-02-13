provider "metakube" {
  // Do not forget to set METAKUBE_API_TOKEN environment variable.
}

resource "metakube_project" "my-project" {
  name = "my-project"

  labels = {
    additionalProp1 = "string"
    additionalProp2 = "string"
    additionalProp3 = "string"
  }
}
