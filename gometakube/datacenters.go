package gometakube

import (
	"context"
	"net/http"
)

const (
	datacentersPath = "/api/v1/dc"
)

// DatacentersService handles communication with datacenters related methods.
type DatacentersService struct {
	client *Client
}

// List requests all datacenters.
func (svc *DatacentersService) List(ctx context.Context) ([]Datacenter, *http.Response, error) {
	ret := make([]Datacenter, 0)
	resp, err := svc.client.resourceList(ctx, datacentersPath, &ret)
	return ret, resp, err
}

// Get returns detailed info on datacenter.
func (svc DatacentersService) Get(ctx context.Context, dc string) (*Datacenter, *http.Response, error) {
	url := datacentersPath + "/" + dc
	ret := new(Datacenter)
	resp, err := svc.client.resourceGet(ctx, url, ret)
	return ret, resp, err
}
