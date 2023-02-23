# File Storage Service

---

[![CI Workflow](https://github.com/moemoe89/file-storage/actions/workflows/ci.yml/badge.svg)](https://github.com/moemoe89/file-storage/actions/workflows/ci.yml) <!-- start-coverage --><img src="https://img.shields.io/badge/coverage-0.0%25-red"><!-- end-coverage -->

File Storage Service handles upload, list and delete related files data into storage.

## Table of Contents

<!-- @import "[TOC]" {cmd="toc" depthFrom=1 depthTo=6 orderedList=false} -->

<!-- code_chunk_output -->

- [Table of Contents](#table-of-contents)
- [Project Summary](#project-summary)
- [Installation](#installation)
    - [1. Set Up Golang Development Environment](#1-set-up-golang-development-environment)
    - [2. Install Development Utility Tools](#2-install-development-utility-tools)
- [Development workflow and guidelines](#development-workflow-and-guidelines)
    - [1. API](#1-api)
    - [2. Instrumentation](#2-instrumentation)
    - [3. Unit Test](#3-unit-test)
    - [4. Linter](#4-linter)
    - [5. Run the service](#5-run-the-service)
    - [6. Test the service](#6-test-the-service)
- [Project Structure](#project-structure)
- [GitHub Actions CI](#github-actions-ci)
- [Documentation](#documentation)
  - [Visualize Code Diagram](#visualize-code-diagram)

<!-- /code_chunk_output -->

## Project Summary

| Item                      | Description                                                                                                           |
|---------------------------|-----------------------------------------------------------------------------------------------------------------------|
| Golang Version            | [1.19](https://golang.org/doc/go1.19)                                                                                 |
| moq                       | [mockgen](https://github.com/golang/mock)                                                                             |
| Linter                    | [GolangCI-Lint](https://github.com/golangci/golangci-lint)                                                            |
| Testing                   | [testing](https://golang.org/pkg/testing/) and [testify/assert](https://godoc.org/github.com/stretchr/testify/assert) |
| API                       | [gRPC](https://grpc.io/docs/tutorials/basic/go/) and [gRPC-Gateway](https://github.com/grpc-ecosystem/grpc-gateway)   |
| Application Architecture  | [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)                    |
| Directory Structure       | [Standard Go Project Layout](https://github.com/golang-standards/project-layout)                                      |
| CI (Lint, Test, Generate) | [GitHubActions](https://github.com/features/actions)                                                                  |
| Visualize Code Diagram    | [go-callviz](https://github.com/ofabry/go-callvis)                                                                    |
| Sequence Diagram          | [Mermaid](https://mermaid.js.org)                                                                                     |
| Protobuf Operations       | [buf](https://buf.build)                                                                                              |
| Instrumentation           | [OpenTelemetry](https://opentelemetry.io) and [Jaeger](https://www.jaegertracing.io)                                  |
| Logger                    | [zap](https://github.com/uber-go/zap)                                                                                 |


## Installation

### 1. Set Up Golang Development Environment

See the following page to download and install Golang.

https://go.dev/doc/install

### 2. Install Development Utility Tools

You can install all tools for development and deployment for this service by running:

```sh
$ go mod download
```

```sh
$ make install
```

---

## Development workflow and guidelines

### 1. API

This project using gRPC and Protocol Buffers, thus all needed data like Service definition, RPC's list, Entities will store in [api/proto](api/proto) directory.

If you unfamiliar with Protocol Buffer, please visit this link for the detail:

* https://protobuf.dev

For generating the Proto files, make sure to have these libs installed on your system, please refer to this link:

* https://buf.build/
* https://grpc.io/docs/protoc-installation
* https://grpc.io/docs/languages/go/quickstart/

The validation for this API using `protoc-gen-validate`, for the detail please refer to this lib:

* https://github.com/bufbuild/protoc-gen-validate

This service also implementing gRPC-Gateway with this library:

* https://github.com/grpc-ecosystem/grpc-gateway

For generating the gRPC-Gateway and OpenAPI files, there's required additional package such as:

* github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
* github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2

Then, generating the Protobuf files can be done my this command:

```sh
$ make protoc
```

#### NOTE:

If you have any difficulties installing all dependencies needed for generating the proto files,
you can easily build the docker image and use that instead.

Here are 2 commands for building and generating:

```shell
make build-protoc
make docker-protoc
```

### 2. Instrumentation

This service implements [https://opentelemetry.io/](https://opentelemetry.io/) to enable instrumentation in order to measure the performance.
The data exported to Jaeger and can be seen in the Jaeger UI [http://localhost:16686](http://localhost:16686)

For running the Jaeger exporter, easily run with docker-compose command:

```sh
$ docker-compose -f ./development/docker-compose.yml up jaeger
```

### 3. Unit Test

You can simply execute the following command to run all test cases in this service:

```sh
$ make test
```

### 4. Linter

For running the linter make sure these libraries already installed in your system:

* https://github.com/golangci/golangci-lint
* https://github.com/yoheimuta/protolint

Then checks the Go and Proto code style using lint can be done with this command:

```sh
$ make lint
```

### 5. Mock

This service using Mock in some places like in the repository, usecase, pkg, etc.
To automatically updating the mock if the interface changed, easily run with `go generate` command:

```sh
$ make mock
```

### 5. Run the service

For running the service, you need the database running and set up some env variables:

```
# app config
export APP_ENV=dev
export SERVER_PORT=8080

# tracing config
export OTEL_AGENT=http://localhost:14268/api/traces
```

Or you can just execute the sh file:

```sh
$ ./scripts/run.sh
```

### 6. Test the service

The example how to call the gRPC service written in Golang can be seen on this [example-client](scripts/example-client) file.

> NOTE: To test this service need the migration to be done. After that you can choose the User ID's from 1 to 5.
 
If you want to test by GUI client, you can use either BloomRPC (although already no longer active) or Postman.
For the detail please visit these links:
* https://github.com/bloomrpc/bloomrpc
* https://www.postman.com

Basically you just need to import the [api/proto/service.proto](api/proto/service.proto) file if you want to test via BloomRPC / Postman.

> NOTE: There will be a possibility issue when importing the proto file to BloomRPC or Postman.
> It is caused by some path issue, the usage of `gRPC Gateway` and `protoc-gen-validate` library.
> To solve this issue, there's need a modification for the proto file.

#### BloomRPC

BloomRPC will have these issues when trying to import the proto file:

```
Error while importing protos
illegal name ';' (/path/file-storage/api/proto/service.proto, line 14)

Error while importing protos
no such type: e.Transaction
```

need to remove gRPC Gateway related annotations:

```protobuf
option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger) = {
  ...
};
```

#### Postman

There's some issue when importing to Postman. Basically we need to do the same things like BloomRPC (remove gRPC Gateway related annotations) and disable the protoc annotations import.

```protobuf
import "protoc-gen-openapiv2/options/annotations.proto";
```

To this:

```protobuf
// import "protoc-gen-openapiv2/options/annotations.proto";
```

Also don't forget to set the import path e.g. `{YOUR-DIR}/file-storage/api/proto`

#### gRPC-Gateway

This service has HTTP server built on gRPC-Gateway, if you prefer to test using HTTP instead HTTP2 protocol,
you can copy the Swagger file here [api/openapiv2/proto/service.swagger.json](api/openapiv2/proto/service.swagger.json) and then copy paste to this URL https://editor.swagger.io/

By default, HTTP server running on gRPC port + 1, if the gRPC port is 8080, then HTTP server will run on 8081.

# NOTE

> If you have any difficulties to run the service, easily just run all dependencies by docker-compose for the example:
> 
> `docker-compose -f ./development/docker-compose.yml up`
>
> Then you will have all services running like `jaeger` and run `file-storage` server.

## Project Structure

This project follow https://github.com/golang-standards/project-layout

However, for have a clear direction when working in this project, here are some small guide about each directory:

* [api](api): contains Protobuf files, generated protobuf, swagger, etc.
* [build](build): Docker file for the service, migration, etc.
* [cmd](cmd): main Go file for running the service, producer, consumer, etc.
* [development](development): file to support development like docker-compose.
* [docs](docs): file about project documentations such as diagram, sequence diagram, etc.
* [internal](internal): internal code that can't be shared.
  * [internal/adapters/grpchandler](internal/adapters/grpchandler): adapter layer that serve into gRPC service.
  * [internal/di](internal/di): dependencies injection for connecting each layer.
* [pkg](pkg): package code that can be shared.
* [scripts](scripts): shell script, go script to help build or testing something.
* [tools](tools): package that need to store on go.mod in order to easily do installation.

## GitHub Actions CI

This project has GitHub Actions CI to do some automation such as:

* [lint](.github/workflows/lint.yml): check the code style.
* [test](.github/workflows/test.yml): run unit testing and uploaded code coverage artifact.
* [generate-proto](.github/workflows/generate-proto.yml): generates protobuf files.
* [generate-rpc-diagram](.github/workflows/generate-rpc-diagram.yml): generates RPC sequence diagram.
* [generate-diagram](.github/workflows/generate-diagram.yml): generates graph code visualization.
* [push-file](.github/workflows/push-file.yml): commit and push generated proto, diagram as github-actions[bot] user.

## Documentation

### Visualize Code Diagram

To help give a better understanding about reading the code
such as relations with packages and types, here are some diagrams listed
generated automatically using [https://github.com/ofabry/go-callvis](https://github.com/ofabry/go-callvis)

<!-- start diagram doc -->
1. [main diagram](docs/diagrams/main.png)
2. [di diagram](docs/diagrams/di.png)
3. [handler diagram](docs/diagrams/handler.png)

<!-- end diagram doc -->
