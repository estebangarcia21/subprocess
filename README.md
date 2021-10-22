# Subprocesses

Spawn subprocesses in Go.

## Example

```go
package main

import (
  "fmt"
  "github.com/estebangarcia21/subprocess"
)

func main() {
  s := subprocess.New("ls -lh")

  err := s.Exec()
  if err != nil {
    fmt.Println("Error while executing subprocess ls")
  }

  fmt.Printf("Exit code: %s\n", s.ExitCode)
}
```
