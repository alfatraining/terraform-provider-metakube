package gometakube

type Datacenter struct {
	Metadata DatacenterMetadata `json:"metadata"`
	Seed     bool               `json:"seed"`
	Spec     *DatacenterSpec    `json:"spec"`
}

type DatacenterMetadata struct {
	Annotations     map[string]string `json:"annotations"`
	Labels          map[string]string `json:"labels"`
	Name            string            `json:"name"`
	ResourceVersion string            `json:"resourceVersion"`
	UID             string            `json:"uid"`
}

type DatacenterSpec struct {
	AWS                 *DatacenterSpecAWS          `json:"aws,omitempty"`
	Azure               *DatacenterSpecAzure        `json:"azure,omitempty"`
	BringYourOwn        *DatacenterSpecBringYourOwn `json:"bringyourown"`
	Country             string                      `json:"country"`
	DigitalOcean        *DatacenterSpecDigitalOcean `json:"digitalocean,omitempty"`
	GCP                 *DatacenterSpecGCP          `json:"gcp,omitempty"`
	Hetzner             *DatacenterSpecHetzner      `json:"hetzner,omitempty"`
	Kubevirt            *DatacenterSpecKubevirt     `json:"kubevirt,omitempty"`
	Location            string                      `json:"location"`
	Openstack           *DatacenterSpecOpenstack    `json:"openstack,omitempty"`
	Packet              *DatacenterSpecPacket       `json:"packet,omitempty"`
	Provider            string                      `json:"provider"`
	RequiredEmailDomain string                      `json:"requiredEmailDomain"`
	Seed                string                      `json:"seed"`
	Vsphare             *DatacenterSpecVsphare      `json:"vsphere,omitempty"`
}

type DatacenterSpecAWS struct {
	Region string `json:"region"`
}

type DatacenterSpecAzure struct {
	Location string `json:"location"`
}

type DatacenterSpecBringYourOwn struct {
}

type DatacenterSpecDigitalOcean struct {
	Region string `json:"region"`
}

type DatacenterSpecGCP struct {
	Region       string   `json:"region"`
	Regional     bool     `json:"regional"`
	ZoneSuffixes []string `json:"zone_suffixes"`
}

type DatacenterSpecHetzner struct {
	Datacenter string `json:"datacenter"`
	Location   string `json:"location"`
}

type DatacenterSpecKubevirt struct{}

type DatacenterSpecOpenstack struct {
	AuthURL           string            `json:"auth_url"`
	AvailabilityZone  string            `json:"availability_zone"`
	EnforceFloatingIP bool              `json:"enforce_floating_ip"`
	Images            map[string]string `json:"images"`
	Region            string            `json:"region"`
}

type DatacenterSpecPacket struct {
	Facilities []string `json:"facilities"`
}

type DatacenterSpecVsphare struct {
	Cluster    string            `json:"cluster"`
	Datacenter string            `json:"datacenter"`
	DataStore  string            `json:"datastore"`
	Endpoint   string            `json:"endpoint"`
	Templates  map[string]string `json:"templates"`
}
