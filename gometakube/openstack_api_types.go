package gometakube

import "time"


type Image struct {
	ID       string
	Created  *time.Time
	MinDisk  uint
	MinRAM   uint
	Name     string
	Progress uint
	Status   string
	Updated  *time.Time
	Metadata ImageMetadata
}

type ImageMetadata struct {
	CIJobID            string `json:"ci_job_id"`
	CIPipelineID       string `json:"ci_pipeline_id"`
	CPUArch            string `json:"cpu_arch"`
	DefaultSSHUsername string `json:"default_ssh_username"`
	Distribution       string `json:"distribution"`
	OSDistro           string `json:"os_distro"`
	OSType             string `json:"os_type"`
	OSVersion          string `json:"os_version"`
	SourceSHA56sum     string `json:"source_sha256sum"`
	SourceURL          string `json:"source_url"`
}

type Tenant struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
