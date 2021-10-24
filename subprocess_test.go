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
	tests := map[string]struct {
		cmd           string
		commandConfig subprocess.CommandConfig
	}{
		"windows": {
			cmd: "dir",
			commandConfig: subprocess.CommandConfig{
				Context: "/",
			},
		},
		"darwin": {
			cmd: "ls",
			commandConfig: subprocess.CommandConfig{
				Args:    []string{"-lh"},
				Context: "/",
			},
		},
		"linux": {
			cmd: "ls",
			commandConfig: subprocess.CommandConfig{
				Args:    []string{"-lh"},
				Context: "/",
			},
		},
	}

	goos := runtime.GOOS

	for platform, tt := range tests {
		if platform != goos {
			continue
		}

		var opts []subprocess.Option
		val, _ := os.LookupEnv("SHOW_TEST_SUBPROCESS_OUTPUT")

		showSubprocessOutput := val == "true"
		if !showSubprocessOutput {
			opts = append(opts, subprocess.HideOutput)
		}

		sp := subprocess.New(tt.cmd, tt.commandConfig)

		if showSubprocessOutput {
			logTitle("Subprocess Output Begin")
		}

		err := sp.Exec()

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
}

const logTitleDiv = "========================================"

func logTitle(msg string) {
	divLen := len(logTitleDiv)

	msgLen := len(msg)
	msgStart := (divLen - msgLen) / 2

	var midStr string
	for n := 0; n < msgStart; n++ {
		midStr += " "
	}

	midStr += strings.ToUpper(msg)

	fmt.Println(logTitleDiv)
	fmt.Println(midStr)
	fmt.Println(logTitleDiv)
}
