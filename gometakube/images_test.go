package gometakube

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

const imageJSON = `{
		"ID": "imageid",
		"Created": "2020-03-03T04:06:52Z",
		"MinDisk": 0,
		"MinRAM": 0,
		"Name": "Ubuntu Bionic 18.04 (2020-03-03)",
		"Progress": 100,
		"Status": "ACTIVE",
		"Updated": "2020-03-03T04:09:07Z",
		"Metadata": {
		  "ci_job_id": "3106028",
		  "ci_pipeline_id": "424463",
		  "cpu_arch": "x86_64",
		  "default_ssh_username": "ubuntu",
		  "distribution": "ubuntu-bionic",
		  "os_distro": "ubuntu",
		  "os_type": "linux",
		  "os_version": "18.04",
		  "source_sha256sum": "bf21a56ba61864122f9893d33ec93db1e8d4dab3db366306115927f81fd2fae7",
		  "source_url": "https://cloud-images.ubuntu.com/bionic/current/bionic-server-cloudimg-amd64.img"
		}
	  }`

var image = Image{
	ID:       "imageid",
	Created:  testParseTime("2020-03-03T04:06:52Z"),
	MinDisk:  0,
	MinRAM:   0,
	Name:     "Ubuntu Bionic 18.04 (2020-03-03)",
	Progress: 100,
	Status:   "ACTIVE",
	Updated:  testParseTime("2020-03-03T04:09:07Z"),
	Metadata: ImageMetadata{
		CIJobID:            "3106028",
		CIPipelineID:       "424463",
		CPUArch:            "x86_64",
		DefaultSSHUsername: "ubuntu",
		Distribution:       "ubuntu-bionic",
		OSDistro:           "ubuntu",
		OSType:             "linux",
		OSVersion:          "18.04",
		SourceSHA56sum:     "bf21a56ba61864122f9893d33ec93db1e8d4dab3db366306115927f81fd2fae7",
		SourceURL:          "https://cloud-images.ubuntu.com/bionic/current/bionic-server-cloudimg-amd64.img",
	},
}

func TestImages_List(t *testing.T) {
	setup()
	defer teardown()

	dcName := "dc"
	username := "theuser"
	password := "pwd"
	domain := "Default"
	mux.HandleFunc("/api/v1/providers/openstack/images", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		if want, got := dcName, r.Header.Get("DatacenterName"); want != got {
			t.Fatalf("want DatacenterName: %v, got: %v", want, got)
		}
		if want, got := username, r.Header.Get("Username"); want != got {
			t.Fatalf("want Username: %v, got: %v", want, got)
		}
		if want, got := password, r.Header.Get("Password"); want != got {
			t.Fatalf("want Password: %v, got: %v", want, got)
		}
		if want, got := domain, r.Header.Get("Domain"); want != got {
			t.Fatalf("want Domain: %v, got: %v", want, got)
		}
		fmt.Fprintf(w, "[%s]", imageJSON)
	})

	got, err := client.Images.List(ctx, dcName, domain, username, password)
	testErrNil(t, err)
	if want := []Image{image}; !reflect.DeepEqual(want, got) {
		t.Fatalf("want: %+v, got: %+v", want, got)
	}
}
