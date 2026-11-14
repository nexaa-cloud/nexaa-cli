package cmd

import (
	"testing"

	"github.com/nexaa-cloud/nexaa-cli/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildContainerModifyInput_PreservesExistingHealthCheck(t *testing.T) {
	oldContainer := api.ContainerResult{
		Name:      "backend",
		Image:     "ghcr.io/org/backend:v1",
		Resources: api.ContainerResourcesCpu1000Ram2000,
		HealthCheck: &api.ContainerResultHealthCheck{
			Port: 2026,
			Path: "/health",
		},
	}

	opts := modifyContainerOptions{
		Namespace: "my-ns",
		Name:      "backend",
		Image:     "ghcr.io/org/backend:v2",
	}

	input := buildContainerModifyInput(oldContainer, opts)

	assert.Equal(t, "backend", input.Name)
	assert.Equal(t, "my-ns", input.Namespace)
	require.NotNil(t, input.Image)
	assert.Equal(t, "ghcr.io/org/backend:v2", *input.Image)

	// Health check must be preserved, NOT nil (which would serialize as null and disable it)
	require.NotNil(t, input.HealthCheck, "existing health check must not be dropped when modifying container")
	assert.Equal(t, 2026, input.HealthCheck.Port)
	assert.Equal(t, "/health", input.HealthCheck.Path)
}

func TestBuildContainerModifyInput_UpdateHealthCheck(t *testing.T) {
	oldContainer := api.ContainerResult{
		Name:      "backend",
		Image:     "ghcr.io/org/backend:v1",
		Resources: api.ContainerResourcesCpu1000Ram2000,
		HealthCheck: &api.ContainerResultHealthCheck{
			Port: 2026,
			Path: "/health",
		},
	}

	opts := modifyContainerOptions{
		Namespace:       "my-ns",
		Name:            "backend",
		HealthCheckPort: 8080,
		HealthCheckPath: "/ready",
	}

	input := buildContainerModifyInput(oldContainer, opts)

	require.NotNil(t, input.HealthCheck)
	assert.Equal(t, 8080, input.HealthCheck.Port)
	assert.Equal(t, "/ready", input.HealthCheck.Path)
}

func TestBuildContainerModifyInput_RemoveHealthCheck(t *testing.T) {
	oldContainer := api.ContainerResult{
		Name:      "backend",
		Image:     "ghcr.io/org/backend:v1",
		Resources: api.ContainerResourcesCpu1000Ram2000,
		HealthCheck: &api.ContainerResultHealthCheck{
			Port: 2026,
			Path: "/health",
		},
	}

	opts := modifyContainerOptions{
		Namespace:         "my-ns",
		Name:              "backend",
		RemoveHealthCheck: true,
	}

	input := buildContainerModifyInput(oldContainer, opts)

	// When explicitly removing, HealthCheck should be nil (serializes to null)
	assert.Nil(t, input.HealthCheck)
}

func TestBuildContainerModifyInput_PreservesCommandAndEntrypoint(t *testing.T) {
	oldContainer := api.ContainerResult{
		Name:       "backend",
		Command:    []string{"start", "--prod"},
		Entrypoint: []string{"/app/entry.sh"},
		Resources:  api.ContainerResourcesCpu1000Ram2000,
		PrivateRegistry: &api.ContainerResultPrivateRegistry{
			Name: "my-registry",
		},
	}

	opts := modifyContainerOptions{
		Namespace: "my-ns",
		Name:      "backend",
	}

	input := buildContainerModifyInput(oldContainer, opts)

	assert.Equal(t, []string{"start", "--prod"}, input.Command)
	assert.Equal(t, []string{"/app/entry.sh"}, input.Entrypoint)
	require.NotNil(t, input.Registry)
	assert.Equal(t, "my-registry", *input.Registry)
	require.NotNil(t, input.Resources)
	assert.Equal(t, api.ContainerResourcesCpu1000Ram2000, *input.Resources)
}

func TestBuildContainerModifyInput_RemoveCommandAndEntrypoint(t *testing.T) {
	oldContainer := api.ContainerResult{
		Name:       "backend",
		Command:    []string{"start"},
		Entrypoint: []string{"/app/entry.sh"},
		Resources:  api.ContainerResourcesCpu1000Ram2000,
	}

	opts := modifyContainerOptions{
		Namespace:        "my-ns",
		Name:             "backend",
		RemoveCommand:    true,
		RemoveEntrypoint: true,
	}

	input := buildContainerModifyInput(oldContainer, opts)

	assert.Empty(t, input.Command)
	assert.Empty(t, input.Entrypoint)
}
