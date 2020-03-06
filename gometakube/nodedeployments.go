package gometakube

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type NodeDeployment struct {
	ID                string                `json:"id,omitempty"`
	Name              string                `json:"name"`
	CreationTimestamp *time.Time            `json:"creationTimestamp,omitempty"`
	Spec              NodeDeploymentSpec    `json:"spec"`
	Status            *NodeDeploymentStatus `json:"status,omitempty"`
}

type NodeDeploymentSpec struct {
	Replicas    uint                       `json:"replicas"`
	MinReplicas uint                       `json:"minReplicas"`
	MacReplicas uint                       `json:"maxReplicas"`
	Template    NodeDeploymentSpecTemplate `json:"template"`
	Paused      bool                       `json:"paused,omitempty"`
}

type NodeDeploymentSpecTemplate struct {
	Cloud           NodeDeploymentSpecTemplateCloud    `json:"cloud"`
	OperatingSystem NodeDeploymentSpecTemplateOS       `json:"operatingSystem"`
	Versions        NodeDeploymentSpecTemplateVersions `json:"versions,omitempty"`
	Labels          map[string]string                  `json:"labels,omitempty"`
}

type NodeDeploymentSpecTemplateCloud struct {
	Openstack NodeDeploymentSpecTemplateCloudOpenstack `json:"openstack"`
}

type NodeDeploymentSpecTemplateCloudOpenstack struct {
	Flavor        string            `json:"flavor"`
	Image         string            `json:"image"`
	Tags          map[string]string `json:"tags"`
	UseFloatingIP bool              `json:"useFloatingIP"`
	// TODO: what is format for DistSize?
}

type NodeDeploymentSpecTemplateOS struct {
	CentOS         *NodeDeploymentSpecTemplateOSOptions `json:"centos,omitempty"`
	Ubuntu         *NodeDeploymentSpecTemplateOSOptions `json:"ubuntu,omitempty"`
	ContainerLinux *NodeDeploymentSpecTemplateOSOptions `json:"containerLinux,omitempty"`
}

type NodeDeploymentSpecTemplateOSOptions struct {
	DisableAutoUpdate *bool `json:"disableAutoUpdate,omitempty"`
	DistUpgradeOnBoot *bool `json:"distUpgradeOnBoot,omitempty"`
}

type NodeDeploymentSpecTemplateVersions struct {
	Kubelet string `json:"kubelet"`
}

type NodeDeploymentStatus struct {
	ObservedGeneration uint `json:"observedGeneration"`
	Replicas           uint `json:"replicas"`
	UpdatedReplicas    uint `json:"updatedReplicas"`
	ReadyReplicas      uint `json:"readyReplicas"`
	AvailableReplicas  uint `json:"availableReplicas"`
}

// NodeDeploymentsService handles communication with node deployment related endpoints.
type NodeDeploymentsService struct {
	client *Client
}

func nodeDeploymentsListPath(prj, dc, cls string) string {
	return fmt.Sprintf("/api/v1/projects/%s/dc/%s/clusters/%s/nodedeployments", prj, dc, cls)
}

func nodeDeploymentsPatchPath(prj, dc, cls, id string) string {
	return fmt.Sprintf("/api/v1/projects/%s/dc/%s/clusters/%s/nodedeployments/%s", prj, dc, cls, id)
}

// List returns list of nodeDeployments.
func (svc *NodeDeploymentsService) List(ctx context.Context, prj, dc, cls string) ([]NodeDeployment, error) {
	path := nodeDeploymentsListPath(prj, dc, cls)
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
	path := nodeDeploymentsPatchPath(prj, dc, cls, id)
	ret := new(NodeDeployment)
	if err := svc.client.resourcePatch(ctx, path, patch, ret); err != nil {
		return nil, err
	}
	return ret, nil
}
