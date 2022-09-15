# docker-gen

docker-gen is a simple tool that generates Dockerfiles from a base template and provided values.

## Build

To build, simply use the provided Makefile:

```bash
make build
```

This will put the binary in the `bin` directory.

## Usage

Full usage information is available by running `docker-gen --help`.

```bash
docker-gen generates validated Dockerfiles from go templated Dockerfiles.

Usage:
  docker-gen [flags]
  docker-gen [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  run         Renders Dockerfiles from templates with provided values
  version     Print the version of docker-gen

Flags:
  -h, --help   help for docker-gen
```

## Examples

To generate a simple Hello World Dockerfile from the base template using the provided values, run the following command:

```bash
bin/docker-gen run -d applications/hello/hello-values.yaml -t templates base
```

This will be default output the final Dockerfile next to the values, in this case `applications/hello/Dockerfile`.

## CI

To test the CI locally, install [act](https://github.com/nektos/act) and run the following command:

```bash
act -P ubuntu-latest=catthehacker/ubuntu:act-latest \
    -s DOCKERHUB_USERNAME=my-username \
    -s DOCKERHUB_TOKEN='my-secret-token'
```

Or simply comment out the `Login to DockerHub` step when testing locally.
