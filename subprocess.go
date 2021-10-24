package subprocess

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
)

// Subprocess represents a monitored process executed by the application.
type Subprocess struct {
	exitCode int       // The exit code of the process.
	stdout   []byte    // The bytes written to stdout.
	stderr   []byte    // The bytes written to stderr.
	cmd      string    // The name of the command to be executed.
	process  *exec.Cmd // The underlying *exec.Cmd that represents the subprocess.

	args       []string // The sanitized command arguments.
	hideStderr bool     // Hide stderr output.
	hideStdout bool     // Hide stdout output.
	shell      bool     // Executes the command directly in the shell with unsanitization.
	context    string   // Where to execute the subprocess.
}

// Option is a configuration argument for a subprocess.
type Option func(s *Subprocess)

// Subprocess options.
var (
	// Arg adds sanitized argument to command.
	Arg = func(arg string) Option {
		return func(s *Subprocess) {
			s.args = append(s.args, arg)
		}
	}
	// Context determines where the subprocess will be executed.
	// A relative or absolute path may be provided.
	Context = func(path string) Option {
		return func(s *Subprocess) {
			s.context = path
		}
	}
	// Silent hides all output from the subprocess.
	Silent Option = func(s *Subprocess) {
		s.hideStdout = true
		s.hideStderr = true
	}
	// HideStout hides the stdout output of the subprocess.
	HideStdout Option = func(s *Subprocess) {
		s.hideStdout = true
	}
	// HideStderr hides the stder output of the subprocess.
	HideStderr Option = func(s *Subprocess) {
		s.hideStderr = true
	}
	// Shell determines whether the command will directly be ran in the shell
	// without paramater sanitization.
	Shell Option = func(s *Subprocess) {
		s.shell = true
	}
)

// New creates a new Subprocess with the default exit code of 1.
func New(cmd string, opts ...Option) *Subprocess {
	s := &Subprocess{
		exitCode: -1,
		cmd:      cmd,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// ExitCode returns the exit code of the subprocess.
// If the process was terminated by a signal or has not finished.
func (s *Subprocess) ExitCode() int {
	return s.exitCode
}

// IsFinished returns true if the process has finished.
func (s *Subprocess) IsFinished() bool {
	return s.process.ProcessState.Exited()
}

// Stderr returns the bytes that the process has sent to stderr.
func (s *Subprocess) Stderr() []byte {
	return s.stderr
}

// StderrText returns the bytes that the process has sent to stderr.
// The bytes are encoded in a new string.
func (s *Subprocess) StderrText() string {
	return string(s.stderr)
}

// Stdout returns the bytes that the process has sent to stdout.
func (s *Subprocess) Stdout() []byte {
	return s.stdout
}

// StdoutText returns the bytes that the process has sent to stdout.
// The bytes are encoded in a new string.
func (s *Subprocess) StdoutText() string {
	return string(s.stdout)
}

// Exec starts the subprocess.
func (s *Subprocess) Exec() error {
	spawner, osName, err := spawnerFromOS()
	if err != nil {
		return fmt.Errorf("operating system %s not found", osName)
	}

	cmd, err := spawner.CreateCommand(s.cmd, s.args, s.shell, osName)
	if err != nil {
		return err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	wd, err := os.Getwd()
	chwd := err == nil && s.context != ""

	if chwd {
		os.Chdir(s.context)
	}

	cmd.Start()

	readBytes(stdout, func(b []byte) {
		s.stdout = append(s.stdout, b...)
		if !s.hideStdout {
			fmt.Print(string(b))
		}
	})
	readBytes(stderr, func(b []byte) {
		s.stderr = append(s.stderr, b...)
		if !s.hideStderr {
			fmt.Print(string(b))
		}
	})

	cmd.Wait()

	if chwd {
		os.Chdir(wd)
	}

	s.exitCode = cmd.ProcessState.ExitCode()

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

func readBytes(closer io.ReadCloser, action func([]byte)) {
	sc := bufio.NewScanner(closer)
	sc.Split(bufio.ScanRunes)

	for sc.Scan() {
		b := sc.Bytes()
		action(b)
	}
}
