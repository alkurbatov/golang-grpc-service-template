# golang-grpc-service-template

Template for fast creation of new `Golang gRPC` service.
Contains many ready-to-use solutions for routine tasks like metrics,logging etc.

## Configuration

The configuration is powered by `viper`.
For additional options and tweaks please refer to [the official documentation](https://github.com/spf13/viper).

### General

| Env variable                            | Type | Description                                                                             | Default value |
|:----------------------------------------|:-----|:----------------------------------------------------------------------------------------|:--------------|
| `GRPC_LISTEN_ADDRESS`                   | str  | `gRPC API` listening address and port                                                   | 0.0.0.0:50051 |
| `HTTP_LISTEN_ADDRESS`                   | str  | `HTTP API` listening address and port to get `Prometheus` metrics.                      | 0.0.0.0:8888  |
| `GRACEFUL_SHUTDOWN_PRESTOP_TIMEOUT_SEC` | int  | Time for service to sleep before starting to reject new requests                        | 10            |
| `GRACEFUL_SHUTDOWN_GRPC_TIMEOUT_SEC`    | int  | Time for `gRPC` server to process all the remaining requests before being shutting down | 280           |

### Logging

| Env variable   | Type | Description                                                         | Default value |
|:---------------|:-----|:--------------------------------------------------------------------|:--------------|
| `LOG_LEVEL`    | str  | Logging level of the application (`debug`, `info`, `warn`, `error`) | `info`        |
| `JSON_LOGGING` | bool | Enable logging in `JSON` format                                     | true          |

### Other

| Env variable          | Type | Description                                                                                                     | Default value |
|:----------------------|:-----|:----------------------------------------------------------------------------------------------------------------|:--------------|
| `PROF_LISTEN_ADDRESS` | str  | `pprof` profiler listening address and port. Should differ from `HTTP_LISTEN_ADDRESS` due to security concerns. |               |

## Project structure

```text
├── cmd                           # main entrypoint, contains main.go files of all components
│
├── deployments                   # project configuration for local deployment
│   ├── docker-compose.debug.yaml
│   ├── docker-compose.dev.yaml
│   ├── docker-compose.test.yaml
│   └── docker-compose.yaml
│
├── docs                          # project documentation
│
├── internal                      # internally used code that shouldn't be imported by library users
│   ├── app
│   │   └── app.go                # main starting point responsible for service lifecycle
│   ├── config
│   ├── controller
│   │   ├── grpc                  # gRPC API handlers
│   │   └── grpc                  # HTTP API handlers (e.g. metrics)
│   └── infra                     # infrastructure layer, contains abstractions of gRPC server etc
│
├── pkg                           # code that could be reused by other projects
│   └── echopb                    # generated code of public service gRPC API
│
├── proto                         # gRPC API protocol files
├── scripts                       # handy build and maintenance scripts
│
├── test
│   └── smoke                     # smoke tests
│
├── .dockerignore
├── .gitignore
│
├── .golangci.yaml                # golangci-lint configuration
├── .hadolint.yaml
├── .pre-commit-config.yaml
│
├── CONTRIBUTING.md               # decribes programing techniques recommended to use in the project
│
├── Dockerfile                    # main Dockerfile to build image for production
├── Dockerfile.debug              # Dockerfile to build image for remote debugging
├── Dockerfile.test               # Dockerfile to build smoke test image
│
├── go.mod                        # package dependencies
├── go.sum                        # dependencies lock file
│
├── Makefile
└── README.md
```

### Additional directories

The directories not used in this template but could be needed in some cases.

#### `test/testdata`

To store artifacts used during unit and smoke tests, e.g. audio samples, `JSON` files, tokens etc.

### Naming convention

The project follows the [`Golang` naming convention][naming-convention].
For detailed explanation of particular folder purpose please refer to
[Standard Go Project Layout](https://github.com/golang-standards/project-layout).

## Creating new service from this project

You may prefer to use the `gonew` utility to create a copy of this project
instead of classic approaches like clone, fork or manual copy-paste.
The following command will create a fresh copy of the library in `new-service`
folder with new module name in `go.mod` and without initialization of version
control system:

``` bash
go run golang.org/x/tools/cmd/gonew@latest github.com/alkurbatov/golang-grpc-service-template@latest github.com/alkurbatov/new-service
```

> :exclamation: Don't forget to rename remains of the old service (e.g. folders and binaries named `templatesrv`).

## Development and testing

### Setup dev environment

1. Install `protoc` (to generate `protobuf` code) as described [here](https://grpc.io/docs/protoc-installation/).
1. Install `gofumpt` (for better code formatting) as described [here](https://github.com/mvdan/gofumpt#installation).
1. Install `gci` (for better code formatting) as described [here](https://github.com/daixiang0/gci?tab=readme-ov-file#installation).
1. Install `golangci-lint` (to lint `Golang` code) according to [this guide](https://golangci-lint.run/welcome/install/).
1. Install `shellcheck` (to lint `bash` scripts) according to [this guide](https://github.com/koalaman/shellcheck#installing).
1. Install `pre-commit` (to run linters before commit) according to [this guide](https://pre-commit.com/#install).

### Prepare to work with the project

1. Clone project repository:

   ```bash
   git clone git@github.com:alkurbatov/golang-grpc-service-template.git
   ```

1. Install additional development tools, `protobuf` plugins and download
   dependencies:

   ```bash
   make install-tools
   ```

1. Install `pre-commit` hooks by running:

   ```bash
   pre-commit install
   ```

### Run the project in Docker

To run the project in Docker using docker-compose execute the following command:

```bash
make run
```

To stop the running project do:

```bash
make stop
```

### Run the project from sources

```bash
go run ./cmd/templatesrv/main.go
```

## Workflow and commands

### Sync project dependencies

To sync project dependencies and lock them:

```bash
go mod tidy
```

### Generate protobuf bindings for `gRPC API`

```bash
make proto
```

### View documentation

> :bulb: When writing `Golang` documentation comments we follow [official style guide](https://tip.golang.org/doc/comment).

Project documentation is available via [`godoc`](https://pkg.go.dev/golang.org/x/tools/cmd/godoc).

```bash
make docs
```

### Linting

Lint the sources with all linters.

```bash
make lint
```

#### False positives

To ignore particular `golangci-lint` errors or false positives inline please
use the inline comments e.g.

```go
//nolint:interfacebloat // no plans to split it right now
//nolint:cyclop,gocyclo // no need in simplification
```

For more powerful settings please follow the recommendations in [documentation of `golangci-lint`](https://golangci-lint.run/usage/false-positives/).

### Run unit tests

To run unit tests execute the following command:

```bash
make unit-tests
```

Unit tests are shuffled on each run. To reproduce previous run in exactly same
order, extract randomization seed from the beginning of test logs and specify it
on next run. E.g. assuming that the seed value was 1725189750482165000:

```bash
make unit-tests SEED=1725189750482165000
```

To update snapshots used in unit tests run the following command:

```bash
make update-snapshots
```

### Smoke tests

Smoke tests contain core operations that have to work. Service that does not
pass the smoke tests should not be passed to QA.

To run smoke tests (requires `docker compose`) do:

```bash
make smoke-tests
```

### Remote debugging

1. Install `delve` debugger according to
   [this instruction](https://github.com/go-delve/delve/tree/master/Documentation/installation).

1. Run the project in docker with enabled remote debugging interface:

   ```bash
   make debug
   ```

1. Attach the debugger:

   ```bash
   dlv connect :2345
   ```

1. While in `delve` remap path to sources to be able to view them during
   debugging, e.g. (assuming that current directory is root of the project):

   ```bash
   config substitute-path /app ./
   ```

### Other commands

To get full list of available commands run:

```bash
make help
```

[naming-convention]: https://www.mohitkhare.com/blog/go-naming-conventions/
