package gometakube

import (
	"context"
	"fmt"
	"net/http"
)

const (
	projectsBasePath = "/api/v1/projects"
)

// ProjectsService handles communication with projects related endpoints.
type ProjectsService struct {
	client *Client
}

// List returns list of all projects.
func (svc *ProjectsService) List(ctx context.Context) ([]Project, *http.Response, error) {
	ret := make([]Project, 0)
	resp, err := svc.client.resourceList(ctx, projectsBasePath, &ret)
	return ret, resp, err
}

// ProjectCreateAndUpdateRequest payload to Create and Update a project.
type ProjectCreateAndUpdateRequest struct {
	Labels map[string]string `json:"labels"`
	Name   string            `json:"name"`
}

// Create creates a project.
func (svc *ProjectsService) Create(ctx context.Context, create *ProjectCreateAndUpdateRequest) (*Project, *http.Response, error) {
	ret := new(Project)
	resp, err := svc.client.resourceCreate(ctx, projectsBasePath, create, ret)
	return ret, resp, err
}

// Get gets projects with given id.
func (svc *ProjectsService) Get(ctx context.Context, id string) (*Project, *http.Response, error) {
	path := projectResourcePath(id)
	ret := new(Project)
	resp, err := svc.client.resourceGet(ctx, path, ret)
	return ret, resp, err
}

// Update updates a project.
func (svc *ProjectsService) Update(ctx context.Context, id string, update *ProjectCreateAndUpdateRequest) (*Project, *http.Response, error) {
	req, err := svc.client.NewRequest(http.MethodPut, projectResourcePath(id), update)
	if err != nil {
		return nil, nil, err
	}
	ret := new(Project)
	resp, err := svc.client.Do(ctx, req, &ret)
	return ret, resp, err
}

// Delete deletes projects with given id.
func (svc *ProjectsService) Delete(ctx context.Context, id string) (*http.Response, error) {
	path := projectResourcePath(id)
	return svc.client.resourceDelete(ctx, path)
}

func projectResourcePath(id string) string {
	return fmt.Sprint(projectsBasePath, "/", id)
}
