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

func clusterNodesUpgradePath(prj, dc, id string) string {
	return fmt.Sprintf("/api/v1/projects/%s/dc/%s/clusters/%s/nodes/upgrades", prj, dc, id)
}

// List returns list of nodeDeployments.
func (svc *NodeDeploymentsService) List(ctx context.Context, prj, dc, cls string) ([]NodeDeployment, *http.Response, error) {
	path := nodeDeploymentsCreateListPath(prj, dc, cls)
	ret := make([]NodeDeployment, 0)
	resp, err := svc.client.resourceList(ctx, path, &ret)
	return ret, resp, err
}

// NodeDeploymentsPatchRequest format of request to patch.
type NodeDeploymentsPatchRequest struct {
	Spec NodeDeploymentSpec `json:"spec"`
}

// Patch updates node deployments spec.
func (svc *NodeDeploymentsService) Patch(ctx context.Context, prj, dc, cls, id string, patch *NodeDeploymentsPatchRequest) (*NodeDeployment, *http.Response, error) {
	path := nodeDeploymentResourcePath(prj, dc, cls, id)
	ret := new(NodeDeployment)
	resp, err := svc.client.resourcePatch(ctx, path, patch, ret)
	return ret, resp, err
}

// Create creates new node deployment.
func (svc *NodeDeploymentsService) Create(ctx context.Context, prj, dc, cls string, v *NodeDeployment) (*NodeDeployment, *http.Response, error) {
	path := nodeDeploymentsCreateListPath(prj, dc, cls)
	ret := new(NodeDeployment)
	resp, err := svc.client.resourceCreate(ctx, path, v, ret)
	return ret, resp, err
}

// Get returns node deployments.
func (svc *NodeDeploymentsService) Get(ctx context.Context, prj, dc, cls, id string) (*NodeDeployment, *http.Response, error) {
	path := nodeDeploymentResourcePath(prj, dc, cls, id)
	ret := new(NodeDeployment)
	resp, err := svc.client.resourceGet(ctx, path, ret)
	return ret, resp, err
}

// Delete deletes node deployments.
func (svc *NodeDeploymentsService) Delete(ctx context.Context, prj, dc, cls, id string) (*http.Response, error) {
	path := nodeDeploymentResourcePath(prj, dc, cls, id)
	return svc.client.resourceDelete(ctx, path)
}

// UpgradeNodesRequest is a body of a request to upgrade.
type UpgradeNodesRequest struct {
	Version string `json:"version,omitempty"`
}

// Upgrade upgrades nodes.
func (svc *NodeDeploymentsService) Upgrade(ctx context.Context, prj, dc, cls string, req *UpgradeNodesRequest) (*http.Response, error) {
	path := clusterNodesUpgradePath(prj, dc, cls)
	return svc.client.resourcePut(ctx, path, req, nil)
}
