package client

import (
	"context"
	"net/http"
	"time"
)

const (
	projectsListhPath = "/api/v1/projects"
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

// ProjectsService handles communication with projects related methods.
type ProjectsService struct {
	client *Client
}

// List returns list of all projects.
func (svc *ProjectsService) List(ctx context.Context) ([]Project, error) {
	req, err := svc.client.NewRequest(http.MethodGet, projectsListhPath)
	if err != nil {
		return nil, err
	}
	ret := make([]Project, 0)
	if err := svc.client.Do(ctx, req, &ret); err != nil {
		return nil, err
	}
	return ret, nil
}
