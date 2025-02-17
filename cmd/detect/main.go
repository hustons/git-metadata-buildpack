package main

import (
	"errors"
	"fmt"
	"github.com/bstick12/git-metadata-buildpack/internal"
	"io/ioutil"
	"os"
	"strings"

	"github.com/bstick12/git-metadata-buildpack/metadata"

	"github.com/buildpack/libbuildpack/buildplan"

	"github.com/cloudfoundry/libcfbuildpack/detect"
	"github.com/cloudfoundry/libcfbuildpack/logger"

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

	if err != nil {
		return detect.FailStatusCode, errors.New("not identified as git project")
	}

	md, err := GetMetadata(context.Logger)

	if err == nil {
		return context.Pass(buildplan.BuildPlan{
			metadata.Dependency: buildplan.Dependency{
				Metadata: buildplan.Metadata{
					"build":  false,
					"launch": true,
					"sha": md.Sha,
					"branch": md.Branch,
					"remote": md.Remote,
				},
			},
		})
	}

	return detect.FailStatusCode, err

}

func GetMetadata(log logger.Logger) (metadata.GitMetadata, error) {

	md := metadata.GitMetadata{}

	log.SubsequentLine("Retrieving GIT Metadata")
	gitsha, err := CmdRunner(nil, nil, nil, "git", "rev-parse", "HEAD").CombinedOutput()
	if err != nil {
		log.Error("Failed to get git SHA [%v]", err)
		return md, err
	}

	md.Sha = strings.TrimSuffix(string(gitsha), "\n")
	branch, err := CmdRunner(nil, nil, nil, "git", "rev-parse", "--abbrev-ref", "--symbolic-full-name", "@{u}").CombinedOutput()
	if err != nil {
		branch = []byte("origin/DETACHED")
		log.SubsequentLine("Failed to get branch assuming detached HEAD")
	}

	md.Branch = strings.TrimSuffix(string(branch), "\n")

	splitBranch := strings.Split(md.Branch, "/")
	remote, err := CmdRunner(nil, nil, nil, "git", "remote", "get-url", splitBranch[0]).CombinedOutput()
	if err != nil {
		log.Error("Failed to get git remote url [%v]", err)
		return md, err
	}
	md.Remote = strings.TrimSuffix(string(remote), "\n")

	return md, nil

}
