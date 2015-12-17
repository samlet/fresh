package runner

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

func build() (string, bool) {
	buildLog("Building...")

	var cmd *exec.Cmd
	var runType = runType()
	if runType == "make" {
		cmd = exec.Command("make")
	} else if runType == "python" {
		cmd = exec.Command("python", "./pre.py")
	} else {
		cmd = exec.Command("go", "build", "-o", buildPath(), root())
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

	io.Copy(os.Stdout, stdout)
	errBuf, _ := ioutil.ReadAll(stderr)

	err = cmd.Wait()
	if err != nil {
		return string(errBuf), false
	}

	return "", true
}
