package subprocess

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

// Subprocess represents a monitored process executed by the application.
type Subprocess struct {
	ExitCode   int
	hideOutput bool
	executable string
	cmd        *exec.Cmd
	shell      bool
	CommandConfig
}

// CommandConfig is the configuration of a command.
type CommandConfig struct {
	Args    []string // Args is the subprocess command arguments.
	Options []Option // Options is the subprocess options.
	Context string   // Context is where the subprocesses will be executed. By default the process will execute where the binary lies.
}

// Args represents a list of command arguments.
type Args []string

// Option is a configuration argument for a subprocess.
type Option func(s *Subprocess)

// HideOutput hides the output of the subprocess.
var HideOutput Option = func(s *Subprocess) {
	s.hideOutput = true
}

// Shell determines whether the command will directly be ran in the shell
// without paramater sanitization.
var Shell Option = func(s *Subprocess) {
	s.shell = true
}

// New creates a new Subprocess with the default exit code of 1.
func New(cmd string, config CommandConfig) *Subprocess {
	s := &Subprocess{
		ExitCode:      -1,
		executable:    cmd,
		CommandConfig: config,
	}
	for _, opt := range config.Options {
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

	cmd, err := spawner.CreateCommand(s.executable, s.Args, s.shell, osName)
	if err != nil {
		return err
	}

	cmd.Stderr = os.Stdout

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	wd, wdErr := os.Getwd()
	if wdErr == nil && s.Context != "" {
		os.Chdir(s.Context)
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

	if wdErr == nil {
		os.Chdir(wd)
	}

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
