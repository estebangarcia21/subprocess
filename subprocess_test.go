package subprocess_test

import (
	"runtime"
	"strings"
	"testing"

	"github.com/estebangarcia21/subprocess"
)

func TestExec(t *testing.T) {
	crossPlatformTestMatrix{
		"windows": subprocess.New("dir"),
		"darwin":  subprocess.New("ls", subprocess.Arg("-lh")),
		"linux":   subprocess.New("ls", subprocess.Arg("-lh")),
	}.Exec(func(s *subprocess.Subprocess) {
		if err := s.Exec(); err != nil {
			t.Fatalf("received error while executing subprocess: %v", err)
		}

		if s.ExitCode() != 0 {
			t.Fatalf("wanted exit code 0; got %d", s.ExitCode())
		}
	})
}

func TestStdoutText(t *testing.T) {
	crossPlatformTestMatrix{
		"windows": subprocess.New("Write-Host Hello world!"),
		"darwin":  subprocess.New("echo Hello world!", subprocess.Shell),
		"linux":   subprocess.New("echo Hello world!", subprocess.Shell),
	}.Exec(func(s *subprocess.Subprocess) {
		if err := s.Exec(); err != nil {
			t.Fatalf("received error while executing subprocess: %v", err)
		}

		stdout := s.StdoutText()

		if i := strings.Index(stdout, "Hello world!"); i == -1 {
			t.Fatal("expected to find \"Hello world!\" in the subprocess stdout")
		}
	})
}

type crossPlatformTestMatrix map[string]*subprocess.Subprocess

func (c crossPlatformTestMatrix) Exec(test func(*subprocess.Subprocess)) {
	for platform, s := range c {
		if platform == runtime.GOOS {
			test(s)
			break
		}
	}
}
