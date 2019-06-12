package main_test

import (
	"errors"
	"github.com/bstick12/git-metadata-buildpack/internal"
	"github.com/bstick12/git-metadata-buildpack/metadata"
	"io"
	"os"
	"testing"

	"github.com/buildpack/libbuildpack/buildplan"

	cmdDetect "github.com/bstick12/git-metadata-buildpack/cmd/detect"
	"github.com/bstick12/git-metadata-buildpack/utils"
	"github.com/cloudfoundry/libcfbuildpack/detect"
	"github.com/cloudfoundry/libcfbuildpack/test"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestUnitDetect(t *testing.T) {
	spec.Run(t, "Detect", testDetect, spec.Report(report.Terminal{}))
}


func testDetect(t *testing.T, when spec.G, it spec.S) {

	var factory *test.DetectFactory

	it.Before(func() {
		RegisterTestingT(t)
		factory = test.NewDetectFactory(t)
	})

	when("when the workspace is a GIT project", func() {
		it("should add git-metadata-buildpack to the buildplan", func() {
			defer utils.ResetEnv(os.Environ())
			os.Clearenv()

			cmdDetect.CmdRunner = func (_, _ io.Writer, _ io.Reader, _ string, _ ...string) internal.Runner {
				return &internal.TestRunner {
					Runner: func() error {
						return nil
					},
					CombinedOutputter: func() ([]byte, error) {
						return []byte{}, nil
					},
				}
			}

			code, err := cmdDetect.RunDetect(factory.Detect)
			Expect(err).NotTo(HaveOccurred())
			Expect(code).To(Equal(detect.PassStatusCode))
			Expect(factory.Output).To(Equal(buildplan.BuildPlan{
				metadata.Dependency: buildplan.Dependency{
					Metadata: buildplan.Metadata{
						"build":  false,
						"launch": true,
					},
				},
			}))
		})
	})

	when("when the workspace is a not GIT project", func() {
		it("should not add git-metadata-buildpack to the buildplan", func() {
			defer utils.ResetEnv(os.Environ())
			os.Clearenv()

			cmdDetect.CmdRunner = func (_, _ io.Writer, _ io.Reader, _ string, _ ...string) internal.Runner {
				return &internal.TestRunner {
					Runner: func() error {
						return errors.New("failed to invoke git status")
					},
					CombinedOutputter: func() ([]byte, error) {
						return []byte{}, nil
					},
				}
			}

			code, err := cmdDetect.RunDetect(factory.Detect)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("not identified as git project"))
			Expect(code).To(Equal(detect.FailStatusCode))
		})
	})

}
