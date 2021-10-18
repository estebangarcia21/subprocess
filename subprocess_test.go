package subprocess_test

import (
	"runtime"
	"testing"

	"github.com/estebangarcia21/subprocess"
)

func TestSubprocess(t *testing.T) {
	var cmdStr string

	if runtime.GOARCH == "windows" {
		cmdStr = "dir"
	} else {
		cmdStr = "ls"
	}

	sp := subprocess.New(subprocess.HideOutput)

	err := sp.Start(cmdStr)

	if sp.ExitCode != 0 {
		t.Fatalf("wanted exit code 0; got %d", sp.ExitCode)
	}

	if err != nil {
		t.Fatalf("received error while executing subprocess: %v", err)
	}
}
