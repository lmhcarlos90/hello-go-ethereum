# testctx

testctx is a Go package with useful methods for simplifying creating test
contexts. The methods ensure a context created within a test is correctly
canceled once the test is completed, to avoid leaking any resources created
during the test.

## Requiring this package

```
$ go get matheusd.com/testctx@latest
```

## Example usage

```go
package mypkg

import (
    "testing"

    "matheusd.com/testctx"
)

func testFunc(ctx context.Context) {
    <-ctx.Done()
}

func TestMyFunc(t *testing.T) {
    t.Parallel()
    testFunc(testctx.New(t))
    
    // testFunc gets GC'd once the test completes.
}
```

## License

This package is licensed under the [copyfree](http://copyfree.org) ISC License.
