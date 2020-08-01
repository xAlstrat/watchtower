package actions_test

import (
	"github.com/containrrr/watchtower/internal/actions"
	"github.com/containrrr/watchtower/pkg/container"
	"github.com/containrrr/watchtower/pkg/container/mocks"
	"github.com/containrrr/watchtower/pkg/types"
	cli "github.com/docker/docker/client"
	"time"

	. "github.com/containrrr/watchtower/internal/actions/mocks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func GetLinkedContainers(containers []container.Container) (linked_names []string) {
	linked_names = []string{}
	for _, container := range containers {
		if container.Linked {
			linked_names = append(linked_names, container.Name())
		}
	}
	return linked_names
}

var _ = Describe("the update action with dependencies", func() {
	var dockerClient cli.CommonAPIClient
	var client MockClientDependency

	BeforeEach(func() {
		server := mocks.NewMockAPIServer()
		dockerClient, _ = cli.NewClientWithOpts(
			cli.WithHost(server.URL),
			cli.WithHTTPClient(server.Client()))
	})

	When("watchtower has been instructed to update", func() {
		BeforeEach(func() {
			pullImages := false
			removeVolumes := false
			client = CreateMockDependencyClient(
				&TestDependencyData{
					NameOfStaleContainer: "test-container-01",
					NameOfContainersShouldNotBeStopped: []string{"test-container-03"},
					Containers: []container.Container{
						CreateMockContainerWithLabels(
							"test-container-01",
							"test-container-01",
							"fake-image-01:latest",
							time.Now().AddDate(0, 0, -1),
							map[string]string{}),
						CreateMockContainerWithLabels(
							"test-container-02",
							"test-container-02",
							"fake-image-02:latest",
							time.Now(),
							map[string]string{
								"com.centurylinklabs.watchtower.depends-on": "test-container-01",
							}),
						CreateMockContainerWithLabels(
							"test-container-03",
							"test-container-03",
							"fake-image-03:latest",
							time.Now(),
							map[string]string{
								"com.centurylinklabs.watchtower.depends-on": "",
							}),
					},
				},
				dockerClient,
				pullImages,
				removeVolumes,
			)
		})

		When("there are linked containers", func() {
			It("should stop the linked containers", func() {
				err := actions.Update(client, types.UpdateParams{})
				Expect(err).NotTo(HaveOccurred())
				Expect(client.TestData.StoppedContainersNames).To(SatisfyAll(ContainElement("test-container-01"), ContainElement("test-container-02"), HaveLen(2)))
			})
			It("should restart the linked containers", func() {
				err := actions.Update(client, types.UpdateParams{})
				Expect(err).NotTo(HaveOccurred())
				Expect(client.TestData.RestartedContainersNames).To(SatisfyAll(ContainElement("test-container-01"), ContainElement("test-container-02"), HaveLen(2)))
			})
			It("should stop the parent container first", func() {
				// test-container-02 depends-on its children test-container-01
				err := actions.Update(client, types.UpdateParams{})
				Expect(err).NotTo(HaveOccurred())
				Expect(client.TestData.StoppedContainersNames[0]).To(Equal("test-container-02"))
			})
			It("should restart the child container first", func() {
				err := actions.Update(client, types.UpdateParams{})
				Expect(err).NotTo(HaveOccurred())
				Expect(client.TestData.RestartedContainersNames[0]).To(Equal("test-container-01"))
			})
			It("shouldn't stop not linked containers", func() {
				err := actions.Update(client, types.UpdateParams{})
				Expect(err).NotTo(HaveOccurred())
				Expect(client.TestData.StoppedContainersNames).NotTo(ContainElement("test-container-03"))
			})
			It("shouldn't restart not linked containers", func() {
				err := actions.Update(client, types.UpdateParams{})
				Expect(err).NotTo(HaveOccurred())
				Expect(client.TestData.RestartedContainersNames).NotTo(ContainElement("test-container-03"))
			})
		})
	})

	When("watchtower has been instructed to update and has a different container order", func() {
		// The result should be the same as before no matter the order of the containers
		BeforeEach(func() {
			pullImages := false
			removeVolumes := false
			client = CreateMockDependencyClient(
				&TestDependencyData{
					NameOfStaleContainer: "test-container-01",
					NameOfContainersShouldNotBeStopped: []string{"test-container-03"},
					Containers: []container.Container{
						CreateMockContainerWithLabels(
							"test-container-02",
							"test-container-02",
							"fake-image-02:latest",
							time.Now(),
							map[string]string{
								"com.centurylinklabs.watchtower.depends-on": "test-container-01",
							}),
						CreateMockContainerWithLabels(
							"test-container-03",
							"test-container-03",
							"fake-image-03:latest",
							time.Now(),
							map[string]string{
								"com.centurylinklabs.watchtower.depends-on": "",
							}),
						CreateMockContainerWithLabels(
							"test-container-01",
							"test-container-01",
							"fake-image-01:latest",
							time.Now().AddDate(0, 0, -1),
							map[string]string{}),
					},
				},
				dockerClient,
				pullImages,
				removeVolumes,
			)
		})

		When("there are linked containers", func() {
			It("should stop the linked containers", func() {
				err := actions.Update(client, types.UpdateParams{})
				Expect(err).NotTo(HaveOccurred())
				Expect(client.TestData.StoppedContainersNames).To(SatisfyAll(ContainElement("test-container-01"), ContainElement("test-container-02"), HaveLen(2)))
			})
			It("should restart the linked containers", func() {
				err := actions.Update(client, types.UpdateParams{})
				Expect(err).NotTo(HaveOccurred())
				Expect(client.TestData.RestartedContainersNames).To(SatisfyAll(ContainElement("test-container-01"), ContainElement("test-container-02"), HaveLen(2)))
			})
			It("should stop the parent container first", func() {
				// test-container-02 depends-on its children test-container-01
				err := actions.Update(client, types.UpdateParams{})
				Expect(err).NotTo(HaveOccurred())
				Expect(client.TestData.StoppedContainersNames[0]).To(Equal("test-container-02"))
			})
			It("should restart the child container first", func() {
				err := actions.Update(client, types.UpdateParams{})
				Expect(err).NotTo(HaveOccurred())
				Expect(client.TestData.RestartedContainersNames[0]).To(Equal("test-container-01"))
			})
			It("shouldn't stop not linked containers", func() {
				err := actions.Update(client, types.UpdateParams{})
				Expect(err).NotTo(HaveOccurred())
				Expect(client.TestData.StoppedContainersNames).NotTo(ContainElement("test-container-03"))
			})
			It("shouldn't restart not linked containers", func() {
				err := actions.Update(client, types.UpdateParams{})
				Expect(err).NotTo(HaveOccurred())
				Expect(client.TestData.RestartedContainersNames).NotTo(ContainElement("test-container-03"))
			})
		})
	})

	When("watchtower has been instructed to update with nested dependency", func() {
		BeforeEach(func() {
			pullImages := false
			removeVolumes := false
			client = CreateMockDependencyClient(
				&TestDependencyData{
					NameOfStaleContainer: "test-container-01",
					NameOfContainersShouldNotBeStopped: []string{"test-container-05"},
					Containers: []container.Container{
						CreateMockContainerWithLabels(
							"test-container-01",
							"test-container-01",
							"fake-image-01:latest",
							time.Now().AddDate(0, 0, -1),
							map[string]string{}),
						CreateMockContainerWithLabels(
							"test-container-02",
							"test-container-02",
							"fake-image-02:latest",
							time.Now(),
							map[string]string{
								"com.centurylinklabs.watchtower.depends-on": "test-container-01",
							}),
						CreateMockContainerWithLabels(
							"test-container-03",
							"test-container-03",
							"fake-image-03:latest",
							time.Now(),
							map[string]string{
								"com.centurylinklabs.watchtower.depends-on": "test-container-02",
							}),
						CreateMockContainerWithLabels(
							"test-container-04",
							"test-container-04",
							"fake-image-04:latest",
							time.Now(),
							map[string]string{
								"com.centurylinklabs.watchtower.depends-on": "test-container-03",
							}),
						CreateMockContainerWithLabels(
							"test-container-05",
							"test-container-05",
							"fake-image-05:latest",
							time.Now(),
							map[string]string{
								"com.centurylinklabs.watchtower.depends-on": "",
							}),
					},
				},
				dockerClient,
				pullImages,
				removeVolumes,
			)
		})

		When("there are linked containers", func() {
			It("should stop the linked containers", func() {
				err := actions.Update(client, types.UpdateParams{})
				Expect(err).NotTo(HaveOccurred())
				Expect(client.TestData.StoppedContainersNames).To(SatisfyAll(
					ContainElement("test-container-04"),
					ContainElement("test-container-03"),
					ContainElement("test-container-02"),
					ContainElement("test-container-01"),
					HaveLen(4)))
			})
			It("should restart the linked containers", func() {
				err := actions.Update(client, types.UpdateParams{})
				Expect(err).NotTo(HaveOccurred())
				Expect(client.TestData.RestartedContainersNames).To(SatisfyAll(
					ContainElement("test-container-01"),
					ContainElement("test-container-02"),
					ContainElement("test-container-03"),
					ContainElement("test-container-04"),
					HaveLen(4)))
			})
			It("should stop containers in the correct order", func() {
				// test-container-02 depends-on its children test-container-01
				err := actions.Update(client, types.UpdateParams{})
				Expect(err).NotTo(HaveOccurred())
				Expect(client.TestData.StoppedContainersNames).To(Equal([]string{
					"test-container-04",
					"test-container-03",
					"test-container-02",
					"test-container-01",
				}))
			})
			It("should restart containers in the correct order", func() {
				err := actions.Update(client, types.UpdateParams{})
				Expect(err).NotTo(HaveOccurred())
				Expect(client.TestData.RestartedContainersNames).To(Equal([]string{
					"test-container-01",
					"test-container-02",
					"test-container-03",
					"test-container-04",
				}))
			})
		})
	})
})
