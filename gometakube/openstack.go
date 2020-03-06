package gometakube

import (
	"context"
	"net/http"
	"time"
)

type Image struct {
	ID       string
	Created  *time.Time
	MinDisk  uint
	MinRAM   uint
	Name     string
	Progress uint
	Status   string
	Updated  *time.Time
	Metadata ImageMetadata
}

type ImageMetadata struct {
	CIJobID            string `json:"ci_job_id"`
	CIPipelineID       string `json:"ci_pipeline_id"`
	CPUArch            string `json:"cpu_arch"`
	DefaultSSHUsername string `json:"default_ssh_username"`
	Distribution       string `json:"distribution"`
	OSDistro           string `json:"os_distro"`
	OSType             string `json:"os_type"`
	OSVersion          string `json:"os_version"`
	SourceSHA56sum     string `json:"source_sha256sum"`
	SourceURL          string `json:"source_url"`
}

const (
	imagesListPath = "/api/v1/providers/openstack/images"
	tenantListPath = "/api/v1/providers/openstack/tenants"
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

type Tenant struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Tenants return list of tenants.
func (svc *OpenstackService) Tenants(ctx context.Context, dc, domain, username, password string) ([]Tenant, error) {
	ret := make([]Tenant, 0)
	if err := svc.listResources(ctx, tenantListPath, dc, domain, username, password, &ret); err != nil {
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
