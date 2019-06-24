package main_test

import (
	"github.com/bstick12/git-metadata-buildpack/metadata"
	"github.com/bstick12/git-metadata-buildpack/utils"
	"os"
	"testing"

	"github.com/buildpack/libbuildpack/buildplan"

	cmdBuild "github.com/bstick12/git-metadata-buildpack/cmd/build"
	"github.com/cloudfoundry/libcfbuildpack/build"
	"github.com/cloudfoundry/libcfbuildpack/test"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestUnitDetect(t *testing.T) {
	spec.Run(t, "Build", testDetect, spec.Report(report.Terminal{}))
}

func testDetect(t *testing.T, when spec.G, it spec.S) {

	var factory *test.BuildFactory

	it.Before(func() {
		RegisterTestingT(t)
		factory = test.NewBuildFactory(t)
	})

	when("building source", func() {
		it("should pass if successful", func() {

			defer utils.ResetEnv(os.Environ())
			os.Clearenv()

			factory.Build.BuildPlan = buildplan.BuildPlan{
				metadata.Dependency: buildplan.Dependency{
					Metadata: buildplan.Metadata{
						"build":  false,
						"launch": true,
						"sha":    "7aa636e253c4115df34b1f2fab526739cbf27570",
						"branch": "fork/master",
						"remote": "git@github.com/example/example.git",
					},
				},
			}
			code, err := cmdBuild.RunBuild(factory.Build)
			Expect(err).NotTo(HaveOccurred())
			Expect(code).To(Equal(build.SuccessStatusCode))
			metadataLayer := factory.Build.Layers.Layer(metadata.Dependency)
			Expect(metadataLayer).To(test.HaveLayerMetadata(false, false, true))
			md := metadata.Metadata{}
			metadataLayer.ReadMetadata(&md)
			Expect(md).To(Equal(metadata.Metadata{
				Sha : "7aa636e253c4115df34b1f2fab526739cbf27570",
				Branch: "fork/master",
				Remote: "git@github.com/example/example.git",
			}))

		})

		it("should fail if it doesn't contribute", func() {
			defer utils.ResetEnv(os.Environ())
			os.Clearenv()
			code, err := cmdBuild.RunBuild(factory.Build)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Failed to find build plan"))
			Expect(code).To(Equal(cmdBuild.FailureStatusCode))
		})
	})

}

