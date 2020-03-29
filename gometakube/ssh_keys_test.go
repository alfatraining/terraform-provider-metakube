package gometakube

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

const sshkeyJSON = `{
  "creationTimestamp": "2020-03-27T19:58:59.707Z",
  "deletionTimestamp": "2020-03-27T19:58:59.707Z",
  "id": "string",
  "name": "string",
  "spec": {
	"fingerprint": null,
	"publicKey": "string"
  }
}`

var sshkey = SSHKey{
	CreationTimestamp: testParseTime("2020-03-27T19:58:59.707Z"),
	DeletionTimestamp: testParseTime("2020-03-27T19:58:59.707Z"),
	ID:                "string",
	Name:              "string",
	Spec: SSHKeySpec{
		Fingerprint: nil,
		PublicKey:   "string",
	},
}

func TestSHHKeys_List(t *testing.T) {
	listJSON := "[" + sshkeyJSON + "]"
	want := []SSHKey{sshkey}
	prj := "project"
	path := fmt.Sprintf("/api/v1/projects/%s/sshkeys", prj)
	testResourceList(t, listJSON, path, want, func() (interface{}, error) {
		return client.SSHKeys.List(ctx, prj)
	})
}

func TestSSHKeys_ListInCluster(t *testing.T) {
	listJSON := "[" + sshkeyJSON + "]"
	want := []SSHKey{sshkey}
	prj := "project"
	dc := "datacenter"
	cls := "cluster"
	path := fmt.Sprintf("/api/v1/projects/%s/dc/%s/clusters/%s/sshkeys", prj, dc, cls)
	testResourceList(t, listJSON, path, want, func() (interface{}, error) {
		return client.SSHKeys.ListAssigned(ctx, prj, dc, cls)
	})
}

func TestSHHKeys_Create(t *testing.T) {
	setup()
	defer teardown()

	prj := "project"
	path := fmt.Sprintf("/api/v1/projects/%s/sshkeys", prj)
	createReq := &sshkey
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		v := new(SSHKey)
		if err := json.NewDecoder(r.Body).Decode(v); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(v, createReq) {
			t.Fatalf("want: %+v, got: %+v", createReq, v)
		}
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, sshkeyJSON)
	})

	got, err := client.SSHKeys.Create(ctx, prj, createReq)
	testErrNil(t, err)
	if !reflect.DeepEqual(got, &sshkey) {
		t.Fatalf("want: %+v, got: %+v", sshkey, got)
	}
}

func TestSSHKeys_Delete(t *testing.T) {
	prj := "prj"
	id := "id"
	path := fmt.Sprintf("/api/v1/projects/%s/sshkeys/%s", prj, id)
	testResourceDelete(t, path, func() error {
		return client.SSHKeys.Delete(ctx, prj, id)
	})
}

func TestSSHKeys_AssingToCluster(t *testing.T) {
	setup()
	defer teardown()

	prj := "theproject"
	dc := "datacenter"
	cls := "thecluster"
	id := "mykeyid"
	path := fmt.Sprintf("/api/v1/projects/%s/dc/%s/clusters/%s/sshkeys/%s", prj, dc, cls, id)
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, sshkeyJSON)
	})

	got, err := client.SSHKeys.AssignToCluster(ctx, prj, dc, cls, id)
	testErrNil(t, err)

	if !reflect.DeepEqual(&sshkey, got) {
		t.Fatalf("want: %+v, got: %+v", sshkey, got)
	}
}

func TestSSHKeys_RemoveFromCluster(t *testing.T) {
	prj := "theproject"
	dc := "datacenter"
	cls := "thecluster"
	id := "mykeyid"
	path := fmt.Sprintf("/api/v1/projects/%s/dc/%s/clusters/%s/sshkeys/%s", prj, dc, cls, id)
	testResourceDelete(t, path, func() error {
		return client.SSHKeys.RemoveFromCluster(ctx, prj, dc, cls, id)
	})
}
