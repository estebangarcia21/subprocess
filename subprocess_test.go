package subprocess_test

import (
	"fmt"
	"math/rand"
	"os/exec"
	"runtime"
	"strings"
	"testing"

	"github.com/estebangarcia21/subprocess"
)

func TestExec(t *testing.T) {
	crossPlatformTestMatrix{
		"windows": subprocess.New("dir", subprocess.HideStdout),
		"darwin":  subprocess.New("ls", subprocess.Arg("-lh"), subprocess.HideStdout),
		"linux":   subprocess.New("ls", subprocess.Arg("-lh"), subprocess.HideStdout),
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
	u := randString(16)
	crossPlatformTestMatrix{
		"windows": subprocess.New(fmt.Sprintf("Write-Host %s", u), subprocess.HideStdout),
		"darwin":  subprocess.New(fmt.Sprintf("echo %s", u), subprocess.Shell, subprocess.HideStdout),
		"linux":   subprocess.New(fmt.Sprintf("echo %s", u), subprocess.Shell, subprocess.HideStdout),
	}.Exec(func(s *subprocess.Subprocess) {
		if err := s.Exec(); err != nil {
			t.Fatalf("received error while executing subprocess: %v", err)
		}

		stdout := s.StdoutText()

		if i := strings.Index(stdout, u); i == -1 {
			t.Fatalf("expected to find \"%s\" in the subprocess stdout", u)
		}
	})
}

func TestStderrText(t *testing.T) {
	u := randString(16)
	crossPlatformTestMatrix{
		"windows": subprocess.New(fmt.Sprintf("Write-Error %s", u), subprocess.HideStderr),
		"darwin":  subprocess.New(fmt.Sprintf(">&2 echo %s", u), subprocess.Shell, subprocess.HideStderr),
		"linux":   subprocess.New(fmt.Sprintf(">&2 echo %s", u), subprocess.Shell, subprocess.HideStderr),
	}.Exec(func(s *subprocess.Subprocess) {
		if err := s.Exec(); err != nil {
			t.Fatalf("received error while executing subprocess: %v", err)
		}

		stderr := s.StderrText()

		if i := strings.Index(stderr, u); i == -1 {
			t.Fatalf("expected to find \"%s\" in the subprocess stderr", u)
		}
	})
}

func TestRandString(t *testing.T) {
	strLen := 16

	a := randString(strLen)
	var b string

	for n := 0; ; n++ {
		if n > 99 {
			t.Fatal("100 invocations of randString did not produce different results")
		}
		b = randString(strLen)
		if a != b {
			break
		}
	}

	if len(a) != strLen && len(b) != strLen {
		t.Fatalf("strings are not of specified length: %d; a=%d; b=%d", strLen, len(a), len(b))
	}
}

func BenchmarkSubprocessLs(b *testing.B) {
	for n := 0; n < b.N; n++ {
		s := subprocess.New("ls", subprocess.HideStdout)
		s.Exec()
		s.Stdout()
	}
}

func BenchmarkExecCommandLs(b *testing.B) {
	for n := 0; n < b.N; n++ {
		exec.Command("ls").Output()
	}
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

var randTokens = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = randTokens[rand.Intn(len(randTokens))]
	}
	return string(b)
}
