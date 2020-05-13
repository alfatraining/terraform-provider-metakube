package gometakube

import "time"

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
	MaxReplicas uint                       `json:"maxReplicas"`
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
	// TODO(furkhat): DistSize (what are the fields).
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
