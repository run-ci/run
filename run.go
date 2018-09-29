package main

import (
	"fmt"
	"os"

	docker "github.com/fsouza/go-dockerclient"
)

func runTask(name string) {
	task, err := LoadTask(name)
	if err != nil {
		printFatal("error loading task %v: %v", name, err)
	}

	printDebug("task %v loaded", task.Name)

	ccfg := &docker.Config{
		Image:        task.Image,
		Cmd:          []string{"echo", "foo"},
		AttachStderr: true,
		AttachStdout: true,
	}

	hcfg := &docker.HostConfig{}

	ncfg := &docker.NetworkingConfig{}

	client, err := docker.NewClient("unix:///var/run/docker.sock")
	if err != nil {
		printFatal("error opening docker socket: %v", err)
	}

	printDebug("docker client initialized")

	cnt, err := client.CreateContainer(docker.CreateContainerOptions{
		Config:           ccfg,
		HostConfig:       hcfg,
		NetworkingConfig: ncfg,
	})
	if err != nil {
		printFatal("error creating container for task: %v", err)
	}

	printDebug("container %v created", cnt.ID)

	go func() {
		printDebug("attaching container %v", cnt.ID)

		attachcfg := docker.AttachToContainerOptions{
			Container: cnt.ID,
			Stderr:    true,
			Stdout:    true,
			Stream:    true,
			Logs:      true,

			OutputStream: os.Stdout,
			ErrorStream:  os.Stderr,
		}

		err = client.AttachToContainer(attachcfg)
		if err != nil {
			msgs := fmt.Sprintf("error attaching to task container: %v", err)

			err := cleanupContainer(client, cnt.ID)
			if err != nil {
				msgs = fmt.Sprintf("%v\nerror cleaning up container %v: %v", msgs, cnt.ID, err)
			}

			printFatal(msgs)
		}
	}()

	printDebug("starting container %v", cnt.ID)

	err = client.StartContainer(cnt.ID, cnt.HostConfig)
	if err != nil {
		msgs := fmt.Sprintf("error starting task container: %v", err)

		err := cleanupContainer(client, cnt.ID)
		if err != nil {
			msgs = fmt.Sprintf("%v\nerror cleaning up container %v: %v", msgs, cnt.ID, err)
		}

		printFatal(msgs)
	}

	status, err := client.WaitContainer(cnt.ID)
	if err != nil {
		msgs := fmt.Sprintf("error running task container: %v", err)

		err := cleanupContainer(client, cnt.ID)
		if err != nil {
			msgs = fmt.Sprintf("%v\nerror cleaning up container %v: %v", msgs, cnt.ID, err)
		}

		printFatal(msgs)
	}

	fmt.Printf("task container exited with status %v\n", status)

	err = cleanupContainer(client, cnt.ID)
	if err != nil {
		printFatal("error cleaning up container %v: %v", cnt.ID, err)
	}
}

func cleanupContainer(client *docker.Client, id string) error {
	return client.RemoveContainer(docker.RemoveContainerOptions{
		ID:            id,
		RemoveVolumes: true,
		Force:         true,
	})
}
