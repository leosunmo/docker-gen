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

```text
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

For the `run` command, the following options are available:

```text
Usage:
  docker-gen run [flags]

Flags:
  -d, --data-file string      File containing data to use in template
  -f, --force                 Print rendered Dockerfiles even if they don't pass validation
  -h, --help                  help for run
  -o, --output-dir string     Directory to write rendered Dockerfiles to. Defaults to data directory if not piped to something else.
  -s, --stdout                Write rendered Dockerfiles to stdout
  -t, --template-dir string   Directory containing templates
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
    -s GITHUB_TOKEN='my-github-token' \
    -s DOCKERHUB_TOKEN='my-secret-token'
```

The Dockerhub token is only required if you want to push the image to Dockerhub, you could simply comment out the `Login to DockerHub` step when testing locally.
The Github token is required in the Docker Tags step as it uses the Github API to get the latest metadata it seems. Can't really get around it in a nice way.

## TODO

- [X] Add built-in Dockerfile validation.

- [ ] Convert to Github Action. If this was to be used by multiple projects, writing a Typescript Github Action would be the way to go.
- [ ] Make it a library + `cmd` CLI. Having it all in the `cmd` dir isn't very neat or flexible.
- [ ] Add more flexible options in the template renderer. Right now the values are hardcoded,
      so we'd need to make it more flexible.
- [ ] CI tidy up. A lot of hardcoded values.
