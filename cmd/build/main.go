package main

import (
	"fmt"
	"github.com/bstick12/git-metadata-buildpack/internal"
	"os"

	"github.com/bstick12/git-metadata-buildpack/metadata"
	"github.com/buildpack/libbuildpack/buildplan"
	"github.com/pkg/errors"

	"github.com/cloudfoundry/libcfbuildpack/build"
)

const (
	FailureStatusCode = 103
)

var CmdRunner = internal.CmdRunner

func main() {

	context, err := build.DefaultBuild()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to create a default build context: %s", err)
		os.Exit(101)
	}

	code, err := RunBuild(context)
	if err != nil {
		context.Logger.Info(err.Error())
	}

	os.Exit(code)

}


func RunBuild(context build.Build) (int, error) {
	context.Logger.FirstLine(context.Logger.PrettyIdentity(context.Buildpack))

	err := metadata.Contribute(context)
	if err != nil {
		return context.Failure(FailureStatusCode), errors.Errorf("Failed to find build plan to create Contributor for %s - [%v]", metadata.Dependency, err)

	}

	return context.Success(buildplan.BuildPlan{})
}
