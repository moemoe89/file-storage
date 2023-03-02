# File Storage Service

---

[![CI Workflow](https://github.com/moemoe89/file-storage/actions/workflows/ci.yml/badge.svg)](https://github.com/moemoe89/file-storage/actions/workflows/ci.yml) <!-- start-coverage --><img src="https://img.shields.io/badge/coverage-98.9%25-brightgreen"><!-- end-coverage -->

File Storage Service handles upload, list and delete related files data into storage.

## Table of Contents

<!-- @import "[TOC]" {cmd="toc" depthFrom=1 depthTo=6 orderedList=false} -->

<!-- code_chunk_output -->

- [Table of Contents](#table-of-contents)
- [Project Summary](#project-summary)
- [Architecture Diagram](#architecture-diagram)
- [Installation](#installation)
    - [1. Set Up Golang Development Environment](#1-set-up-golang-development-environment)
    - [2. Install Development Utility Tools](#2-install-development-utility-tools)
- [Development workflow and guidelines](#development-workflow-and-guidelines)
    - [1. API](#1-api)
    - [2. Object Storage](#2-object-storage)
    - [3. Instrumentation](#3-instrumentation)
    - [4. Unit Test](#4-unit-test)
    - [5. Linter](#5-linter)
    - [6. Mock](#6-mock)
    - [7. Run the service](#7-run-the-service)
    - [8. Test the service](#8-test-the-service)
    - [9. Load Testing](#9-load-testing)
- [CLI](#cli)
- [Project Structure](#project-structure)
- [GitHub Actions CI](#github-actions-ci)
- [Documentation](#documentation)
  - [Visualize Code Diagram](#visualize-code-diagram)
  - [RPC Sequence Diagram](#rpc-sequence-diagram)

<!-- /code_chunk_output -->

## Project Summary

| Item                      | Description                                                                                                          |
|---------------------------|----------------------------------------------------------------------------------------------------------------------|
| Golang Version            | [1.19](https://golang.org/doc/go1.19)                                                                                |
| Object Storage            | [MinIO](https://min.io) and [minio-go](https://github.com/minio/minio-go)                                            |
| moq                       | [mockgen](https://github.com/golang/mock)                                                                            |
| Linter                    | [GolangCI-Lint](https://github.com/golangci/golangci-lint)                                                           |
| Testing                   | [testing](https://golang.org/pkg/testing) and [testify/assert](https://godoc.org/github.com/stretchr/testify/assert) |
| Load Testing              | [ghz](https://ghz.sh)                                                                                                |
| API                       | [gRPC](https://grpc.io/docs/tutorials/basic/go) and [gRPC-Gateway](https://github.com/grpc-ecosystem/grpc-gateway)   |
| CLI                       | [flag](https://pkg.go.dev/flag)                                                                                      |
| Application Architecture  | [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)                   |
| Directory Structure       | [Standard Go Project Layout](https://github.com/golang-standards/project-layout)                                     |
| CI (Lint, Test, Generate) | [GitHubActions](https://github.com/features/actions)                                                                 |
| Visualize Code Diagram    | [go-callviz](https://github.com/ofabry/go-callvis)                                                                   |
| Sequence Diagram          | [Mermaid](https://mermaid.js.org)                                                                                    |
| Protobuf Operations       | [buf](https://buf.build)                                                                                             |
| Instrumentation           | [OpenTelemetry](https://opentelemetry.io) and [Jaeger](https://www.jaegertracing.io)                                 |
| Logger                    | [zap](https://github.com/uber-go/zap)                                                                                |


## Architecture Diagram

---

[Excalidraw link](https://excalidraw.com/#json=Nj7TrMA5pQPIKY4Jze8df,vygEYUvYlvnY1QSksZVAZQ)

![Architecture-Diagram](https://user-images.githubusercontent.com/7221739/221416201-8473c385-de37-486f-b21f-3bf2e626bfb9.png)

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

### 2. Object Storage

![MinIO](https://user-images.githubusercontent.com/7221739/222305823-23040742-d07d-48fa-a8d7-b16221bb6ced.png)


Instead storing the files to disk, this project use [MinIO](https://min.io) as an Object Storage to store the files.
You can follow the installation from the official page, or easily run this docker-compose command to setup everything including the bucket. 

```sh
$ docker-compose -f ./development/docker-compose.yml up minio
```

For create the default bucket:

```sh
$ docker-compose -f ./development/docker-compose.yml up createbuckets
```

For create the backup:

```sh
$ docker-compose -f ./development/docker-compose.yml up createbackup
```

### 3. Instrumentation

![JaegerUI](https://user-images.githubusercontent.com/7221739/222305605-64a77902-8ce9-40b6-9298-5d733f5f1316.png)

This service implements [OpenTelemetry](https://opentelemetry.io/) to enable instrumentation in order to measure the performance.
The data exported to Jaeger and can be seen in the Jaeger UI [http://localhost:16686](http://localhost:16686)

For running the Jaeger exporter, easily run with docker-compose command:

```sh
$ docker-compose -f ./development/docker-compose.yml up jaeger
```

### 4. Unit Test

You can simply execute the following command to run all test cases in this service:

```sh
$ make test
```

### 5. Linter

For running the linter make sure these libraries already installed in your system:

* https://github.com/golangci/golangci-lint
* https://github.com/yoheimuta/protolint

Then checks the Go and Proto code style using lint can be done with this command:

```sh
$ make lint
```

### 6. Mock

This service using Mock in some places like in the repository, usecase, pkg, etc.
To automatically updating the mock if the interface changed, easily run with `go generate` command:

```sh
$ make mock
```

### 7. Run the service

For running the service, you need the database running and set up some env variables:

```
# app config
export APP_ENV=dev
export SERVER_PORT=8080

# minio config
export MINIO_HOST=localhost:9000
export MINIO_ACCESS_KEY_ID=minioadmin
export MINIO_SECRET_ACCESS_KEY=minioadmin

# tracing config
export OTEL_AGENT=http://localhost:14268/api/traces
```

Or you can just execute the sh file:

```sh
$ ./scripts/run.sh
```

### 8. Test the service

The example how to call the gRPC service written in Golang can be seen on these 6 examples:

1. [upload-from-file](scripts/upload-from-file/main.go) file.
2. [upload-from-url](scripts/upload-from-url/main.go) file.
3. [concurrent-upload-file](scripts/concurrent-upload-file/main.go) file.
4. [concurrent-upload-url](scripts/concurrent-upload-url/main.go) file.
5. [list-file](scripts/list-file/main.go) file.
6. [delete-file](scripts/delete-file/main.go) file.

> NOTE: To test this service need MinIO running in order to store the file.
 
If you want to test by GUI client, you can use either BloomRPC (although already no longer active) or Postman.
For the detail please visit these links:
* https://github.com/bloomrpc/bloomrpc
* https://www.postman.com

Basically you just need to import the [api/proto/service.proto](api/proto/service.proto) file if you want to test via BloomRPC / Postman.

> NOTE: There will be a possibility issue when importing the proto file to BloomRPC or Postman.
> It is caused by some path issue, the usage of `gRPC Gateway` and `protoc-gen-validate` library.
> To solve this issue, there's need a modification for the proto file.

#### BloomRPC

![BloomRPC](https://user-images.githubusercontent.com/7221739/222306402-369e6288-229c-43e4-ad2c-d09c9cb465dc.png)

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

![Postman](https://user-images.githubusercontent.com/7221739/222306916-56e28d36-1756-4187-b9c6-1bd0e20eda0d.png)

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

![Swagger](https://user-images.githubusercontent.com/7221739/222307821-2f0844cb-9f76-491d-8a1c-1e8158cdc0e7.png)


This service has HTTP server built on gRPC-Gateway, if you prefer to test using HTTP instead HTTP2 protocol,
you can copy the Swagger file here [api/openapiv2/proto/service.swagger.json](api/openapiv2/proto/service.swagger.json) and then copy paste to this URL https://editor.swagger.io/

By default, HTTP server running on gRPC port + 1, if the gRPC port is 8080, then HTTP server will run on 8081.

# NOTE

> If you have any difficulties to run the service, easily just run all dependencies by docker-compose for the example:
> 
> `docker-compose -f ./development/docker-compose.yml up`
>
> Then you will have all services running like `minio`, `createbuckets`, `createbackup`,`jaeger` and run `file-storage` server.

## CLI

![CLI](https://user-images.githubusercontent.com/7221739/222308036-545f1fa6-c51e-4cd2-ac55-e059d636eb29.png)


This project has CLI to simplify interact with the File-Storage server. You need to build the binary before running the command:

```shell
make build-cli
```

You will have a binary named `fs-store`.

> NOTE:
> If the binary is not in your PATH, you need to run it directly like this: ./fs-store list-files

Then, here are several commands available:

#### File Upload
```shell
// Use default bucket, upload from file path
fs-store file-upload /path/test.txt

// Use default bucket, upload from URL
fs-store -source=url file-upload /path/test.txt

// Use specific bucket
fs-store -bucket=my-bucket file-upload /path/test.txt

// Use specific filename
fs-store -filename=my-file.txt file-upload /path/test.txt

// Use validations
fs-store -content_type=image/jpeg,image/png -max_size=1000 file-upload /path/test.txt
```

#### List File
```shell
// Use default bucket
file-storage list-files

// With specific bucket
file-storage -bucket=my-bucket list-files
```

#### File Upload
```shell
// Use default bucket
file-storage file-delete test.txt

// With specific bucket
file-storage -bucket=my-bucket file-delete test.txt
```

### 9. Load Testing

![ghz](https://user-images.githubusercontent.com/7221739/222305995-964c6f8e-65c1-4788-9ce5-392a617528c1.png)


In order to make sure the service ready to handle a big traffic, it will better if we can do Load Testing to see the performance.

Since the service running in gRPC, we need the tool that support to do HTTP2 request.
In this case we can use https://ghz.sh/ because it is very simple and can generate various output report type.

> NOTE: Like importing the proto file to BloomRPC / Postman,
> when running the `ghz` there's will be issue shown due to the tool can't read the path & validate lib.

Here are some possibility issues when we're trying to run the `ghz` commands:
* `./api/proto/service.proto:5:8: open api/proto/proto/entity.proto: no such file or directory`
* `./api/proto/service.proto:7:8: open api/proto/validate/validate.proto: no such file or directory`

To fix this issue, you need to change some file in proto file:

```protobuf
import "google/api/annotations.proto";
import "protoc-gen-openapiv2/options/annotations.proto";
```

To this:

```protobuf
// import "google/api/annotations.proto";
// import "protoc-gen-openapiv2/options/annotations.proto";
```

and remove gRPC Gateway related annotations:

```protobuf
option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger) = {
  ...
};
```

Then, you can run this `ghz` command to do Load Testing for specific RPC, for the example:

#### 1. Upload RPC:

```shell
ghz --insecure --proto ./api/proto/service.proto --call FileStorageService.Upload -d '{ "type": 1, "filename": "test.txt", "bucket": "default", "file": { "data": "dGVzdA==", "offset": 0 } }' 0.0.0.0:8080 -O html -o load_testing_upload_file.html
```


#### 2. List RPC:

```sh
ghz --insecure --proto ./api/proto/service.proto --call FileStorageService.List -d '{ "bucket": "default" }' 0.0.0.0:8080 -O html -o load_testing_list_files.html
```

#### 3. Delete RPC:

```sh
ghz --insecure --proto ./api/proto/service.proto --call FileStorageService.Delete -d '{ "bucket": "default", "object": "file.txt" }' 0.0.0.0:8080 -O html -o load_testing_delete_file.html
```

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
  * [internal/usecases](internal/usecases): business logic that connect to repository layer, RPC & HTTP client, etc.
* [pkg](pkg): package code that can be shared.
* [scripts](scripts): shell script, go script to help build or testing something.
* [tools](tools): package that need to store on go.mod in order to easily do installation.

## GitHub Actions CI

![GitHubActionsCI](https://user-images.githubusercontent.com/7221739/222308704-51a26273-e398-47d2-9619-a2b8b4fc0988.png)


This project has GitHub Actions CI to do some automation such as:

* [lint](.github/workflows/lint.yml): check the code style.
* [test](.github/workflows/test.yml): run unit testing and uploaded code coverage artifact.
* [generate-proto](.github/workflows/generate-proto.yml): generates protobuf files.
* [generate-diagram](.github/workflows/generate-diagram.yml): generates graph code visualization.
* [push-file](.github/workflows/push-file.yml): commit and push generated proto, diagram as github-actions[bot] user.

## Documentation

### Visualize Code Diagram

![GraphDiagram](https://user-images.githubusercontent.com/7221739/222308530-674b4561-c9fe-4529-acb8-876da7027a18.png)

To help give a better understanding about reading the code
such as relations with packages and types, here are some diagrams listed
generated automatically using [https://github.com/ofabry/go-callvis](https://github.com/ofabry/go-callvis)

<!-- start diagram doc -->
1. [main diagram](docs/diagrams/main.png)
2. [di diagram](docs/diagrams/di.png)
3. [handler diagram](docs/diagrams/handler.png)

<!-- end diagram doc -->

### RPC Sequence Diagram

![SequenceDiagram](https://user-images.githubusercontent.com/7221739/222308588-37ebd33e-84db-43ae-a029-0f78b6c8098e.png)

To help give a better understanding about reading the RPC flow
such as relations with usecases and repositories, here are some sequence diagrams (generated automatically) listed in Markdown file and written in Mermaid JS [https://mermaid-js.github.io/mermaid/](https://mermaid-js.github.io/mermaid/) format.

To generate the RPC sequence diagram, there's a Makefile command that can be use:

1. Run this command to generate specific RPC `make sequence-diagram RPC=GetData`.
2. For generates multiple RPC's, just adding the other RPC by comma `make sequence-diagram RPC=GetData,GetList`.
3. For generates all RPC's, use wildcard * in the parameter `make sequence-diagram RPC=*`.

<!-- start rpc sequence diagram doc -->
1. [Delete RPC - Sequence Diagram](docs/sequence-diagrams/rpc/delete.md)
2. [List RPC - Sequence Diagram](docs/sequence-diagrams/rpc/list.md)
3. [Upload RPC - Sequence Diagram](docs/sequence-diagrams/rpc/upload.md)

<!-- end rpc sequence diagram doc -->
