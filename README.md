# Subprocesses

Spawn subprocesses in Go.

```
go get -u "github.com/estebangarcia21/subprocess"
```

## Sanitized mode

```go
package main

import (
	"log"

	"github.com/estebangarcia21/subprocess"
)

func main() {
	s := subprocess.New("ls", subprocess.Arg("-lh"))

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
