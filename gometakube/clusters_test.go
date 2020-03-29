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
		  "auditLogging": {
			  "enabled": false
		  }
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
	clustersJSON := fmt.Sprintf("[%s]", clusterJSON)
	prj := "foo"
	path := fmt.Sprintf("/api/v1/projects/%s/clusters", prj)

	testResourceList(t, clustersJSON, path, []Cluster{cluster}, func() (interface{}, error) {
		return client.Clusters.List(ctx, prj)
	})
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
	path := fmt.Sprintf("/api/v1/projects/%s/dc/%s/clusters", prj, dc)
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
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
		t.Fatalf("want: %+v, got: %+v", want, got)
	}
}

func TestClusters_Delete(t *testing.T) {
	prj := "the-proj"
	dc := "bk11"
	cls := "the-cluster"
	path := fmt.Sprintf("/api/v1/projects/%s/dc/%s/clusters/%s", prj, dc, cls)
	testResourceDelete(t, path, func() error {
		return client.Clusters.Delete(ctx, prj, dc, cls)
	})
}

func TestClusters_Get(t *testing.T) {
	setup()
	defer teardown()

	prj := "the-prj"
	dc := "thedc"
	cls := "thecluster"
	path := fmt.Sprintf("/api/v1/projects/%s/dc/%s/clusters/%s", prj, dc, cls)
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, clusterJSON)
	})

	got, err := client.Clusters.Get(ctx, prj, dc, cls)
	testErrNil(t, err)
	if want := &cluster; !reflect.DeepEqual(want, got) {
		t.Fatalf("want: %+v, got: %+v", want, got)
	}
}

func TestClusters_Patch(t *testing.T) {
	setup()
	defer teardown()

	prj := "the-prj"
	dc := "thedc"
	cls := "thecluster"
	path := fmt.Sprintf("/api/v1/projects/%s/dc/%s/clusters/%s", prj, dc, cls)
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		fmt.Fprint(w, clusterJSON)
	})
	patch := &PatchClusterRequest{
		Name: "edited",
		Labels: map[string]string{
			"newkey": "newvalue",
		},
		Spec: &PatchClusterRequestSpec{
			AuditLogging: &ClusterSpecAuditLogging{
				Enabled: false,
			},
		},
	}
	got, err := client.Clusters.Patch(ctx, prj, dc, cls, patch)
	testErrNil(t, err)
	if want := &cluster; !reflect.DeepEqual(want, got) {
		t.Fatalf("want: %v, got: %v", want, got)
	}
}

const clusterHealthJSON = `
  {
	"apiserver": 0,
	"cloudProviderInfrastructure": 0,
	"controller": 0,
	"etcd": 0,
	"machineController": 0,
	"scheduler": 0,
	"userClusterControllerManager": 0
  }
`

var clusterHealth = ClusterHealth{
	APIServer:                    0,
	CloudProviderInfrastructure:  0,
	Controller:                   0,
	Etcd:                         0,
	MachineController:            0,
	Scheduler:                    0,
	UserClusterControllerManager: 0,
}

func TestClusters_Health(t *testing.T) {
	setup()
	defer teardown()

	prj := "the-proj"
	dc := "thedc"
	cls := "thecluster"
	path := fmt.Sprintf("/api/v1/projects/%s/dc/%s/clusters/%s/health", prj, dc, cls)
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, clusterHealthJSON)
	})

	got, err := client.Clusters.Health(ctx, prj, dc, cls)
	testErrNil(t, err)

	if want := &clusterHealth; !reflect.DeepEqual(want, got) {
		t.Fatalf("want: %+v, err: %+v", want, got)
	}
}

func TestClusters_Upgrades(t *testing.T) {
	setup()
	defer teardown()
	upgradesWant := []ClusterUpgrade{
		{Version: "1.16.8"},
		{Version: "1.17.4", Defailt: true},
	}
	path := "/api/v1/upgrades/cluster"
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `[{ "version": "1.16.8" }, { "version": "1.17.4", "default": true }]`)
	})

	got, err := client.Clusters.Upgrades(ctx)
	testErrNil(t, err)
	if !reflect.DeepEqual(upgradesWant, got) {
		t.Fatalf("want upgrades: %+v, got: %+v", upgradesWant, got)
	}
}

func TestClusters_Upgrade(t *testing.T) {
	setup()
	defer teardown()
	prj := "theproject"
	dc := "somedc"
	id := "id"
	path := fmt.Sprintf("/api/v1/projects/%s/dc/%s/clusters/%s/upgrades", prj, dc, id)
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `[{ "version": "1.16.8" }]`)
	})

	got, err := client.Clusters.ClusterUpgrades(ctx, prj, dc, id)
	testErrNil(t, err)
	want := []ClusterUpgrade{{Version: "1.16.8"}}
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("want cluster upgrades: %v, got: %v", want, got)
	}
}
