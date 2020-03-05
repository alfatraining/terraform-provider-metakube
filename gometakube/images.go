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

const imagesListPath = "/api/v1/providers/openstack/images"

// ImagesService handles communication with image related endpoints.
type ImagesService struct {
	client *Client
}

// List returns list of images.
func (svc *ImagesService) List(ctx context.Context, dc, domain, username, password string) ([]Image, error) {
	req, err := svc.client.NewRequest(http.MethodGet, imagesListPath, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("DatacenterName", dc)
	req.Header.Set("Username", username)
	req.Header.Set("Password", password)
	req.Header.Set("Domain", domain)
	ret := make([]Image, 0)
	if resp, err := svc.client.Do(ctx, req, &ret); err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, unexpectedResponseError(resp)
	}
	return ret, nil
}
