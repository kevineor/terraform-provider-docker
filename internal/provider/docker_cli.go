package provider

import (
	"fmt"
	"log"
	"strings"

	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/flags"
	"github.com/docker/docker/client"
)

// createAndInitDockerCli creates and initializes a Docker CLI instance backed by the
// provided client. originalHost is the host string from the provider configuration
// (e.g. "ssh://user@host"); it is used verbatim when SSH is configured so that the
// CLI endpoint reflects the SSH address rather than the internal HTTP helper URL that
// the Docker client reports via DaemonHost().
func createAndInitDockerCli(dockerClient *client.Client, originalHost string) (*command.DockerCli, error) {
	dockerCli, err := command.NewDockerCli()
	if err != nil {
		return nil, fmt.Errorf("failed to create Docker CLI: %w", err)
	}

	host := dockerClient.DaemonHost()
	if strings.HasPrefix(originalHost, "ssh://") {
		// For SSH connections the Docker client's DaemonHost() returns an internal
		// HTTP URL (e.g. "http://docker.example.com") that is only reachable through
		// the SSH dialer. Passing that URL to the CLI causes Compose to try a direct
		// HTTP connection and fail. Use the original SSH URL instead so the CLI
		// endpoint is correct and any subsequent connections use SSH.
		host = originalHost
	}

	log.Printf("[DEBUG] Docker CLI initialized with host: %s", host)
	options := flags.NewClientOptions()
	if host != "" {
		options.Hosts = []string{host}
	}

	if err := dockerCli.Initialize(options, command.WithAPIClient(dockerClient)); err != nil {
		return nil, fmt.Errorf("failed to initialize Docker CLI: %w", err)
	}
	return dockerCli, nil
}
