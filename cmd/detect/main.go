package main

import (
	"errors"
	"fmt"
	"github.com/bstick12/git-metadata-buildpack/internal"
	"io/ioutil"
	"os"

	"github.com/bstick12/git-metadata-buildpack/metadata"

	"github.com/buildpack/libbuildpack/buildplan"

	"github.com/cloudfoundry/libcfbuildpack/detect"
)

func main() {
	context, err := detect.DefaultDetect()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create a default detection context: %s", err)
		os.Exit(detect.FailStatusCode)
	}

	if err := context.BuildPlan.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize Build Plan: %s\n", err)
		os.Exit(detect.FailStatusCode)
	}

	code, err := RunDetect(context)
	if err != nil {
		context.Logger.Info(err.Error())
	}

	os.Exit(code)
}

var CmdRunner = internal.CmdRunner

func RunDetect(context detect.Detect) (int, error) {

	err := CmdRunner(ioutil.Discard, ioutil.Discard, nil, "git", "status").Run()

	if err == nil {
		return context.Pass(buildplan.BuildPlan{
			metadata.Dependency: buildplan.Dependency{
				Metadata: buildplan.Metadata{
					"build":  false,
					"launch": true,
				},
			},
		})
	}

	return detect.FailStatusCode, errors.New("not identified as git project")

}
