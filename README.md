# run

"It works on my workstation" should be a valid excuse.

It shouldn't take endless commits messaged "testing CI", or endless broken
builds, to verify if a pipeline is working.

It shouldn't take days to get a workstation set up, or ramped up with a
project.

`run` solves these problems by using shrinkwrapped sets of dependencies
in a defined set of tasks that allow engineers to discover how to work
with a repo, locally or in CI.

## Quick Start

```
go install github.com/run-ci/run
# create a task YAML file in `./tasks/`
run -l # list all tasks
run -d TASK # describe TASK
run TASK # actually runs the task
```

## What is run?

`run` is a program that runs tasks in containers, defined by YAML files which
include metadata like a `description` and a `summary`, as well as arguments
defined by environment variables or paths to remote configuration.

`run` guarantees that tasks run the same way on different machines. Given that
the container image's version is pinned, there's no way that one person's setup,
or even a pipeline,  can be broken while another person's setup is working.

## Contributing

Before you begin, make sure you have the following installed:

- Docker (https://docs.docker.com/install/)
- docker-compose (https://docs.docker.com/compose/)
- `run`

1. Clone the repo

```
git clone https://github.com/run-ci/run
```

2. See what tasks are available!

```
run list
```

3. Optionally, run `docker-compose` for other services like Vault.

```
source env/local
docker-compose up
```
