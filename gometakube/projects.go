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
func (svc *ProjectsService) List(ctx context.Context) ([]Project, error) {
	ret := make([]Project, 0)
	if resp, err := svc.client.resourceList(ctx, projectsBasePath, &ret); err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, unexpectedResponseError(resp)
	}
	return ret, nil
}

// ProjectCreateAndUpdateRequest payload to Create and Update a project.
type ProjectCreateAndUpdateRequest struct {
	Labels map[string]string `json:"labels"`
	Name   string            `json:"name"`
}

// Create creates a project.
func (svc *ProjectsService) Create(ctx context.Context, create *ProjectCreateAndUpdateRequest) (*Project, error) {
	ret := new(Project)
	if resp, err := svc.client.resourceCreate(ctx, projectsBasePath, create, ret); err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusCreated {
		return nil, unexpectedResponseError(resp)
	}
	return ret, nil
}

// Get gets projects with given id.
func (svc *ProjectsService) Get(ctx context.Context, id string) (*Project, error) {
	url := projectResourcePath(id)
	ret := new(Project)
	if resp, err := svc.client.resourceGet(ctx, url, ret); err != nil {
		return nil, err
	} else if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	} else if resp.StatusCode != http.StatusOK {
		return nil, unexpectedResponseError(resp)
	}
	return ret, nil
}

// Update updates a project.
func (svc *ProjectsService) Update(ctx context.Context, id string, update *ProjectCreateAndUpdateRequest) (*Project, error) {
	req, err := svc.client.NewRequest(http.MethodPut, projectResourcePath(id), update)
	if err != nil {
		return nil, err
	}
	ret := new(Project)
	if resp, err := svc.client.Do(ctx, req, &ret); err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, unexpectedResponseError(resp)
	}
	return ret, nil
}

// Delete deletes projects with given id.
func (svc *ProjectsService) Delete(ctx context.Context, id string) error {
	url := projectResourcePath(id)
	if resp, err := svc.client.resourceDelete(ctx, url); err != nil {
		return err
	} else if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotFound {
		return unexpectedResponseError(resp)
	}
	return nil
}

func projectResourcePath(id string) string {
	return fmt.Sprint(projectsBasePath, "/", id)
}
