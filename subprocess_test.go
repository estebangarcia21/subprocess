package subprocess_test

import (
	"fmt"
	"os"
	"runtime"
	"strings"
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
	val, _ := os.LookupEnv("SHOW_TEST_SUBPROCESS_OUTPUT")

	showSubprocessOutput := val == "true"
	if !showSubprocessOutput {
		fmt.Println("SHOW_TEST_SUBPROCESS_OUTPUT is enabled. Subprocess STDOUT will be shown...")
		opts = append(opts, subprocess.HideOutput)
	}

	sp := subprocess.New(opts...)

	if showSubprocessOutput {
		logTitle("Subprocess Output Begin")
	}

	err := sp.Start(cmdStr)

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
