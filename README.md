# Tests

## Unit
```bash
go test -v ./...
```
## Acceptance tests 
IMPORTANT: this tests provision real resources.

Set environment token variable
```bash
export METAKUBE_API_TOKEN=<YOUR TOKEN>
```

Run
```bash
TF_ACC=1 go test -v ./...
```

# Manually

There is simple project resource example at `./examples` directory to be used for manual testing.

First, build provider.
```bash
go build -o terraform-provider-metakube
```

Init terraform (so it knows about metakube provider)
```bash
terraform init ./examples
```

Apply
```bash
terraform apply ./examples
```

Check resources created
```bash
curl 'https://metakube.syseleven.de/api/v1/projects' -H "authorization: Bearer ${METAKUBE_API_TOKEN}" -H 'accept: application/json'
```

Destroy
```bash
terraform destroy ./examples
```
