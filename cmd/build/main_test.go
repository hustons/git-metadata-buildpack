package main_test

import (
	"errors"
	"github.com/bstick12/git-metadata-buildpack/internal"
	"github.com/bstick12/git-metadata-buildpack/metadata"
	"github.com/bstick12/git-metadata-buildpack/utils"
	"io"
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

			metadata.CmdRunner = func (_, _ io.Writer, _ io.Reader, _ string, _ ...string) internal.Runner {
				return &internal.TestRunner {
					Runner: func() error {
						return nil
					},
					CombinedOutputter: func() ([]byte, error) {
						return []byte{}, nil
					},
				}
			}

			factory.Build.BuildPlan = buildplan.BuildPlan{
				metadata.Dependency: buildplan.Dependency{
					Metadata: buildplan.Metadata{
						"launch": true,
					},
				},
			}
			code, err := cmdBuild.RunBuild(factory.Build)
			Expect(err).NotTo(HaveOccurred())
			Expect(code).To(Equal(build.SuccessStatusCode))

			metadataLayer := factory.Build.Layers.Layer(metadata.Dependency)
			Expect(metadataLayer).To(test.HaveLayerMetadata(false, false, true))

		})

		it("should fail if it doesn't contribute", func() {

			defer utils.ResetEnv(os.Environ())
			os.Clearenv()

			metadata.CmdRunner = func (_, _ io.Writer, _ io.Reader, _ string, _ ...string) internal.Runner {
				return &internal.TestRunner {
					Runner: func() error {
						return errors.New("error")
					},
					CombinedOutputter: func() ([]byte, error) {
						return []byte{}, nil
					},
				}
			}

			code, err := cmdBuild.RunBuild(factory.Build)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Failed to find build plan"))
			Expect(code).To(Equal(cmdBuild.FailureStatusCode))
		})
	})

}

