package gometakube

import "time"

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
	CreationTimestamp *time.Time      `json:"creationTimestamp"`
	DeletionTimestamp *time.Time      `json:"deletionTimestamp"`
	Email             string          `json:"email"`
	ID                string          `json:"id"`
	Name              string          `json:"name"`
	Projects          []OwnerProjects `json:"projects"`
}

type OwnerProjects struct {
	Group string `json:"group"`
	ID    string `json:"id"`
}
