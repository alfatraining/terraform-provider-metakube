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
func (svc *OpenstackService) Images(ctx context.Context, dc, domain, username, password string) ([]Image, error) {
	ret := make([]Image, 0)
	if err := svc.listResources(ctx, imagesListPath, dc, domain, username, password, &ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// Tenants return list of tenants.
func (svc *OpenstackService) Tenants(ctx context.Context, dc, domain, username, password string) ([]Tenant, error) {
	ret := make([]Tenant, 0)
	if err := svc.listResources(ctx, tenantsListPath, dc, domain, username, password, &ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func (svc *OpenstackService) listResources(ctx context.Context, path, dc, domain, username, password string, ret interface{}) error {
	req, err := svc.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return err
	}
	req.Header.Set("DatacenterName", dc)
	req.Header.Set("Username", username)
	req.Header.Set("Password", password)
	req.Header.Set("Domain", domain)
	if resp, err := svc.client.Do(ctx, req, &ret); err != nil {
		return err
	} else if resp.StatusCode != http.StatusOK {
		return unexpectedResponseError(resp)
	}
	return nil
}
