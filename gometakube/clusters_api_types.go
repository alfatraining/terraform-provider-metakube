package gometakube

import "time"

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
	AvailabilityZone       string                               `json:"availabilityZone"`
	CredentialsReference   *ClusterSpecCloudCredentialReference `json:"credentialsReference"`
	InstanceProfileName    string                               `json:"instanceProfileName"`
	OpenstackBillingTenant string                               `json:"openstackBillingTenant"`
	RoleARN                string                               `json:"roleARN"`
	RoleName               string                               `json:"roleName"`
	RouteTableId           string                               `json:"routeTableId"`
	SecretAccessKey        string                               `json:"secretAccessKey"`
	SecurityGroupID        string                               `json:"securityGroupID"`
	VPCId                  string                               `json:"vpcId"`
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

type ClusterHealth struct {
	APIServer                    uint8 `json:"apiserver"`
	CloudProviderInfrastructure  uint8 `json:"cloudProviderInfrastructure"`
	Controller                   uint8 `json:"controller"`
	Etcd                         uint8 `json:"etcd"`
	MachineController            uint8 `json:"machineController"`
	Scheduler                    uint8 `json:"scheduler"`
	UserClusterControllerManager uint8 `json:"userClusterControllerManager"`
}
