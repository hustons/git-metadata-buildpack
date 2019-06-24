package metadata

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/buildpack/libbuildpack/buildplan"
	"github.com/cloudfoundry/libcfbuildpack/build"
	"github.com/cloudfoundry/libcfbuildpack/layers"
)

const (
	Dependency = "git-metadata"
)

type GitMetadata struct {
	Sha    string `toml:"sha"`
	Branch string `toml:"branch"`
	Remote string `toml:"remote"`
}

func (md GitMetadata) Identity() (string, string) {
	return md.Sha, Dependency
}

func Contribute(context build.Build) error {

	dependency, wantLayer := context.BuildPlan[Dependency]
	if !wantLayer {
		return errors.New(fmt.Sprintf("layer %s is not wanted", Dependency))
	}

	layer := context.Layers.Layer(Dependency)


	md := GitMetadata{
		Sha: dependency.Metadata["sha"].(string),
		Branch: dependency.Metadata["branch"].(string),
		Remote: dependency.Metadata["remote"].(string),
	}

	var metadataDependencyLayerContributor = func(layer layers.Layer) error {
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

	if err := layer.Contribute(md, metadataDependencyLayerContributor, flags(dependency)...); err != nil {
		layer.Logger.Error("Failed to contribute layer [%v]", err)
		return err
	}

	return nil
}

func flags(plan buildplan.Dependency) []layers.Flag {
	var flags []layers.Flag
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
