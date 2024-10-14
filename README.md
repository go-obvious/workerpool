# Obvious `workpool` implementation

[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-2.1-4baaaa.svg)](CODE-OF-CONDUCT.md)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
![GitHub release](https://img.shields.io/github/release/go-obvious/env.svg)

No frills workerpool implementation.

## How to Use

### Installation

```sh
go get github.com/go-obvious/workerpool
```

### Example Usage

```go
package main

import (
    "fmt"
    "github.com/go-obvious/workerpool"
)

func main() {
    pool := workerpool.New(10, 10)
    pool.Start()
    defer pool.Stop()

    var wg sync.WaitGroup
    values := make([]int, itemCount+1)
    for i := 0; i < itemCount; i++ {
        i := i
        pipevalue := func() {
            defer wg.Done()
            values[i] = i
        }
        wg.Add(1)
        inst.WorkCh <- pipevalue
    }
    wg.Wait()

    // workers are done, print the values
    for entry := range values {
        fmt.Printf("%d\n", entry)
    }
}
```
