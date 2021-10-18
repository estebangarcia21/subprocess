package subprocess_test

import (
	"os"
	"runtime"
	"testing"

	"github.com/estebangarcia21/subprocess"
)

func TestSubprocess(t *testing.T) {
	var cmdStr string

	if runtime.GOOS == "windows" {
		cmdStr = "dir"
	} else {
		cmdStr = "ls"
	}

	var opts []subprocess.Option
	if val, _ := os.LookupEnv("SHOW_TEST_SUBPROCESS_OUTPUT"); val != "true" {
		opts = append(opts, subprocess.HideOutput)
	}

	sp := subprocess.New(opts...)

	err := sp.Start(cmdStr)

	if sp.ExitCode != 0 {
		t.Fatalf("wanted exit code 0; got %d", sp.ExitCode)
	}

	if err != nil {
		t.Fatalf("received error while executing subprocess: %v", err)
	}
}
