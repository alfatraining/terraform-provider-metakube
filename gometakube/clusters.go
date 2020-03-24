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
func (svc *ClustersService) List(ctx context.Context, project string) ([]Cluster, error) {
	ret := make([]Cluster, 0)
	if resp, err := svc.client.resourceList(ctx, clustersListPath(project), &ret); err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, unexpectedResponseError(resp)
	}
	return ret, nil
}

// CreateClusterRequest used to create a cluster.
type CreateClusterRequest struct {
	Cluster        Cluster        `json:"cluster"`
	NodeDeployment NodeDeployment `json:"nodeDeployment"`
}

// Create creates a cluster.
func (svc *ClustersService) Create(ctx context.Context, prj, dc string, create *CreateClusterRequest) (*Cluster, error) {
	ret := new(Cluster)
	if err := svc.client.resourceCreate(ctx, createClusterPath(prj, dc), create, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// Delete deletes cluster.
func (svc *ClustersService) Delete(ctx context.Context, prj, dc, clusterID string) error {
	path := clusterResourcePath(prj, dc, clusterID)
	return svc.client.resourceDelete(ctx, path)
}

// Get returns cluster details.
func (svc *ClustersService) Get(ctx context.Context, prj, dc, clusterID string) (*Cluster, error) {
	path := clusterResourcePath(prj, dc, clusterID)
	ret := new(Cluster)
	if resp, err := svc.client.resourceGet(ctx, path, ret); err != nil {
		return nil, err
	} else if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	} else if resp.StatusCode != http.StatusOK {
		return nil, unexpectedResponseError(resp)
	}
	return ret, nil
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
func (svc *ClustersService) Patch(ctx context.Context, prj, dc, clusterID string, patch *PatchClusterRequest) (*Cluster, error) {
	path := clusterResourcePath(prj, dc, clusterID)
	ret := new(Cluster)
	if err := svc.client.resourcePatch(ctx, path, patch, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// Healthy returns whether all cluster components are ready.
func (h *ClusterHealth) Healthy() bool {
	return (h.APIServer & h.CloudProviderInfrastructure &
		h.Controller & h.Etcd & h.MachineController &
		h.Scheduler & h.UserClusterControllerManager) == 1
}

func (svc *ClustersService) Health(ctx context.Context, prj, dc, id string) (*ClusterHealth, error) {
	path := clusterResourceHealthPath(prj, dc, id)
	ret := new(ClusterHealth)
	req, err := svc.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	if resp, err := svc.client.Do(ctx, req, ret); err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, unexpectedResponseError(resp)
	}
	return ret, nil
}

// ClusterUpgrade is a cluster version possible to upgrade into.
type ClusterUpgrade struct {
	Version string `json:"version"`
	Defailt bool   `json:"default"`
}

// Upgrades lists all versions which don't result in automatic updates.
func (svc *ClustersService) Upgrades(ctx context.Context) ([]ClusterUpgrade, error) {
	ret := make([]ClusterUpgrade, 0)
	if resp, err := svc.client.resourceList(ctx, "/api/v1/upgrades/cluster", &ret); err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, unexpectedResponseError(resp)
	}
	return ret, nil
}

// ClusterUpgrades returns upgrades for a cluster.
func (svc *ClustersService) ClusterUpgrades(ctx context.Context, prj, dc, id string) ([]ClusterUpgrade, error) {
	ret := make([]ClusterUpgrade, 0)
	path := clusterUpgradesPath(prj, dc, id)
	if resp, err := svc.client.resourceList(ctx, path, &ret); err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, unexpectedResponseError(resp)
	}
	return ret, nil
}
