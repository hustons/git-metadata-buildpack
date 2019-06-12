package utils_test

import (
	"os"
	"testing"

	"github.com/bstick12/git-metadata-buildpack/utils"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestUnitDetect(t *testing.T) {
	spec.Run(t, "Utils", testUtils, spec.Report(report.Terminal{}))
}

func testUtils(t *testing.T, when spec.G, it spec.S) {

	it.Before(func() {
		RegisterTestingT(t)
	})

	when("reseting environment", func() {
		it("has the same value prior to reset", func() {
			envVars := os.Environ()
			os.Clearenv()
			Expect(os.Environ()).To(HaveLen(0))
			utils.ResetEnv(envVars)
			Expect(os.Environ()).To(Equal(envVars))
		})
	})
}
