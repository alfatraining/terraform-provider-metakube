package gometakube

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type NodeDeployment struct {
	ID                string                `json:"id"`
	Name              string                `json:"name"`
	CreationTimestamp *time.Time            `json:"creationTimestamp,omitempty"`
	Spec              NodeDeploymentSpec    `json:"spec"`
	Status            *NodeDeploymentStatus `json:"status,omitempty"`
}

type NodeDeploymentSpec struct {
	Replicas uint                       `json:"replicas"`
	Template NodeDeploymentSpecTemplate `json:"template"`
	Paused   bool                       `json:"paused"`
}

type NodeDeploymentSpecTemplate struct {
	Cloud           NodeDeploymentSpecTemplateCloud    `json:"cloud"`
	OperatingSystem NodeDeploymentSpecTemplateOS       `json:"operatingSystem"`
	Versions        NodeDeploymentSpecTemplateVersions `json:"versions"`
	Labels          map[string]string                  `json:"labels"`
}

type NodeDeploymentSpecTemplateCloud struct {
	Openstack NodeDeploymentSpecTemplateCloudOpenstack `json:"openstack"`
}

type NodeDeploymentSpecTemplateCloudOpenstack struct {
	Flavor        string            `json:"flavor"`
	Image         string            `json:"image"`
	Tags          map[string]string `json:"tags"`
	UseFloatingIP bool              `json:"useFloatingIP"`
	// TODO: what foramt for DistSize?
}

type NodeDeploymentSpecTemplateOS struct {
	CentOS         *NodeDeploymentSpecTemplateOSOptions `json:"centos,omitempty"`
	Ubuntu         *NodeDeploymentSpecTemplateOSOptions `json:"ubuntu,omitempty"`
	ContainerLinux *NodeDeploymentSpecTemplateOSOptions `json:"containerLinux,omitempty"`
}

type NodeDeploymentSpecTemplateOSOptions struct {
	DisableAutoUpdate bool `json:"disableAutoUpdate"`
	DistUpgradeOnBoot bool `json:"distUpgradeOnBoot"`
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

type NodeDeploymentsService struct {
	client *Client
}

func nodeDeploymentsListURL(prj, dc, cls string) string {
	return fmt.Sprintf("/api/v1/projects/%s/dc/%s/clusters/%s/nodedeployments", prj, dc, cls)
}

func (svc *NodeDeploymentsService) List(ctx context.Context, prj, dc, cls string) ([]NodeDeployment, error) {
	url := nodeDeploymentsListURL(prj, dc, cls)
	req, err := svc.client.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	ret := make([]NodeDeployment, 0)
	if err := svc.client.Do(ctx, req, &ret); err != nil {
		return nil, err
	}
	return ret, nil
}
