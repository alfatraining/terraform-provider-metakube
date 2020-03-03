package gometakube

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Cluster struct {
	CreationTimestamp *time.Time        `json:"creationTimestamp,omitempty"`
	DeletionTimestamp *time.Time        `json:"deletionTimestamp,omitempty"`
	Credential        string            `json:"credential"`
	ID                string            `json:"id,omitempty"`
	Labels            map[string]string `json:"labels,omitempty"`
	Name              string            `json:"name"`
	Spec              *ClusterSpec      `json:"spec"`
	Status            *ClusterStatus    `json:"status,omitempty"`
	Type              string            `json:"type"`
	SSHKeys           []string          `json:"sshKeys"`
}

type ClusterStatus struct {
	URL     string `json:"url"`
	Version string `json:"version"`
}

type ClusterSpec struct {
	AuditLogging                        ClusterSpecAuditLogging     `json:"auditLogging"`
	Cloud                               *ClusterSpecCloud           `json:"cloud"`
	ClusterNetwork                      *ClusterSpecClusterNetwork  `json:"clusterNetwork,omitempty"`
	MachineNetworks                     []ClusterSpecMachineNetwork `json:"machineNetworks"`
	OIDC                                *ClusterSpecOIDC            `json:"oidc,omitempty"`
	Openshift                           *ClusterSpecOpenShift       `json:"openshift,omitempty"`
	Sys11Auth                           *ClusterSpecSys11Auth       `json:"sys11auth,omitempty"`
	UpdateWindow                        *ClusterSpecUpdateWindow    `json:"updateWindow,omitempty"`
	UsePodSecurityPolicyAdmissionPlugin *bool                       `json:"usePodSecurityPolicyAdmissionPlugin"`
	Version                             string                      `json:"version"`
}

type ClusterSpecAuditLogging struct {
	Enabled bool `json:"enabled"`
}

type ClusterSpecCloud struct {
	AWS          *ClusterSpecCloudAWS          `json:"aws,omitempty"`
	Azure        *ClusterSpecCloudAzure        `json:"azure,omitempty"`
	BringYourOwn *ClusterSpecCloudBringYourOwn `json:"bringyourown,omitempty"`
	DataCenter   string                        `json:"dc"`
	DigitalOcean *ClusterSpecCloudDigitalOcean `json:"digitalocean,omitempty"`
	Fake         *ClusterSpecCloudFake         `json:"fake,omitempty"`
	GCP          *ClusterSpecCloudGCP          `json:"gcp,omitempty"`
	Hetzner      *ClusterSpecCloudHetzner      `json:"hetzner,omitempty"`
	Kubevirt     *ClusterSpecCloudKubevirt     `json:"kubevirt,omitempty"`
	OpenStack    *ClusterSpecCloudOpenstack    `json:"openstack,omitempty"`
	Packet       *ClusterSpecCloudPacket       `json:"packet,omitempty"`
	Vsphere      *ClusterSpecCloudVsphere      `json:"vsphere,omitempty"`
}

type ClusterSpecCloudAWS struct {
	AccessKeyId            string                               `json:"accessKeyId"`
	availabilityZone       string                               `json:"availabilityZone"`
	credentialsReference   *ClusterSpecCloudCredentialReference `json:"credentialsReference"`
	instanceProfileName    string                               `json:"instanceProfileName"`
	openstackBillingTenant string                               `json:"openstackBillingTenant"`
	roleARN                string                               `json:"roleARN"`
	roleName               string                               `json:"roleName"`
	routeTableId           string                               `json:"routeTableId"`
	secretAccessKey        string                               `json:"secretAccessKey"`
	securityGroupID        string                               `json:"securityGroupID"`
	vpcId                  string                               `json:"vpcId"`
}

type ClusterSpecCloudAzure struct {
	AvailabilitySet        string                               `json:"availabilitySet"`
	ClientID               string                               `json:"clientID"`
	ClientSecret           string                               `json:"clientSecret"`
	CredentialsReference   *ClusterSpecCloudCredentialReference `json:"credentialsReference"`
	OpenstackBillingTenant string                               `json:"openstackBillingTenant"`
	ResourceGroup          string                               `json:"resourceGroup"`
	RouteTable             string                               `json:"routeTable"`
	SecurityGroup          string                               `json:"securityGroup"`
	Subnet                 string                               `json:"subnet"`
	SubscriptionID         string                               `json:"subscriptionID"`
	TenantID               string                               `json:"tenantID"`
	VNet                   string                               `json:"vnet"`
}

type ClusterSpecCloudBringYourOwn struct{}

type ClusterSpecCloudDigitalOcean struct {
	CredentialsReference *ClusterSpecCloudCredentialReference `json:"credentialsReference"`
	Token                string                               `json:"token"`
}

type ClusterSpecCloudFake struct {
	Token string `json:"token"`
}

type ClusterSpecCloudGCP struct {
	CredentialsReference *ClusterSpecCloudCredentialReference `json:"credentialsReference"`
	Network              string                               `json:"network"`
	ServiceAccount       string                               `json:"serviceAccount"`
	Subnetwork           string                               `json:"subnetwork"`
}

type ClusterSpecCloudHetzner struct {
	CredentialsReference *ClusterSpecCloudCredentialReference `json:"credentialsReference"`
	Token                string                               `json:"token"`
}

type ClusterSpecCloudKubevirt struct {
	CredentialsReference *ClusterSpecCloudCredentialReference `json:"credentialsReference"`
	Kubeconfig           string                               `json:"kubeconfig"`
}

