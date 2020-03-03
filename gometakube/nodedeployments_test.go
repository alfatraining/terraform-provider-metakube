package gometakube

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

const nodeDeploymentJSON = `
  {
    "id": "metakube-worker-2xkvd",
    "name": "metakube-worker-2xkvd",
    "creationTimestamp": "2020-02-20T08:17:22Z",
    "spec": {
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
    },
    "status": {
      "observedGeneration": 1,
      "replicas": 3,
      "updatedReplicas": 3,
      "readyReplicas": 3,
      "availableReplicas": 3
    }
  }`

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
	setup()
	defer teardown()

	nodeDeploymentsJSON := fmt.Sprintf("[%s]", nodeDeploymentJSON)
	prj := "5hrnkmpmp4"
	dc := "bki1"
	cls := "j7f2svjll8"
	url := fmt.Sprintf("/api/v1/projects/%s/dc/%s/clusters/%s/nodedeployments", prj, dc, cls)
	mux.HandleFunc(url, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, nodeDeploymentsJSON)
	})

	got, err := client.NodeDeployments.List(ctx, prj, dc, cls)
	testErrNil(t, err)

	if want := []NodeDeployment{nodeDeployment}; !reflect.DeepEqual(want, got) {
		t.Fatalf("want: %+v, got: %+v", want, got)
	}
}
