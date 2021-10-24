# Subprocesses

Spawn subprocesses in Go.

## Sanitized mode

```go
package main

import (
	"log"

	"github.com/estebangarcia21/subprocess"
)

func main() {
	s := subprocess.New("ls", subprocess.Arg("-lh"), subprocess.Context("/"))

	if err := s.Exec(); err != nil {
		log.Fatal(err)
	}
}
```

## Shell mode

```go
package main

import (
	"log"

	"github.com/estebangarcia21/subprocess"
)

func main() {
	s := subprocess.New("ls -lh", subprocess.Shell)

	if err := s.Exec(); err != nil {
		log.Fatal(err)
	}
}
```
