package gometakube

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type SSHKey struct {
	CreationTimestamp *time.Time `json:"creationTimestamp,omitempty"`
	DeletionTimestamp *time.Time `json:"deletionTimestamp,omitempty"`
	ID                string     `json:"id"`
	Name              string     `json:"name"`
	Spec              SSHKeySpec `json:"spec"`
}

type SSHKeySpec struct {
	Fingerprint *string `json:"fingerprint"`
	PublicKey   string  `json:"publicKey"`
}

// SSHKeysService handle communication with sshkeys resource endpoints.
type SSHKeysService struct {
	client *Client
}

func projectSSHKeysPath(prj string) string {
	return fmt.Sprintf("/api/v1/projects/%s/sshkeys", prj)
}

func deleteSSHKeyPath(prj, id string) string {
	return fmt.Sprintf("/api/v1/projects/%s/sshkeys/%s", prj, id)
}

func clusterSSHKeysPath(prj, dc, cls string) string {
	return fmt.Sprintf("/api/v1/projects/%s/dc/%s/clusters/%s/sshkeys", prj, dc, cls)
}

func clusterSSHKeyPath(prj, dc, cls, id string) string {
	return fmt.Sprintf("/api/v1/projects/%s/dc/%s/clusters/%s/sshkeys/%s", prj, dc, cls, id)
}

// List returns list of sshkeys in a project.
func (svc *SSHKeysService) List(ctx context.Context, prj string) ([]SSHKey, *http.Response, error) {
	ret := make([]SSHKey, 0)
	resp, err := svc.client.resourceList(ctx, projectSSHKeysPath(prj), &ret)
	return ret, resp, err
}

// ListAssigned returns list of sshkeys assigned to a cluster.
func (svc *SSHKeysService) ListAssigned(ctx context.Context, prj, dc, cls string) ([]SSHKey, *http.Response, error) {
	ret := make([]SSHKey, 0)
	resp, err := svc.client.resourceList(ctx, clusterSSHKeysPath(prj, dc, cls), &ret)
	return ret, resp, err
}

// Create adds sshkey to a project.
func (svc *SSHKeysService) Create(ctx context.Context, prj string, params *SSHKey) (*SSHKey, *http.Response, error) {
	ret := new(SSHKey)
	resp, err := svc.client.resourceCreate(ctx, projectSSHKeysPath(prj), params, ret)
	return ret, resp, err
}

// Delete deletes a sshkey from a project.
func (svc *SSHKeysService) Delete(ctx context.Context, prj, id string) (*http.Response, error) {
	path := deleteSSHKeyPath(prj, id)
	return svc.client.resourceDelete(ctx, path)
}

// AssignToCluster assign ssh key to a cluster.
func (svc SSHKeysService) AssignToCluster(ctx context.Context, prj, dc, cls, id string) (*SSHKey, *http.Response, error) {
	ret := new(SSHKey)
	path := clusterSSHKeyPath(prj, dc, cls, id)
	req, err := svc.client.NewRequest(http.MethodPut, path, nil)
	if err != nil {
		return nil, nil, err
	}
	resp, err := svc.client.Do(ctx, req, ret)
	return ret, resp, err
}

// RemoveFromCluster unassigns sshkey from a cluster.
func (svc *SSHKeysService) RemoveFromCluster(ctx context.Context, prj, dc, cls, id string) (*http.Response, error) {
	path := clusterSSHKeyPath(prj, dc, cls, id)
	return svc.client.resourceDelete(ctx, path)
}
