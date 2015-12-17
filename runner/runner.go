package runner

import (
	"io"
	"os/exec"
)

func run() bool {
	runnerLog("Running...")

	var cmd *exec.Cmd
	runType := runType()

	if runType == "test" {
		cmd = exec.Command("go", "test")
	} else if runType == "python" {
		cmd = exec.Command("python", "./runner.py")
	} else {
		cmd = exec.Command(buildPath())
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		fatal(err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fatal(err)
	}

	err = cmd.Start()
	if err != nil {
		fatal(err)
	}

	go io.Copy(appLogWriter{}, stderr)
	go io.Copy(appLogWriter{}, stdout)

	go func() {
		<-stopChannel
		pid := cmd.Process.Pid
		runnerLog("Killing PID %d", pid)
		cmd.Process.Kill()
	}()

	return true
}
