package gometakube

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

const (
	clusterJSON = `
	{
		"id": "idstring",
		"name": "thecluster",
		"creationTimestamp": "2020-02-20T08:15:47Z",
		"type": "kubernetes",
		"spec": {
		  "cloud": {
			"dc": "syseleven-datacenter",
			"openstack": {
			  "floatingIpPool": "ext-net",
			  "subnetCIDR": "192.168.1.0/24"
			}
		  },
		  "clusterNetwork": {
			"services": {
			  "cidrBlocks": [
				"10.240.16.0/20"
			  ]
			},
			"pods": {
			  "cidrBlocks": [
				"172.25.0.0/16"
			  ]
			},
			"dnsDomain": "cluster.local",
			"proxyMode": "ipvs"
		  },
		  "version": "1.17.2",
		  "oidc": {},
		  "sys11auth": {},
		  "auditLogging": {}
		},
		"status": {
		  "version": "1.17.2",
		  "url": "https://url"
		}
	  }	  
	`
)

var (
	cluster = Cluster{
		ID:                "idstring",
		Name:              "thecluster",
		CreationTimestamp: testParseTime("2020-02-20T08:15:47Z"),
		Type:              "kubernetes",
		Spec: &ClusterSpec{
			Cloud: &ClusterSpecCloud{
				DataCenter: "syseleven-datacenter",
				OpenStack: &ClusterSpecCloudOpenstack{
					FloatingIPPool: "ext-net",
					SubnetCIDR:     "192.168.1.0/24",
				},
			},
			ClusterNetwork: &ClusterSpecClusterNetwork{
				Services: &ClusterSpecClusterNetworkServices{
					CIDRBlocks: []string{"10.240.16.0/20"},
				},
				Pods: &ClusterSpecClusterNetworkPods{
					CIDRBlocks: []string{"172.25.0.0/16"},
				},
				DNSDomain: "cluster.local",
				ProxyMode: "ipvs",
			},
			Version: "1.17.2",
		},
		Status: &ClusterStatus{
			URL:     "https://url",
			Version: "1.17.2",
		},
	}
)

func TestClusters_List(t *testing.T) {
	setup()
	defer teardown()

	clustersJSON := fmt.Sprintf("[%s]", clusterJSON)
	prj := "foo"
	mux.HandleFunc("/api/v1/projects/foo/clusters", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, clustersJSON)
	})

	got, err := client.Clusters.List(ctx, prj)
	testErrNil(t, err)

	if want := []Cluster{cluster}; !reflect.DeepEqual(want, got) {
		t.Fatalf("want: %+v, got: %+v", want, got)
	}
}

func TestClusters_Create(t *testing.T) {
	setup()
	defer teardown()

	createRequest := &CreateClusterRequest{
		Cluster:        Cluster{ID: "id-cluster"},
		NodeDeployment: NodeDeployment{ID: "id-nodeDeployment"},
	}

	prj := "the-proj"
	dc := "bki1"
	url := fmt.Sprintf("/api/v1/projects/%s/dc/%s/clusters", prj, dc)
	mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		v := &CreateClusterRequest{}
		if err := json.NewDecoder(r.Body).Decode(v); err != nil {
			t.Fatalf("want: %v, got: %v", *createRequest, *v)
		}
		if !reflect.DeepEqual(createRequest, v) {
			t.Fatalf("want: %v, got: %v", *createRequest, *v)
		}
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, `{"id": "id-cluster"}`)
	})

	got, err := client.Clusters.Create(ctx, prj, dc, createRequest)
	testErrNil(t, err)

	if want := createRequest.Cluster; want.ID != got.ID {
		t.Fatalf("want: %v, got: %v", want, got)
	}
}
