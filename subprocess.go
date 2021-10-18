package subprocess

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

// Subprocess represents a monitored process executed by the application.
type Subprocess struct {
	ExitCode   int // ExitCode is the exit code of the process. Defaults to -1.
	hideOutput bool
	cmdStr     string
	cmd        *exec.Cmd // cmd is the underlying command being executed.
}

// Option is a configuration argument for a subprocess.
type Option func(s *Subprocess)

// HideOutput hides the output of the subprocess.
var HideOutput Option = func(s *Subprocess) {
	s.hideOutput = true
}

// New creates a new Subprocess.
func New(cmdStr string, opts ...Option) *Subprocess {
	s := &Subprocess{
		ExitCode:   -1,
		hideOutput: false,
		cmdStr:     cmdStr,
		cmd:        nil,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// IsFinished returns true if the process has finished.
func (s *Subprocess) IsFinished() bool {
	return s.cmd.ProcessState.Exited()
}

// Exec starts the subprocess.
func (s *Subprocess) Exec() error {
	spawner, osName, err := spawnerFromOS()
	if err != nil {
		return fmt.Errorf("operating system %s not found", osName)
	}

	cmd, err := spawner.CreateCommand(s.cmdStr)
	if err != nil {
		return err
	}

	cmd.Stderr = os.Stdout

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	_ = cmd.Start()

	if !s.hideOutput {
		scanner := bufio.NewScanner(stdout)
		scanner.Split(bufio.ScanRunes)

		for scanner.Scan() {
			m := scanner.Text()
			fmt.Print(m)
		}
	}

	_ = cmd.Wait()

	s.ExitCode = cmd.ProcessState.ExitCode()

	return nil
}

// ExecAsync starts the subprocess asynchronously.
func (s *Subprocess) ExecAsync() chan error {
	ch := make(chan error)
	go func(s *Subprocess) {
		ch <- s.Exec()
	}(s)
	return ch
}
