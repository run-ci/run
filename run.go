package run

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/run-ci/run/pkg/run"
)

func RunTask(name string) {
	task, err := LoadTask(name)
	if err != nil {
		log.Fatalf("error loading task %v: %v", name, err)
	}

	log.Debugf("task %#v loaded", task)

	env, err := task.GetEnv()
	if err != nil {
		log.Fatalf("error loading task arguments: %v", err)
	}

	agent, err := run.NewAgent(nil)
	if err != nil {
		log.Fatalf("error creating run agent: %v", err)
	}

	log.Debugf("run agent initialized")

	err = agent.VerifyImagePresent(task.Image, true)
	if err != nil {
		log.Fatalf("error verifying image present: %v", err)
	}

	log.Debugf("image %v is present in the cache", task.Image)

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("error getting current working directory: %v", err)
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

	log.Debugf("running container with spec: %#v", spec)

	id, status, err := agent.RunContainer(spec)
	if err != nil {
		log.Fatalf("error running container with id %v: %v", id, err)
	}

	log.Debugf("task container exited with status %v", status)

	os.Exit(status)
}