type ClusterSpecCloudOpenstack struct {
	CredentialsReference *ClusterSpecCloudCredentialReference `json:"credentialsReference"`
	Domain               string                               `json:"domain"`
	FloatingIPPool       string                               `json:"floatingIpPool"`
	Network              string                               `json:"network"`
	Password             string                               `json:"password"`
	RouterID             string                               `json:"routerID,omitempty"`
	SecurityGroups       string                               `json:"securityGroups"`
	SubnetCIDR           string                               `json:"subnetCIDR,omitempty"`
	SubnetID             string                               `json:"subnetID"`
	Tenant               string                               `json:"tenant"`
	TenantID             string                               `json:"tenantID"`
	Username             string                               `json:"username"`
}

type ClusterSpecCloudPacket struct {
	ApiKey               string                               `json:"apiKey"`
	BillingCycle         string                               `json:"billingCycle"`
	CredentialsReference *ClusterSpecCloudCredentialReference `json:"credentialsReference"`
	ProjectID            string                               `json:"projectID"`
}

type ClusterSpecCloudVsphere struct {
	CredentialsReference *ClusterSpecCloudCredentialReference        `json:"credentialsReference"`
	Folder               string                                      `json:"folder"`
	InfraManagementUser  *ClusterSpecCloudVsphereInfraManagementUser `json:"infraManagementUser"`
	Password             string                                      `json:"password"`
	Username             string                                      `json:"username"`
	VMNetName            string                                      `json:"vmNetName"`
}

type ClusterSpecCloudVsphereInfraManagementUser struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type ClusterSpecCloudCredentialReference struct {
	ApiVersion      string `json:"apiVersion"`
	FieldPath       string `json:"fieldPath"`
	Key             string `json:"key"`
	Kind            string `json:"kind"`
	Name            string `json:"name"`
	Namespace       string `json:"namespace"`
	ResourceVersion string `json:"resourceVersion"`
	UID             string `json:"uid"`
}

type ClusterSpecClusterNetwork struct {
	DNSDomain string                             `json:"dnsDomain"`
	Pods      *ClusterSpecClusterNetworkPods     `json:"pods"`
	ProxyMode string                             `json:"proxyMode"`
	Services  *ClusterSpecClusterNetworkServices `json:"services"`
}

type ClusterSpecClusterNetworkPods struct {
	CIDRBlocks []string `json:"cidrBlocks"`
}

type ClusterSpecClusterNetworkServices struct {
	CIDRBlocks []string `json:"cidrBlocks"`
}

type ClusterSpecMachineNetwork struct {
	CIDR      string   `json:"cidr"`
	NSServers []string `json:"dnsServers"`
	Gateway   string   `json:"gateway"`
}

type ClusterSpecOIDC struct {
	ClientId      string `json:"clientId,omitempty"`
	lientSecret   string `json:"clientSecret,omitempty"`
	ExtraScopes   string `json:"extraScopes,omitempty"`
	GroupsClaim   string `json:"groupsClaim,omitempty"`
	IssuerUrl     string `json:"issuerUrl,omitempty"`
	RequiredClaim string `json:"requiredClaim,omitempty"`
	UsernameClaim string `json:"usernameClaim,omitempty"`
}

type ClusterSpecOpenShift struct {
	ImagePullSecret string `json:"imagePullSecret"`
}

type ClusterSpecSys11Auth struct {
	Realm string `json:"sys11auth"`
}

type ClusterSpecUpdateWindow struct {
	Length string `json:"length,omitempty"`
	Start  string `json:"start,omitempty"`
}

const (
	clusterListURLTpl = "/api/v1/projects/%s/clusters"
)

func createClusterPath(prj, dc string) string {
	return fmt.Sprintf("/api/v1/projects/%s/dc/%s/clusters", prj, dc)
}

func clusterResourcePath(prj, dc, clusterID string) string {
	return fmt.Sprintf("/api/v1/projects/%s/dc/%s/clusters/%s", prj, dc, clusterID)
}

// ClustersService handles comminication with cluster related endpoints.
type ClustersService struct {
	client *Client
}

// List returns list of clusters in project.
func (svc *ClustersService) List(ctx context.Context, project string) ([]Cluster, error) {
	url := fmt.Sprintf(clusterListURLTpl, project)
	ret := make([]Cluster, 0)
	if resp, err := svc.client.serviceList(ctx, url, &ret); err != nil {
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
	req, err := svc.client.NewRequest(http.MethodPost, createClusterPath(prj, dc), create)
	if err != nil {
		return nil, err
	}
	ret := new(Cluster)
	if resp, err := svc.client.Do(ctx, req, &ret); err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusCreated {
		return nil, unexpectedResponseError(resp)
	}
	return ret, nil
}

// Delete deletes cluster.
func (svc *ClustersService) Delete(ctx context.Context, prj, dc, clusterID string) error {
	url := clusterResourcePath(prj, dc, clusterID)
	if resp, err := svc.client.resourceDelete(ctx, url); err != nil {
		return fmt.Errorf("could not delete cluster: %v", err)
	} else if resp.StatusCode != http.StatusOK {
		return unexpectedResponseError(resp)
	}
	return nil
}

// Get returns cluster details.
func (svc *ClustersService) Get(ctx context.Context, prj, dc, clusterID string) (*Cluster, error) {
	url := clusterResourcePath(prj, dc, clusterID)
	ret := new(Cluster)
	if resp, err := svc.client.resourceGet(ctx, url, ret); err != nil {
		return nil, err
	} else if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	} else if resp.StatusCode != http.StatusOK {
		return nil, unexpectedResponseError(resp)
	}
	return ret, nil
}
