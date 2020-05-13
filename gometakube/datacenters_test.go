package gometakube

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

const (
	datacenterJSON = `
	{
		"metadata": {
			"annotations": {
			"additionalProp1": "string",
			"additionalProp2": "string",
			"additionalProp3": "string"
			},
			"labels": {
			"additionalProp1": "string",
			"additionalProp2": "string",
			"additionalProp3": "string"
			},
			"name": "string",
			"resourceVersion": "string",
			"uid": "string"
		},
		"seed": true,
		"spec": {
			"aws": {
				"region": "string"
			},
			"azure": {
				"location": "string"
			},
			"bringyourown": {},
			"country": "string",
			"digitalocean": {
			"region": "string"
			},
			"gcp": {
				"region": "string",
				"regional": true,
				"zone_suffixes": [
					"string"
				]
			},
			"hetzner": {
				"datacenter": "string",
				"location": "string"
			},
			"kubevirt": {},
			"location": "string",
			"openstack": {
				"auth_url": "string",
				"availability_zone": "string",
				"enforce_floating_ip": true,
				"images": {
					"additionalProp1": "string",
					"additionalProp2": "string",
					"additionalProp3": "string"
				},
				"region": "string"
			},
			"packet": {
				"facilities": [
					"string"
				]
			},
			"provider": "string",
			"requiredEmailDomain": "string",
			"seed": "string",
			"vsphere": {
				"cluster": "string",
				"datacenter": "string",
				"datastore": "string",
				"endpoint": "string",
				"templates": {
					"additionalProp1": "string",
					"additionalProp2": "string",
					"additionalProp3": "string"
				}
			}
		}
	}
`
)

var datacenter = Datacenter{
	Metadata: DatacenterMetadata{
		Annotations: map[string]string{
			"additionalProp1": "string",
			"additionalProp2": "string",
			"additionalProp3": "string",
		},
		Labels: map[string]string{
			"additionalProp1": "string",
			"additionalProp2": "string",
			"additionalProp3": "string",
		},
		Name:            "string",
		ResourceVersion: "string",
		UID:             "string",
	},
	Seed: true,
	Spec: &DatacenterSpec{
		AWS: &DatacenterSpecAWS{
			Region: "string",
		},
		Azure: &DatacenterSpecAzure{
			Location: "string",
		},
		BringYourOwn: &DatacenterSpecBringYourOwn{},
		Country:      "string",
		DigitalOcean: &DatacenterSpecDigitalOcean{
			Region: "string",
		},
		GCP: &DatacenterSpecGCP{
			Region:       "string",
			Regional:     true,
			ZoneSuffixes: []string{"string"},
		},
		Hetzner: &DatacenterSpecHetzner{
			Datacenter: "string",
			Location:   "string",
		},
		Kubevirt: &DatacenterSpecKubevirt{},
		Location: "string",
		Openstack: &DatacenterSpecOpenstack{
			AuthURL:           "string",
			AvailabilityZone:  "string",
			EnforceFloatingIP: true,
			Images: map[string]string{
				"additionalProp1": "string",
				"additionalProp2": "string",
				"additionalProp3": "string",
			},
			Region: "string",
		},
		Packet: &DatacenterSpecPacket{
			Facilities: []string{"string"},
		},
		Provider:            "string",
		RequiredEmailDomain: "string",
		Seed:                "string",
		Vsphare: &DatacenterSpecVsphare{
			Cluster:    "string",
			Datacenter: "string",
			DataStore:  "string",
			Endpoint:   "string",
			Templates: map[string]string{
				"additionalProp1": "string",
				"additionalProp2": "string",
				"additionalProp3": "string",
			},
		},
	},
}

func TestDatacenters_List(t *testing.T) {
	datacentersJSON := fmt.Sprintf(`[%s]`, datacenterJSON)
	testResourceList(t, datacentersJSON, "/api/v1/dc", []Datacenter{datacenter}, func() (interface{}, error) {
		l, _, err := client.Datacenters.List(ctx)
		return l, err
	})
}

func TestDatacenters_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/v1/dc/"+datacenter.Metadata.Name, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, datacenterJSON)
	})

	got, _, err := client.Datacenters.Get(ctx, datacenter.Metadata.Name)
	testErrNil(t, err)

	if !reflect.DeepEqual(&datacenter, got) {
		t.Fatalf("want: %+v, got: %+v", &datacenter, got)
	}
}
