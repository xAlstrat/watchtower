package mocks

import (
	"errors"
	"github.com/containrrr/watchtower/pkg/container"
	"time"

	t "github.com/containrrr/watchtower/pkg/types"
	cli "github.com/docker/docker/client"
)

// MockClientDependency is a mock that passes as a watchtower Client
type MockClientDependency struct {
	TestData      *TestDependencyData
	api           cli.CommonAPIClient
	pullImages    bool
	removeVolumes bool
}

// TestDependencyData is the data used to perform the test
type TestDependencyData struct {
	NameOfStaleContainer				string
	NameOfContainersShouldNotBeStopped	[]string
	Containers              			[]container.Container
	StoppedContainersNames				[]string
	RestartedContainersNames			[]string
}

// CreateMockDependencyClient creates a mock watchtower Client for usage in tests
func CreateMockDependencyClient(data *TestDependencyData, api cli.CommonAPIClient, pullImages bool, removeVolumes bool) MockClientDependency {
	return MockClientDependency{
		data,
		api,
		pullImages,
		removeVolumes,
	}
}

// ListContainers is a mock method returning the provided container testdata
func (client MockClientDependency) ListContainers(f t.Filter) ([]container.Container, error) {
	return client.TestData.Containers, nil
}

// StopContainer is a mock method
func (client MockClientDependency) StopContainer(c container.Container, d time.Duration) error {
	if stringInSlice(c.Name(), client.TestData.NameOfContainersShouldNotBeStopped) {
		return errors.New("tried to stop an instance that should not be stopped")
	}
	client.TestData.StoppedContainersNames = append(client.TestData.StoppedContainersNames, c.Name())
	return nil
}

// StartContainer is a mock method
func (client MockClientDependency) StartContainer(c container.Container) (string, error) {
	if stringInSlice(c.Name(), client.TestData.NameOfContainersShouldNotBeStopped) {
		return "", errors.New("tried to start an instance that should not be stopped")
	}
	client.TestData.RestartedContainersNames = append(client.TestData.RestartedContainersNames, c.Name())
	return "", nil
}

// RenameContainer is a mock method
func (client MockClientDependency) RenameContainer(c container.Container, s string) error {
	return nil
}

// RemoveImageByID increments the TriedToRemoveImageCount on being called
func (client MockClientDependency) RemoveImageByID(id string) error {
	return nil
}

// GetContainer is a mock method
func (client MockClientDependency) GetContainer(containerID string) (container.Container, error) {
	return container.Container{}, nil
}

// ExecuteCommand is a mock method
func (client MockClientDependency) ExecuteCommand(containerID string, command string, timeout int) error {
	return nil
}

// IsContainerStale is always true for the mock client
func (client MockClientDependency) IsContainerStale(c container.Container) (bool, error) {
	if c.Name() == client.TestData.NameOfStaleContainer {
		return true, nil
	}
	return false, nil
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
