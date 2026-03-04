// Package collector wraps github.com/aceextensions/collector/docker
// for use inside the devops-autopilot application.
package collector

import (
	"github.com/aceextensions/collector/docker"
)

// DockerStats re-exports the upstream type.
type DockerStats = docker.DockerStats

// ContainerInfo re-exports the upstream type.
type ContainerInfo = docker.ContainerInfo

// CollectDocker collects Docker container stats using the aceextensions/collector package.
func CollectDocker() (*DockerStats, error) {
	return docker.CollectDocker()
}
