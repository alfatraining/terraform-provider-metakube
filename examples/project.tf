provider "metakube" {
}

resource "metakube_project" "my-project" {
  name = "my-project"

  labels = {
    additionalProp1 = "string"
    additionalProp2 = "string"
    additionalProp3 = "string"
  }
}
