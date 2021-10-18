package subprocess

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Subprocess represents a monitored process executed by the application.
type Subprocess struct {
	ExitCode int       // ExitCode is the exit code of the process. Defaults to -1.
	cmd      *exec.Cmd // cmd is the underlying command being executed.
}

// New creates a new Subprocess.
func New() *Subprocess {
	return &Subprocess{ExitCode: -1}
}

// IsFinished returns true if the process has finished.
func (s *Subprocess) IsFinished() bool {
	return s.cmd.ProcessState.Exited()
}

// Start starts the process.
func (s *Subprocess) Start(cmdStr string) error {
	spawner, osName, err := spawnerFromOS()
	if err != nil {
		return fmt.Errorf("operating system %s not found", osName)
	}

	t := strings.Split(cmdStr, " ")

	cmd, err := spawner.CreateCommand(t[0], strings.Join(t[0:], " "))
	if err != nil {
		return err
	}

	cmd.Stderr = os.Stdout

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	_ = cmd.Start()

	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanRunes)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Print(m)
	}

	_ = cmd.Wait()

	s.ExitCode = cmd.ProcessState.ExitCode()

	return nil
}
