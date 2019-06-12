package internal

import (
	"io"
	"os/exec"
)


type CmdFunctionParams struct {
	Stdout io.Writer
	StdErr io.Writer
	Stdin io.Reader
	Command string
	Args []string
	Return error
	Output []byte
}

type TestRunner struct {
	Runner func () error
	CombinedOutputter func () ([]byte, error)
}

func(tr *TestRunner) Run() error {
	return tr.Runner()
}

func(tr *TestRunner) CombinedOutput() ([]byte, error) {
	return tr.CombinedOutputter()
}

type Runner interface {
	Run() error
	CombinedOutput() ([]byte, error)
}

func CmdRunner(stdout, stderr io.Writer, stdin io.Reader, command string, args ...string) Runner {
	cmd := exec.Command(command, args...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.Stdin = stdin
	return cmd
}

