default: build

build:
	go build -o terraform-provider-metakube

test:
	go test ./... -v

testacc:
	TF_ACC=1 go test ./metakube -v $(TESTARGS) -timeout 120m

.PHONY: build test testacc
