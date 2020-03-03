package gometakube

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Cluster struct {
	CreationTimestamp *time.Time        `json:"creationTimestamp"`
	DeletionTimestamp *time.Time        `json:"deletionTimestamp"`
	Credential        string            `json:"credential"`
	ID                string            `json:"id"`
	Labels            map[string]string `json:"labels"`
	Name              string            `json:"name"`
	Spec              *ClusterSpec      `json:"spec"`
	Status            *ClusterStatus    `json:"status"`
	Type              string            `json:"type"`
}

type ClusterStatus struct {
	URL     string `json:"url"`
	Version string `json:"version"`
}

type ClusterSpec struct {
	AuditLogging                        ClusterSpecAuditLogging     `json:"auditLogging"`
	Cloud                               *ClusterSpecCloud           `json:"cloud"`
	ClusterNetwork                      *ClusterSpecClusterNetwork  `json:"clusterNetwork"`
	MachineNetworks                     []ClusterSpecMachineNetwork `json:"machineNetworks"`
	OIDC                                ClusterSpecOIDC             `json:"oidc"`
	Openshift                           *ClusterSpecOpenShift       `json:"openshift"`
	Sys11Auth                           ClusterSpecSys11Auth        `json:"sys11auth"`
	UpdateWindow                        *ClusterSpecUpdateWindow    `json:"updateWindow,omitempty"`
	UsePodSecurityPolicyAdmissionPlugin bool                        `json:"usePodSecurityPolicyAdmissionPlugin"`
	Version                             string
}

type ClusterSpecAuditLogging struct {
	Enabled bool `json:"enabled,omitempty"`
}

type ClusterSpecCloud struct {
	AWS          *ClusterSpecCloudAWS          `json:"aws,omitempty"`
	Azure        *ClusterSpecCloudAzure        `json:"azure,omitempty"`
	BringYourOwn *ClusterSpecCloudBringYourOwn `json:"bringyourown"`
	DataCenter   string                        `json:"dc"`
	DigitalOcean *ClusterSpecCloudDigitalOcean `json:"digitalocean"`
	Fake         *ClusterSpecCloudFake         `json:"fake"`
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
	RouterID             string                               `json:"routerID"`
	SecurityGroups       string                               `json:"securityGroups"`
	SubnetCIDR           string                               `json:"subnetCIDR"`
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

// ClustersService handles comminication with cluster related endpoints.
type ClustersService struct {
	client *Client
}

// List returns list of clusters in project.
func (svc *ClustersService) List(ctx context.Context, project string) ([]Cluster, error) {
	url := fmt.Sprintf(clusterListURLTpl, project)
	ret := make([]Cluster, 0)
	if err := svc.client.serviceList(ctx, url, &ret); err != nil {
		return nil, err
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
	if err := svc.client.Do(ctx, req, &ret); err != nil {
		return nil, err
	}
	return ret, nil
}
