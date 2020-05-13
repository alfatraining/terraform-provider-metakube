package gometakube

import (
	"context"
	"fmt"
	"net/http"
)

func clustersListPath(prj string) string {
	return fmt.Sprintf("/api/v1/projects/%s/clusters", prj)
}

func createClusterPath(prj, dc string) string {
	return fmt.Sprintf("/api/v1/projects/%s/dc/%s/clusters", prj, dc)
}

func clusterResourcePath(prj, dc, clusterID string) string {
	return fmt.Sprintf("/api/v1/projects/%s/dc/%s/clusters/%s", prj, dc, clusterID)
}

func clusterResourceHealthPath(prj, dc, clusterID string) string {
	return fmt.Sprintf("/api/v1/projects/%s/dc/%s/clusters/%s/health", prj, dc, clusterID)
}

func clusterUpgradesPath(prj, dc, clusterID string) string {
	return fmt.Sprintf("/api/v1/projects/%s/dc/%s/clusters/%s/upgrades", prj, dc, clusterID)
}

// ClustersService handles comminication with cluster related endpoints.
type ClustersService struct {
	client *Client
}

// List returns list of clusters in project.
func (svc *ClustersService) List(ctx context.Context, project string) ([]Cluster, *http.Response, error) {
	ret := make([]Cluster, 0)
	resp, err := svc.client.resourceList(ctx, clustersListPath(project), &ret)
	return ret, resp, err
}

// CreateClusterRequest used to create a cluster.
type CreateClusterRequest struct {
	Cluster        Cluster        `json:"cluster"`
	NodeDeployment NodeDeployment `json:"nodeDeployment"`
}

// Create creates a cluster.
func (svc *ClustersService) Create(ctx context.Context, prj, dc string, create *CreateClusterRequest) (*Cluster, *http.Response, error) {
	ret := new(Cluster)
	resp, err := svc.client.resourceCreate(ctx, createClusterPath(prj, dc), create, ret)
	return ret, resp, err
}

// Delete deletes cluster.
func (svc *ClustersService) Delete(ctx context.Context, prj, dc, clusterID string) (*http.Response, error) {
	path := clusterResourcePath(prj, dc, clusterID)
	return svc.client.resourceDelete(ctx, path)
}

// Get returns cluster details.
func (svc *ClustersService) Get(ctx context.Context, prj, dc, clusterID string) (*Cluster, *http.Response, error) {
	path := clusterResourcePath(prj, dc, clusterID)
	ret := new(Cluster)
	resp, err := svc.client.resourceGet(ctx, path, ret)
	return ret, resp, err
}

// PatchClusterRequest specifies fields to be changed on cluster.
// Only patchable fields are specified.
type PatchClusterRequest struct {
	Name   string                   `json:"name,omitempty"`
	Labels map[string]string        `json:"labels,omitempty"`
	Spec   *PatchClusterRequestSpec `json:"spec,omitempty"`
}

// PatchClusterRequestSpec fields allowed to change on cluster spec in place.
type PatchClusterRequestSpec struct {
	Version      string                   `json:"version,omitempty"`
	AuditLogging *ClusterSpecAuditLogging `json:"auditLogging,omitempty"`
}

// Patch updates cluster.
func (svc *ClustersService) Patch(ctx context.Context, prj, dc, clusterID string, patch *PatchClusterRequest) (*Cluster, *http.Response, error) {
	path := clusterResourcePath(prj, dc, clusterID)
	ret := new(Cluster)
	resp, err := svc.client.resourcePatch(ctx, path, patch, ret)
	return ret, resp, err
}

// Healthy returns whether all cluster components are ready.
func (h *ClusterHealth) Healthy() bool {
	return (h.APIServer & h.CloudProviderInfrastructure &
		h.Controller & h.Etcd & h.MachineController &
		h.Scheduler & h.UserClusterControllerManager) == 1
}

// Health requests cluster's helth.
func (svc *ClustersService) Health(ctx context.Context, prj, dc, id string) (*ClusterHealth, *http.Response, error) {
	path := clusterResourceHealthPath(prj, dc, id)
	ret := new(ClusterHealth)
	req, err := svc.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	resp, err := svc.client.Do(ctx, req, ret)
	return ret, resp, err
}

// ClusterUpgrade is a cluster version possible to upgrade into.
type ClusterUpgrade struct {
	Version string `json:"version"`
	Defailt bool   `json:"default"`
}

// Upgrades lists all versions which don't result in automatic updates.
func (svc *ClustersService) Upgrades(ctx context.Context) ([]ClusterUpgrade, *http.Response, error) {
	ret := make([]ClusterUpgrade, 0)
	resp, err := svc.client.resourceList(ctx, "/api/v1/upgrades/cluster", &ret)
	return ret, resp, err
}

// ClusterUpgrades returns upgrades for a cluster.
func (svc *ClustersService) ClusterUpgrades(ctx context.Context, prj, dc, id string) ([]ClusterUpgrade, *http.Response, error) {
	ret := make([]ClusterUpgrade, 0)
	path := clusterUpgradesPath(prj, dc, id)
	resp, err := svc.client.resourceList(ctx, path, &ret)
	return ret, resp, err
}
