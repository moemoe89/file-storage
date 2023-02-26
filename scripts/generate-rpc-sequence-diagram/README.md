# Generate RPC Sequence Diagram

## Requirements

- [Go](https://golang.org/dl/)

## Annotations

To easily find the struct handler and struct usecase, need to add some annotation ended with
`is a struct for handler` and `is a struct for usecase` e.g

```go
// xxxHandler is a struct for handler
type xxxHandler struct {
}
```

```go
// yyyUsecase is a struct for usecase
type yyyUsecase struct {
}
```

## How to run the script

1. Go to root directory.
2. Run this command to generate specific RPC `go run ./scripts/generate-rpc-sequence-diagram -RPC=GetData`.
3. For generates multiple RPC's, just adding the other RPC by comma `go run ./scripts/generate-rpc-sequence-diagram -RPC=GetData,GetList`.
4. For generates all RPC's, use wildcard * in the parameter `go run ./scripts/generate-rpc-sequence-diagram -RPC=*`.
5. Check the [docs/sequence-diagrams/rpc](../../docs/sequence-diagrams/rpc) from the root directory.
