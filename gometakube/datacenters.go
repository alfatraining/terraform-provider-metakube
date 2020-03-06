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
func (svc *DatacentersService) List(ctx context.Context) ([]Datacenter, error) {
	ret := make([]Datacenter, 0)
	if resp, err := svc.client.resourceList(ctx, datacentersPath, &ret); err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, unexpectedResponseError(resp)
	}
	return ret, nil
}

// Get returns detailed info on datacenter.
func (svc DatacentersService) Get(ctx context.Context, dc string) (*Datacenter, error) {
	url := datacentersPath + "/" + dc
	ret := new(Datacenter)
	if resp, err := svc.client.resourceGet(ctx, url, ret); err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, unexpectedResponseError(resp)
	}
	return ret, nil
}
