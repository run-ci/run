package run

import (
	"errors"
	"os"
	"strings"

	"github.com/fsouza/go-dockerclient"
)

// Agent is an abstraction on top of a Docker client that enforces certain
// conventions when running containers. These conventions are geared towards
// running tasks with containers.
type Agent struct {
	client *docker.Client
}

// NewAgent returns an initialized Agent. By default, this agent is
// initialized with a Docker client pointing to the Docker socket at
// `/var/run/docker.sock`, but a custom Docker client can be passed
// in for use instead.
func NewAgent(client *docker.Client) (*Agent, error) {
	if client == nil {
		var err error
		client, err = docker.NewClient("unix:///var/run/docker.sock")
		if err != nil {
			return &Agent{}, err
		}
	}

	return &Agent{
		client: client,
	}, nil
}

// VerifyImagePresent checks to see if the image with the specified tag is
// present and pulls it if necessary. Set `verbose` to `true` to see the
// Docker engine's output in stdout.
func (ag *Agent) VerifyImagePresent(imgref string, verbose bool) error {
	imgsegs := strings.Split(imgref, ":")

	imgs, err := ag.client.ListImages(docker.ListImagesOptions{
		All:    true,
		Filter: imgsegs[0],
	})
	if err != nil {
		return err
	}

	if len(imgs) > 0 {
		return nil
	}

	pullopts := docker.PullImageOptions{
		Repository: imgsegs[0],
	}

	if len(imgsegs) > 1 {
		pullopts.Tag = imgsegs[1]
	}

	if verbose {
		pullopts.OutputStream = os.Stdout
	}

	return ag.client.PullImage(pullopts, docker.AuthConfiguration{})
}

// RunContainer runs a container with the given spec. The
// container's working directory will be the same as the
// `Mount.Point` on the container.
//
// The first return argument is the container ID. If an error
// occurred before the container was able to be created, an
// empty string is returned instead.
//
// The second return argument is the task return status.
//
// If the mount type is not one of "bind" or "volume", it
// returns -1 and an error. If something ever goes wrong, it
// returns -1 and the underlying error. If all goes well, it
// returns the status code returned by the container process.
func (ag *Agent) RunContainer(spec ContainerSpec) (string, int, error) {
	switch spec.Mount.Type {
	case "bind", "volume":
	default:
		return "", -1, errors.New("unknown mount type")
	}

	ccfg := &docker.Config{
		Image:        spec.Imgref,
		Cmd:          spec.Cmd,
		AttachStderr: true,
		AttachStdout: true,
		Env:          spec.GetEnvArray(),
		Volumes: map[string]struct{}{
			spec.Mount.Point: struct{}{},
		},
		WorkingDir: spec.Mount.Point,
	}

	hcfg := &docker.HostConfig{
		Mounts: []docker.HostMount{
			docker.HostMount{
				Target: spec.Mount.Point,
				Source: spec.Mount.Src,
				Type:   spec.Mount.Type,
			},

			docker.HostMount{
				Target: "/var/run/docker.sock",
				Source: "/var/run/docker.sock",
				Type:   "bind",
			},
		},
	}

	ncfg := &docker.NetworkingConfig{}

	cnt, err := ag.client.CreateContainer(docker.CreateContainerOptions{
		Config:           ccfg,
		HostConfig:       hcfg,
		NetworkingConfig: ncfg,
	})
	if err != nil {
		return "", -1, err
	}

	defer ag.CleanupContainer(cnt.ID, spec.Mount.Cleanup)

	go func() {
		attachcfg := docker.AttachToContainerOptions{
			Container: cnt.ID,
			Stderr:    true,
			Stdout:    true,
			Stream:    true,
			Logs:      true,

			OutputStream: os.Stdout,
			ErrorStream:  os.Stderr,
		}

		err = ag.client.AttachToContainer(attachcfg)
		if err != nil {
			// TODO: decide what to do here. Probably just quit the whole thing.
		}
	}()

	err = ag.client.StartContainer(cnt.ID, cnt.HostConfig)
	if err != nil {
		return "", -1, err
	}

	status, err := ag.client.WaitContainer(cnt.ID)
	return cnt.ID, status, err
}

// CleanupContainer forcibly removes a container and, if specified by `cleanVols`,
// removes its volumes.
func (ag *Agent) CleanupContainer(id string, cleanVols bool) error {
	return ag.client.RemoveContainer(docker.RemoveContainerOptions{
		ID:            id,
		RemoveVolumes: cleanVols,
		Force:         true,
	})
}
