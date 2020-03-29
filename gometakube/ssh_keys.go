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
func (svc *SSHKeysService) List(ctx context.Context, prj string) ([]SSHKey, error) {
	ret := make([]SSHKey, 0)
	if resp, err := svc.client.resourceList(ctx, projectSSHKeysPath(prj), &ret); err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, unexpectedResponseError(resp)
	}
	return ret, nil
}

// ListAssigned returns list of sshkeys assigned to a cluster.
func (svc *SSHKeysService) ListAssigned(ctx context.Context, prj, dc, cls string) ([]SSHKey, error) {
	ret := make([]SSHKey, 0)
	if resp, err := svc.client.resourceList(ctx, clusterSSHKeysPath(prj, dc, cls), &ret); err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		return nil, unexpectedResponseError(resp)
	}
	return ret, nil
}

// Create adds sshkey to a project.
func (svc *SSHKeysService) Create(ctx context.Context, prj string, params *SSHKey) (*SSHKey, error) {
	ret := new(SSHKey)
	return ret, svc.client.resourceCreate(ctx, projectSSHKeysPath(prj), params, ret)
}

// Delete deletes a sshkey from a project.
func (svc *SSHKeysService) Delete(ctx context.Context, prj, id string) error {
	path := deleteSSHKeyPath(prj, id)
	return svc.client.resourceDelete(ctx, path)
}

// AssignToCluster assign ssh key to a cluster.
func (svc SSHKeysService) AssignToCluster(ctx context.Context, prj, dc, cls, id string) (*SSHKey, error) {
	ret := new(SSHKey)
	path := clusterSSHKeyPath(prj, dc, cls, id)
	req, err := svc.client.NewRequest(http.MethodPut, path, nil)
	if err != nil {
		return nil, err
	}
	if resp, err := svc.client.Do(ctx, req, ret); err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusCreated {
		return nil, unexpectedResponseError(resp)
	} else {
		return ret, nil
	}
}

// RemoveFromCluster unassigns sshkey from a cluster
func (svc *SSHKeysService) RemoveFromCluster(ctx context.Context, prj, dc, cls, id string) error {
	path := clusterSSHKeyPath(prj, dc, cls, id)
	return svc.client.resourceDelete(ctx, path)
}
