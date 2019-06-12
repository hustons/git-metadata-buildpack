package metadata_test

import (
	"errors"
	"fmt"
	"github.com/bstick12/git-metadata-buildpack/internal"
	"github.com/bstick12/git-metadata-buildpack/metadata"
	"github.com/bstick12/git-metadata-buildpack/utils"
	"github.com/buildpack/libbuildpack/buildplan"
	"github.com/cloudfoundry/libcfbuildpack/test"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestUnitDetect(t *testing.T) {
	spec.Run(t, "Contributed", testContribute, spec.Report(report.Terminal{}))
}

var cmdFunctions map[string]internal.CmdFunctionParams

func testContribute(t *testing.T, when spec.G, it spec.S) {

	var factory *test.BuildFactory

	var statusArgs = []string{"status"}
	var shaArgs = []string{"rev-parse", "HEAD"}
    var	branchArgs = []string{"rev-parse", "--abbrev-ref", "--symbolic-full-name", "@{u}"}
    var	remoteArgs = []string{"remote", "get-url", "fork"}

	it.Before(func() {
		RegisterTestingT(t)
		factory = test.NewBuildFactory(t)

		factory.Build.BuildPlan = buildplan.BuildPlan{
			metadata.Dependency: buildplan.Dependency{
				Metadata: buildplan.Metadata{
					"build": false,
					"launch": true,
				},

			},
		}
	})

	when("workspace is GIT", func() {

		it.Before(func() {
			cmdFunctions = make(map[string]internal.CmdFunctionParams)
			cmdFunctions[lookupKey(statusArgs)] = internal.CmdFunctionParams{
				Stdout: ioutil.Discard,
				StdErr: ioutil.Discard,
				Stdin: nil,
				Args: statusArgs,
				Return: nil,
			}
			cmdFunctions[lookupKey(shaArgs)] = internal.CmdFunctionParams{
				Stdout: nil,
				StdErr: nil,
				Stdin: nil,
				Args: shaArgs,
				Return: nil,
				Output: []byte("7aa636e253c4115df34b1f2fab526739cbf27570\n"),
			}
			cmdFunctions[lookupKey(branchArgs)] = internal.CmdFunctionParams{
				Stdout: nil,
				StdErr: nil,
				Stdin: nil,
				Args: branchArgs,
				Return: nil,
				Output: []byte("fork/master\n"),
			}
			cmdFunctions[lookupKey(remoteArgs)] = internal.CmdFunctionParams{
				Stdout: nil,
				StdErr: nil,
				Stdin: nil,
				Args: remoteArgs,
				Return: nil,
			}
			})

		it("should execute required commands", func() {
			defer utils.ResetEnv(os.Environ())
			os.Clearenv()
			metadata.CmdRunner = CmdSuccess
			err := metadata.Contribute(factory.Build)
			Expect(err).To(BeNil())
		})

		when("commands fail", func() {
			it("should fail to get status", func() {
				ret := errors.New("not identified as git project")
				changeCmdReturn(lookupKey(statusArgs), ret)
				defer utils.ResetEnv(os.Environ())
				os.Clearenv()

				metadata.CmdRunner = CmdFailure
				err := metadata.Contribute(factory.Build)
				Expect(ret).To(Equal(err))
			})

			it("should fail to get SHA", func() {
				ret := errors.New("Failed to get git SHA")
				changeCmdReturn(lookupKey(shaArgs), ret)
				defer utils.ResetEnv(os.Environ())
				os.Clearenv()

				metadata.CmdRunner = CmdFailure
				err := metadata.Contribute(factory.Build)
				Expect(ret).To(Equal(err))
			})

			it("should fail to get branch", func() {
				ret := errors.New("Failed to get branch")
				changeCmdReturn(lookupKey(branchArgs), ret)
				defer utils.ResetEnv(os.Environ())
				os.Clearenv()

				metadata.CmdRunner = CmdFailure
				err := metadata.Contribute(factory.Build)
				Expect(ret).To(Equal(err))
			})

			it("should fail to get remote url", func() {
				ret := errors.New("Failed to get remote url")
				changeCmdReturn(lookupKey(remoteArgs), ret)
				defer utils.ResetEnv(os.Environ())
				os.Clearenv()

				metadata.CmdRunner = CmdFailure
				err := metadata.Contribute(factory.Build)
				Expect(ret).To(Equal(err))
			})

		})
	})

}

func changeCmdReturn(command string, ret error) {
	cmdFunction := cmdFunctions[command]
	cmdFunction.Return = ret
	cmdFunctions[command] = cmdFunction

}

func CmdSuccess (stdout, stderr io.Writer, stdin io.Reader, command string, args ...string) internal.Runner {
	return &internal.TestRunner {
		CombinedOutputter: func() ([]byte, error) {
			return checkCommand(stdout, stderr, stdin, command, args)
		},
		Runner: func() (err error) {
			_, err = checkCommand(stdout, stderr, stdin, command, args)
			return
		},
	}
}

func checkCommand (stdout, stderr io.Writer, stdin io.Reader, command string, args []string) ([]byte, error) {
	description := fmt.Sprintf("%s with args %q", command, args)
	cmdFunction, ok := cmdFunctions[lookupKey(args)]
	Expect(ok).To(BeTrue(), fmt.Sprintf("Failed to find command %s ", description))
	isEqual(stdout, cmdFunction.Stdout, "stdout", description)
	isEqual(stderr, cmdFunction.StdErr, "stderr", description)
	isEqual(stdin, cmdFunction.Stdin, "stdin", description)
	Expect(args).To(Equal(cmdFunction.Args), description)
	return cmdFunction.Output, cmdFunction.Return
}


func isEqual(actual interface{}, expected interface{}, object string, description string) {
	if expected == nil {
		Expect(actual).To(BeNil(), fmt.Sprintf("%s - %s", object, description))
	} else {
		Expect(actual).To(Equal(expected), fmt.Sprintf("%s - %s", object, description))
	}
}

func CmdFailure (_, _ io.Writer, _ io.Reader, command string, args ...string) internal.Runner {
	return &internal.TestRunner{
		Runner: func() error {
			cmdFunction, ok := cmdFunctions[lookupKey(args)]
			Expect(ok).To(BeTrue(), fmt.Sprintf("Failed to find command %s with args %q", command, args))
			return cmdFunction.Return
		},
		CombinedOutputter: func() ([]byte, error) {
			cmdFunction, ok := cmdFunctions[lookupKey(args)]
			Expect(ok).To(BeTrue(), fmt.Sprintf("Failed to find command %s with args %q", command, args))
			return cmdFunction.Output, cmdFunction.Return
		},
	}
}

func lookupKey(args []string) string {
	return strings.Join(args,"|")
}
