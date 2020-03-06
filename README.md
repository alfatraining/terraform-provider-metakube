# Resources

* `metakube_project` metakube project
* `matekube_cluster` represents k8s cluster on openstack provider.


Example terraform file [./examples/main.tf](/examples/main.tf)

# Running

## Unit tests

Run tests:
```bash
go test -v ./...
```
## Acceptance tests 

IMPORTANT: this tests provision real resources.

Set required environment variables:
```
export METAKUBE_API_TOKEN=<token>
export ACC_PROVIDER_DC=<openstack datacenter name>
export ACC_TENANT=<tenant>
export ACC_PROVIDER_USERNAME=<username>
export ACC_PROVIDER_PASSWORD=<password>
```

Run
```bash
TF_ACC=1 go test -v ./...
```

## Manually

There is simple configuration file to use as base located at `./examples` directory. It is sinle project consisting of single cluster with single node deployment.

First, build provider.
```bash
go build -o terraform-provider-metakube
```

Init terraform (so it knows about metakube provider)
```bash
terraform init ./examples
```

Make desired changes to base config file [./examples/main.tf](/examples/main.tf).

Apply
```bash
terraform apply ./examples
```

Check resources created
list projects:
```bash
curl 'https://metakube.syseleven.de/api/v1/projects' -H "authorization: Bearer ${METAKUBE_API_TOKEN}" -H 'accept: application/json'
```
list clusters in the project (substitute your project id):
```bash
curl 'https://metakube.syseleven.de/api/v1/projects/<id>/clusters' -H "authorization: Bearer ${METAKUBE_API_TOKEN}" -H 'accept: application/json'
```
OR if you use your user account's bearer token, you can inspect everything on UI.

Destroy
```bash
terraform destroy ./examples
```
