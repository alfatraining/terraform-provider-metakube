package client

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
	"time"
)

func TestProjects_List(t *testing.T) {
	client := New()
	ctx := context.TODO()
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	url, _ := url.Parse(server.URL)
	client.BaseURL = url

	projectsJSON := `[
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
	  ]`

	mux.HandleFunc("/api/v1/projects", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, projectsJSON)
	})

	got, err := client.Projects.List(ctx)
	testErrNil(t, err)

	want := []Project{
		{
			CreationTimestamp: testParseTime(t, "2020-02-12T13:51:53.926Z"),
			DeletionTimestamp: testParseTime(t, "2020-02-12T13:51:53.926Z"),
			ID:                "string",
			Labels: map[string]string{
				"additionalProp1": "string",
				"additionalProp2": "string",
				"additionalProp3": "string",
			},
			Name: "string",
			Owners: []ProjectOwner{
				{
					CreationTimestamp: testParseTime(t, "2020-02-12T13:51:53.926Z"),
					DeletionTimestamp: testParseTime(t, "2020-02-12T13:51:53.926Z"),
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
		},
	}

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("want: %+v, got: %+v", want, got)
	}
}

func testParseTime(t *testing.T, s string) *time.Time {
	t.Helper()
	ret, err := time.Parse(time.RFC3339, s)
	if err != nil {
		t.Fatalf("failed to parse time string `%s`: %v", s, err)
	}
	return &ret
}
