# mini

A tiny, dependency-free HTTP client for Go.

mini exposes a straightforward request/response API with no third-party
dependencies. It keeps memory use small, so it is comfortable in embedded
runtimes, and it can be vendored directly into any project.

## Installation

```sh
go get github.com/libnudget/mini@v0.1.0
```

## Usage

```go
import "github.com/libnudget/mini"

resp, err := mini.Get("https://example.com/health")
if err != nil {
    panic(err)
}
fmt.Printf("%d %s\n", resp.StatusCode, resp.Body)
```

With a client:

```go
client := mini.New(
    mini.WithTimeout(5 * time.Second),
    mini.WithHeader("Authorization", "Bearer token"),
)
resp, err := client.Post("https://example.com/api", []byte(`{"a":1}`))
```

## Development

```sh
go vet ./...
go test ./...
gofmt -l .
```

## License

MIT
