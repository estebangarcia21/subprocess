package subprocess

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// spawner is a command located in the system's PATH that begins the
// process. Each major operating system has its own spawner.
//
//
// Each key is PATH command that can open the process. Their position in
// the map determines their priority. For example, if the first
// command is not present in the PATH then it will attempt to use the
// second command and so on.
//
// Each value is an optional set of flags that will be passed to the command.
type spawner map[string][]string

var (
	windowsSpawner = spawner{
		"cmd":  {"powershell -Command"},
		"pwsh": {"powershell -Command"},
	}
	macSpawner = spawner{
		"bash": {"-c"},
		"sh":   {"-c"},
	}
	linuxSpawner = spawner{
		"bash": {"-c"},
		"sh":   {"-c"},
	}
)

// CreateCommand creates an exec.Cmd that is prepared with the root command.
func (s spawner) CreateCommand(cmd string, args string) (*exec.Cmd, error) {
	spawnCmd, err := s.getAvaiableSpawnCommand()
	if err != nil {
		return nil, err
	}
	return exec.Command(spawnCmd, append(s[spawnCmd], "\""+args+"\"")...), nil
}

// getAvaiableSpawnCommand gets the first available spawn command. It returns
// an error if no command is available.
func (s spawner) getAvaiableSpawnCommand() (string, error) {
	var cmd string
	for k := range s {
		if _, err := exec.LookPath(k); err != nil {
			continue
		}
		cmd = k
		break
	}
	if cmd == "" {
		keys := make([]string, 0, len(s))
		for k := range s {
			keys = append(keys, k)
		}
		return "", fmt.Errorf("no available subprocess spawner found in the system PATH. Available spawners for your OS: %s", strings.Join(keys, ", "))
	}
	return cmd, nil
}

// spawnerFromOS returns the proper process spawner and OS (computed from GOOS).
func spawnerFromOS() (spawner, string, error) {
	os := runtime.GOOS
	var s spawner

	switch os {
	case "windows":
		s = windowsSpawner
	case "darwin":
		s = macSpawner
	case "linux":
		s = linuxSpawner
	default:
		return s, os, fmt.Errorf("unsupported operating system %s", os)
	}

	return s, os, nil
}
