package gometakube

import (
	"context"
	"fmt"
	"net/http"
)

// NodeDeploymentsService handles communication with node deployment related endpoints.
type NodeDeploymentsService struct {
	client *Client
}

func nodeDeploymentsCreateListPath(prj, dc, cls string) string {
	return fmt.Sprintf("/api/v1/projects/%s/dc/%s/clusters/%s/nodedeployments", prj, dc, cls)
}

func nodeDeploymentResourcePath(prj, dc, cls, id string) string {
	return fmt.Sprintf("/api/v1/projects/%s/dc/%s/clusters/%s/nodedeployments/%s", prj, dc, cls, id)
}

// List returns list of nodeDeployments.
func (svc *NodeDeploymentsService) List(ctx context.Context, prj, dc, cls string) ([]NodeDeployment, error) {
	path := nodeDeploymentsCreateListPath(prj, dc, cls)
	ret := make([]NodeDeployment, 0)
	if resp, err := svc.client.resourceList(ctx, path, &ret); err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, unexpectedResponseError(resp)
	}
	return ret, nil
}

// NodeDeploymentsPatchRequest format of request to patch.
type NodeDeploymentsPatchRequest struct {
	Spec NodeDeploymentSpec `json:"spec"`
}

// Patch updates node deployments spec.
func (svc *NodeDeploymentsService) Patch(ctx context.Context, prj, dc, cls, id string, patch *NodeDeploymentsPatchRequest) (*NodeDeployment, error) {
	path := nodeDeploymentResourcePath(prj, dc, cls, id)
	ret := new(NodeDeployment)
	if err := svc.client.resourcePatch(ctx, path, patch, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// Create creates new node deployment.
func (svc *NodeDeploymentsService) Create(ctx context.Context, prj, dc, cls string, v *NodeDeployment) (*NodeDeployment, error) {
	path := nodeDeploymentsCreateListPath(prj, dc, cls)
	ret := new(NodeDeployment)
	if err := svc.client.resourceCreate(ctx, path, v, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// Get returns node deployments.
func (svc *NodeDeploymentsService) Get(ctx context.Context, prj, dc, cls, id string) (*NodeDeployment, error) {
	path := nodeDeploymentResourcePath(prj, dc, cls, id)
	ret := new(NodeDeployment)
	if resp, err := svc.client.resourceGet(ctx, path, ret); err != nil {
		return nil, err
	} else if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	} else if resp.StatusCode != http.StatusOK {
		return nil, unexpectedResponseError(resp)
	}
	return ret, nil
}

// Delete deletes node deployments.
func (svc *NodeDeploymentsService) Delete(ctx context.Context, prj, dc, cls, id string) error {
	path := nodeDeploymentResourcePath(prj, dc, cls, id)
	return svc.client.resourceDelete(ctx, path)
}
