# run

The universal task runner.

## Quick Start

```
go install github.com/run-ci/run
# create a task YAML file in `./tasks/`
run -l # list all tasks
run -d TASK # describe TASK
run TASK # actually runs the task
```

## Why run?

Task runners help make life easier. When you're new to a project, looking at
the different tasks available for working with the repo can help you get up
to speed much faster. At the same time, experienced engineers don't have to
be wasting time doing things that can just be automated.

Different toolchains have different opinions on what a task runner should be.
Ruby has `rake`, Go has `mage` and `grift`, and some projects just use
different assortments of shell scripts or Python scripts. This creates a lot
of unnecessary overhead in a couple of ways.

The first way is by not being consistent from toolchain to toolchain. `rake`
thinks you should use it one way, while `mage` and `grift` also have their
own opinions. This makes it just a bit harder to change between different
projects that use different toolchains.

The second way is by having to manage the task runners' dependencies, plus
all the dependencies being used in their tasks. For `rake` you have to run
`bundle install`, and even then gem versions aren't necessarily guaranteed
to stay the same which can cause two different engineers to have different
workspaces, one of which can be broken while the other one works perfectly.

Containers can solve problem number two, but what about number one? And what
about the value the task runners can provide?

`run` aims to solve the dependency problem by using containers, while also
providing the same value that task runners can provide, and providing a
consistent workflow regardless of the toolchain a project is using.

## What is run?

`run` is a program that runs tasks in containers, defined by YAML files which
include metadata like a `description` and a `summary`, as well as arguments
defined by environment variables or paths to remote configuration.

`run` guarantees that tasks run the same way on different machines. Given that
the container image's version is pinned, there's no way that one person's
setup can be broken while another person's setup is working.

## Contributing

Before you begin, make sure you have the following installed:

- Docker (https://docs.docker.com/install/)
- docker-compose (https://docs.docker.com/compose/)
- `run`

1. Clone the repo

```
git clone https://github.com/run-ci/run
```

2. See what's available!

```
run -l
```

3. Optionally, run `docker-compose` for other services like Vault.

```
source env/local
docker-compose up
```
