package gometakube

import (
	"context"
	"net/http"
)

const (
	imagesListPath  = "/api/v1/providers/openstack/images"
	tenantsListPath = "/api/v1/providers/openstack/tenants"
)

// OpenstackService handles communication with image related endpoints.
type OpenstackService struct {
	client *Client
}

// Images returns list of images.
func (svc *OpenstackService) Images(ctx context.Context, dc, domain, username, password string) ([]Image, *http.Response, error) {
	ret := make([]Image, 0)
	resp, err := svc.listResources(ctx, imagesListPath, dc, domain, username, password, &ret)
	return ret, resp, err
}

// Tenants return list of tenants.
func (svc *OpenstackService) Tenants(ctx context.Context, dc, domain, username, password string) ([]Tenant, *http.Response, error) {
	ret := make([]Tenant, 0)
	resp, err := svc.listResources(ctx, tenantsListPath, dc, domain, username, password, &ret)
	return ret, resp, err
}

func (svc *OpenstackService) listResources(ctx context.Context, path, dc, domain, username, password string, ret interface{}) (*http.Response, error) {
	req, err := svc.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("DatacenterName", dc)
	req.Header.Set("Username", username)
	req.Header.Set("Password", password)
	req.Header.Set("Domain", domain)
	return svc.client.Do(ctx, req, &ret)
}
