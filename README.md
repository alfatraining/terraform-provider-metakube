# Resources

* `metakube_project` metakube project
* `matekube_cluster` represents k8s cluster on openstack provider
* `metakube_sshkey` ssh key to upload to cloud.


Example terraform file [./examples/main.tf](/examples/main.tf)

# Running

## Unit tests

Run tests:
```bash
make test
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
make testacc
```

## Manually

You can find example configuration file to use as base at [./examples/main.tf](/examples/main.tf). It is single project consisting of single cluster with single node deployment.

To compile the provider, run `make`. This will build the provider and put it in the current working directory.
```bash
make
```

Init terraform (so it knows about metakube provider)
```bash
terraform init ./examples
```

Make changes to base config file [./examples/main.tf](/examples/main.tf). Minimal changes would be setting values for `tenant`, `provider_username` and `provider_password` fields of a `matkube_cluster` resource which are left empty in the example file.

Apply
```bash
terraform apply ./examples
```

List created projects:
```bash
curl 'https://metakube.syseleven.de/api/v1/projects' -H "authorization: Bearer ${METAKUBE_API_TOKEN}" -H 'accept: application/json'
```
list clusters in the project (substitute your project id):
```bash
curl 'https://metakube.syseleven.de/api/v1/projects/<YOUR PROJECT ID>/clusters' -H "authorization: Bearer ${METAKUBE_API_TOKEN}" -H 'accept: application/json'
```
OR if you use your user account's bearer token, you can inspect everything on UI.

Now you can keep changing and keep applying.

To cleanup, run Destroy
```bash
terraform destroy ./examples
```
