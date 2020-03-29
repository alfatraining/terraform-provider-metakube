package gometakube

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

const (
	nodeDeploymentSpecJSON = `{
	"replicas": 3,
	"template": {
	  "cloud": {
		"openstack": {
		  "flavor": "m1.small",
		  "image": "Ubuntu Bionic 18.04 (2020-02-19)",
		  "tags": {
			"metakube-cluster": "j7f2svjll8",
			"system-cluster": "j7f2svjll8",
			"system-project": "5hrnkmpmp4"
		  },
		  "useFloatingIP": true
		}
	  },
	  "operatingSystem": {
		"ubuntu": {
		  "distUpgradeOnBoot": false
		}
	  },
	  "versions": {
		"kubelet": "1.17.2"
	  },
	  "labels": {
		"system/cluster": "j7f2svjll8",
		"system/project": "5hrnkmpmp4"
	  }
	},
	"paused": false
  }`
	nodeDeploymentJSON = `
  {
    "id": "metakube-worker-2xkvd",
    "name": "metakube-worker-2xkvd",
    "creationTimestamp": "2020-02-20T08:17:22Z",
    "spec": ` + nodeDeploymentSpecJSON + `,
    "status": {
      "observedGeneration": 1,
      "replicas": 3,
      "updatedReplicas": 3,
      "readyReplicas": 3,
      "availableReplicas": 3
    }
  }`
	prj = "theproj"
	dc  = "thedc"
	cls = "theclust"
)

var nodeDeployment = NodeDeployment{
	ID:                "metakube-worker-2xkvd",
	Name:              "metakube-worker-2xkvd",
	CreationTimestamp: testParseTime("2020-02-20T08:17:22Z"),
	Spec: NodeDeploymentSpec{
		Replicas: 3,
		Template: NodeDeploymentSpecTemplate{
			Cloud: NodeDeploymentSpecTemplateCloud{
				Openstack: NodeDeploymentSpecTemplateCloudOpenstack{
					Flavor: "m1.small",
					Image:  "Ubuntu Bionic 18.04 (2020-02-19)",
					Tags: map[string]string{
						"metakube-cluster": "j7f2svjll8",
						"system-cluster":   "j7f2svjll8",
						"system-project":   "5hrnkmpmp4",
					},
					UseFloatingIP: true,
				},
			},
			OperatingSystem: NodeDeploymentSpecTemplateOS{
				Ubuntu: &NodeDeploymentSpecTemplateOSOptions{
					DistUpgradeOnBoot: new(bool),
				},
			},
			Versions: NodeDeploymentSpecTemplateVersions{
				Kubelet: "1.17.2",
			},
			Labels: map[string]string{
				"system/cluster": "j7f2svjll8",
				"system/project": "5hrnkmpmp4",
			},
		},
		Paused: false,
	},
	Status: &NodeDeploymentStatus{
		ObservedGeneration: 1,
		Replicas:           3,
		UpdatedReplicas:    3,
		ReadyReplicas:      3,
		AvailableReplicas:  3,
	},
}

func TestNodeDeployments_List(t *testing.T) {
	nodeDeploymentsJSON := fmt.Sprintf("[%s]", nodeDeploymentJSON)
	path := fmt.Sprintf("/api/v1/projects/%s/dc/%s/clusters/%s/nodedeployments", prj, dc, cls)
	want := []NodeDeployment{nodeDeployment}
	testResourceList(t, nodeDeploymentsJSON, path, want, func() (interface{}, error) {
		return client.NodeDeployments.List(ctx, prj, dc, cls)
	})
}

func TestNodeDeployments_Patch(t *testing.T) {
	setup()
	defer teardown()

	url := fmt.Sprintf("/api/v1/projects/%s/dc/%s/clusters/%s/nodedeployments/%s", prj, dc, cls, nodeDeployment.ID)
	mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPatch)
		fmt.Fprint(w, nodeDeploymentJSON)
	})

	got, err := client.NodeDeployments.Patch(ctx, prj, dc, cls, nodeDeployment.ID, &NodeDeploymentsPatchRequest{nodeDeployment.Spec})
	testErrNil(t, err)
	if want := &nodeDeployment; !reflect.DeepEqual(want, got) {
		t.Fatalf("want: %+v, got: %+v", want, got)
	}
}

func TestNodeDeployments_Create(t *testing.T) {
	setup()
	defer teardown()

	path := fmt.Sprintf("/api/v1/projects/%s/dc/%s/clusters/%s/nodedeployments", prj, dc, cls)
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		v := new(NodeDeployment)
		if err := json.NewDecoder(r.Body).Decode(v); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(v, &nodeDeployment) {
			t.Fatalf("want receive node deployment: %v, got: %v", nodeDeployment, v)
		}
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, nodeDeploymentJSON)
	})

	got, err := client.NodeDeployments.Create(ctx, prj, dc, cls, &nodeDeployment)
	testErrNil(t, err)
	if !reflect.DeepEqual(got, &nodeDeployment) {
		t.Fatalf("want: %+v, got: %+v", nodeDeployment, got)
	}
}

func TestNodeDeployments_Get(t *testing.T) {
	setup()
	defer teardown()

	path := fmt.Sprintf("/api/v1/projects/%s/dc/%s/clusters/%s/nodedeployments/%s", prj, dc, cls, nodeDeployment.ID)
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, nodeDeploymentJSON)
	})

	got, err := client.NodeDeployments.Get(ctx, prj, dc, cls, nodeDeployment.ID)
	testErrNil(t, err)
	if !reflect.DeepEqual(got, &nodeDeployment) {
		t.Fatalf("want: %+v, got: %+v", nodeDeployment, got)
	}
}

func TestNodeDeployments_Delete(t *testing.T) {
	path := fmt.Sprintf("/api/v1/projects/%s/dc/%s/clusters/%s/nodedeployments/%s", prj, dc, cls, nodeDeployment.ID)
	testResourceDelete(t, path, func() error {
		return client.NodeDeployments.Delete(ctx, prj, dc, cls, nodeDeployment.ID)
	})
}

func TestNodeDeployments_Upgrade(t *testing.T) {
	setup()
	defer teardown()

	want := &UpgradeNodesRequest{
		Version: "1.17.5",
	}
	got := new(UpgradeNodesRequest)
	path := fmt.Sprintf("/api/v1/projects/%s/dc/%s/clusters/%s/nodes/upgrades", prj, dc, cls)
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		json.NewDecoder(r.Body).Decode(got)
	})

	err := client.NodeDeployments.Upgrade(ctx, prj, dc, cls, want)
	testErrNil(t, err)
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("wanted upgrade request: %+v, got: %+v", want, got)
	}
}
