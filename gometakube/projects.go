package gometakube

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

const (
	projectsBasePath = "/api/v1/projects"
)

type Project struct {
	CreationTimestamp *time.Time        `json:"creationTimestamp"`
	DeletionTimestamp *time.Time        `json:"deletionTimestamp"`
	ID                string            `json:"id"`
	Labels            map[string]string `json:"labels"`
	Name              string            `json:"name"`
	Owners            []ProjectOwner    `json:"owners"`
	Status            string            `json:"status"`
}

type ProjectOwner struct {
	CreationTimestamp *time.Time     `json:"creationTimestamp"`
	DeletionTimestamp *time.Time     `json:"deletionTimestamp"`
	Email             string         `json:"email"`
	ID                string         `json:"id"`
	Name              string         `json:"name"`
	Projects          []OwnerProject `json:"projects"`
}

type OwnerProject struct {
	Group string `json:"group"`
	ID    string `json:"id"`
}

// ProjectsService handles communication with projects related endpoints.
type ProjectsService struct {
	client *Client
}

// List returns list of all projects.
func (svc *ProjectsService) List(ctx context.Context) ([]Project, error) {
	ret := make([]Project, 0)
	if resp, err := svc.client.serviceList(ctx, projectsBasePath, &ret); err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, unexpectedResponseError(resp)
	}
	return ret, nil
}

// ProjectCreateRequest payload to Create and Update a project.
type ProjectCreateRequest struct {
	Labels map[string]string `json:"labels"`
	Name   string            `json:"name"`
}

// Create creates a project.
func (svc *ProjectsService) Create(ctx context.Context, create *ProjectCreateRequest) (*Project, error) {
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
func (svc *ProjectsService) Update(ctx context.Context, id string, update *ProjectCreateRequest) (*Project, error) {
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
