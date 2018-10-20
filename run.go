package main

import (
	"os"

	"github.com/run-ci/run/pkg/run"
)

func runTask(name string) {
	task, err := LoadTask(name)
	if err != nil {
		printFatal("error loading task %v: %v", name, err)
	}

	printDebug("task %#v loaded", task)

	env, err := task.GetEnv()
	if err != nil {
		printFatal("error loading task arguments: %v", err)
	}

	agent, err := run.NewAgent(nil)
	if err != nil {
		printFatal("error creating run agent: %v", err)
	}

	printDebug("run agent initialized")

	err = agent.VerifyImagePresent(task.Image, true)
	if err != nil {
		printFatal("error verifying image present: %v", err)
	}

	printDebug("image %v is present in the cache", task.Image)

	pwd, err := os.Getwd()
	if err != nil {
		printFatal("error getting current working directory: %v", err)
	}

	spec := run.ContainerSpec{
		Imgref: task.Image,
		Cmd:    task.GetCmd(),
		Env:    env,
		Mount: run.Mount{
			Src:     pwd,
			Point:   task.Mount,
			Type:    "bind",
			Cleanup: true,
		},
	}

	printDebug("running container with spec: %#v", spec)

	id, status, err := agent.RunContainer(spec)
	if err != nil {
		printFatal("error running container with id %v: %v", id, err)
	}

	printDebug("task container exited with status %v", status)

	os.Exit(status)
}
