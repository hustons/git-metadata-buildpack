package metadata

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/bstick12/git-metadata-buildpack/internal"

	"github.com/BurntSushi/toml"
	"github.com/buildpack/libbuildpack/buildplan"
	"github.com/cloudfoundry/libcfbuildpack/build"
	"github.com/cloudfoundry/libcfbuildpack/layers"
)

const (
	Dependency = "git-metadata"
)

type Metadata struct {
	Sha    string `toml:"sha"`
	Branch string `toml:"branch"`
	Remote string `toml:"remote"`
}

func (md Metadata) Identity() (string, string) {
	return md.Sha, Dependency
}

var CmdRunner = internal.CmdRunner

func Contribute(context build.Build) error {
	err := CmdRunner(ioutil.Discard, ioutil.Discard, nil, "git", "status").Run()
	if err != nil {
		return errors.New("not identified as git project")
	}

	dependency, wantLayer := context.BuildPlan[Dependency]
	if !wantLayer {
		return errors.New(fmt.Sprintf("layer %s is not wanted", Dependency))
	}

	layer := context.Layers.Layer(Dependency)

	md := Metadata{}

	layer.Logger.SubsequentLine("Retrieving GIT Metadata")
	gitsha, err := CmdRunner(nil, nil, nil, "git", "rev-parse", "HEAD").CombinedOutput()
	if err != nil {
		layer.Logger.Error("Failed to get git SHA [%v]", err)
		return err
	}

	md.Sha = strings.TrimSuffix(string(gitsha), "\n")
	branch, err := CmdRunner(nil, nil, nil, "git", "rev-parse", "--abbrev-ref", "--symbolic-full-name", "@{u}").CombinedOutput()
	if err != nil {
		layer.Logger.Error("Failed to get git branch [%v]", err)
		return err
	}

	md.Branch = strings.TrimSuffix(string(branch), "\n")

	splitBranch := strings.Split(md.Branch, "/")
	remote, err := CmdRunner(nil, nil, nil, "git", "remote", "get-url", splitBranch[0]).CombinedOutput()
	if err != nil {
		layer.Logger.Error("Failed to get git remote url [%v]", err)
		return err
	}
	md.Remote = strings.TrimSuffix(string(remote), "\n")

	var metadataHelperLayerContributor = func(layer layers.Layer) error {
		layer.Touch()
		l := layer.Layer
		filename := filepath.Join(l.Root, "git-metadata.toml")
		if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
			return err
		}

		file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		defer file.Close()

		return toml.NewEncoder(file).Encode(md)
	}

	if err := layer.Contribute(md, metadataHelperLayerContributor, flags(dependency)...); err != nil {
		layer.Logger.Error("Failed to contribute layer [%v]", err)
		return err
	}

	return nil
}

func flags(plan buildplan.Dependency) []layers.Flag {
	flags := []layers.Flag{}
	cache, _ := plan.Metadata["cache"].(bool)
	if cache {
		flags = append(flags, layers.Cache)
	}
	build, _ := plan.Metadata["build"].(bool)
	if build {
		flags = append(flags, layers.Build)
	}
	launch, _ := plan.Metadata["launch"].(bool)
	if launch {
		flags = append(flags, layers.Launch)
	}
	return flags
}
