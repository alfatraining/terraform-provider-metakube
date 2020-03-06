package gometakube

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
				Projects: []OwnerProjects{
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

	mux.HandleFunc("/api/v1/projects", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, projectsJSON)
	})

	got, err := client.Projects.List(ctx)
	testErrNil(t, err)

	if want := []Project{project}; !reflect.DeepEqual(want, got) {
		t.Fatalf("want: %+v, got: %+v", want, got)
	}
}

func TestProjects_Create(t *testing.T) {
	setup()
	defer teardown()

	createRequest := &ProjectCreateAndUpdateRequest{
		Labels: map[string]string{
			"additionalProp1": "string",
			"additionalProp2": "string",
			"additionalProp3": "string",
		},
		Name: "myproject",
	}
	mux.HandleFunc("/api/v1/projects", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)
		v := &ProjectCreateAndUpdateRequest{}
		if err := json.NewDecoder(r.Body).Decode(v); err != nil {
			t.Fatalf("unexpected request parse error: %v", err)
		}
		if !reflect.DeepEqual(createRequest, v) {
			t.Fatalf("want: %v, got: %v", createRequest, v)
		}
		w.WriteHeader(http.StatusCreated)
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

func TestProject_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/v1/projects/"+project.ID, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
	})

	testErrNil(t, client.Projects.Delete(ctx, project.ID))
}

func TestProject_Update(t *testing.T) {
	setup()
	defer teardown()

	update := ProjectCreateAndUpdateRequest{
		Name: "name",
		Labels: map[string]string{
			"label1": "string",
		},
	}
	mux.HandleFunc("/api/v1/projects/"+project.ID, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPut)
		got := &ProjectCreateAndUpdateRequest{}
		if err := json.NewDecoder(r.Body).Decode(got); err != nil {
			t.Fatalf("unexpected request parse error: %v", err)
		}
		if want := &update; !reflect.DeepEqual(want, got) {
			t.Fatalf("want: %+v, got: %+v", want, got)
		}
		fmt.Fprint(w, projectJSON)
	})

	_, err := client.Projects.Update(ctx, project.ID, &update)
	testErrNil(t, err)
}
