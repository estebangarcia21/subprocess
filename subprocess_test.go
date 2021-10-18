package subprocess_test

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/estebangarcia21/subprocess"
)

func TestExec(t *testing.T) {
	subprocessTest(t, testSubprocess{
		WindowsCmd: "dir",
		LinuxCmd:   "ls",
		TestFunc: func(s *subprocess.Subprocess) error {
			return s.Exec()
		},
	})
}

func TestExecAsync(t *testing.T) {
	subprocessTest(t, testSubprocess{
		WindowsCmd: "dir",
		LinuxCmd:   "ls",
		TestFunc: func(s *subprocess.Subprocess) error {
			return <-s.ExecAsync()
		},
	})
}

type testSubprocess struct {
	WindowsCmd string
	LinuxCmd   string
	TestFunc   testFunc
}

type testFunc func(s *subprocess.Subprocess) error

func subprocessTest(t *testing.T, testOpts testSubprocess) {
	var cmdStr string

	if runtime.GOOS == "windows" {
		cmdStr = testOpts.WindowsCmd
	} else {
		cmdStr = testOpts.LinuxCmd
	}

	var opts []subprocess.Option
	val, _ := os.LookupEnv("SHOW_TEST_SUBPROCESS_OUTPUT")

	showSubprocessOutput := val == "true"
	if !showSubprocessOutput {
		opts = append(opts, subprocess.HideOutput)
	}

	sp := subprocess.New(cmdStr, opts...)

	if showSubprocessOutput {
		logTitle("Subprocess Output Begin")
	}

	err := testOpts.TestFunc(sp)

	if showSubprocessOutput {
		logTitle("Subprocess Output End")
	}

	if sp.ExitCode != 0 {
		t.Fatalf("wanted exit code 0; got %d", sp.ExitCode)
	}

	if err != nil {
		t.Fatalf("received error while executing subprocess: %v", err)
	}

}

func logTitle(msg string) {
	div := "========================================"
	divLen := len(div)

	msgLen := len(msg)
	msgStart := (divLen - msgLen) / 2

	var midStr string
	for n := 0; n < msgStart; n++ {
		midStr += " "
	}

	midStr += strings.ToUpper(msg)

	fmt.Println(div)
	fmt.Println(midStr)
	fmt.Println(div)
}
