package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

const (
	projectJSON = `
	{
		"creationTimestamp": "2020-02-12T13:51:53.926Z",
		"deletionTimestamp": "2020-02-12T13:51:53.926Z",
		"id": "string",
		"labels": {
		  "additionalProp1": "string",
		  "additionalProp2": "string",
		  "additionalProp3": "string"
		},
		"name": "string",
		"owners": [
		  {
			"creationTimestamp": "2020-02-12T13:51:53.926Z",
			"deletionTimestamp": "2020-02-12T13:51:53.926Z",
			"email": "string",
			"id": "string",
			"name": "string",
			"projects": [
			  {
				"group": "string",
				"id": "string"
			  }
			]
		  }
		],
		"status": "string"
	}
	`
)

var (
	project = Project{
		CreationTimestamp: testParseTime("2020-02-12T13:51:53.926Z"),
		DeletionTimestamp: testParseTime("2020-02-12T13:51:53.926Z"),
		ID:                "string",
		Labels: map[string]string{
			"additionalProp1": "string",
			"additionalProp2": "string",
			"additionalProp3": "string",
		},
		Name: "string",
		Owners: []ProjectOwner{
			{
				CreationTimestamp: testParseTime("2020-02-12T13:51:53.926Z"),
				DeletionTimestamp: testParseTime("2020-02-12T13:51:53.926Z"),
				Email:             "string",
				ID:                "string",
				Name:              "string",
				Projects: []OwnerProject{
					{
						Group: "string",
						ID:    "string",
					},
				},
			},
		},
		Status: "string",
	}
)

func TestProjects_List(t *testing.T) {
	setup()
	defer teardown()

	projectsJSON := fmt.Sprintf("[%s]", projectJSON)

	mux.HandleFunc(projectsListhPath, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, projectsJSON)
	})

	got, err := client.Projects.List(ctx)
	testErrNil(t, err)

	want := []Project{project}

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("want: %+v, got: %+v", want, got)
	}
}

func TestProjects_Create(t *testing.T) {
	setup()
	defer teardown()

	createRequest := &ProjectCreateRequest{
		Labels: map[string]string{
			"additionalProp1": "string",
			"additionalProp2": "string",
			"additionalProp3": "string",
		},
		Name: "myproject",
	}
	mux.HandleFunc(projectsCreatePath, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		v := &ProjectCreateRequest{}
		if err := json.NewDecoder(r.Body).Decode(v); err != nil {
			t.Fatalf("unexpected request parse error: %v", err)
		}
		if !reflect.DeepEqual(createRequest, v) {
			t.Fatalf("want: %v, got: %v", createRequest, v)
		}
		fmt.Fprint(w, projectJSON)
	})

	got, err := client.Projects.Create(ctx, createRequest)
	testErrNil(t, err)

	if want := &project; !reflect.DeepEqual(want, got) {
		t.Fatalf("want: %v, got: %v", want, got)
	}
}

func TestProject_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/v1/projects/"+project.ID, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, projectJSON)
	})

	got, err := client.Projects.Get(ctx, project.ID)
	testErrNil(t, err)

	if want := &project; !reflect.DeepEqual(want, got) {
		t.Fatalf("want: %v, got: %v", want, got)
	}
}
