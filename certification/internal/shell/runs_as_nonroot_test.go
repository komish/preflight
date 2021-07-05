package shell

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RunAsNonRoot", func() {
	var (
		engine         FakePodmanEngine
		run_as_nonroot RunAsNonRootCheck
	)

	BeforeEach(func() {
		engine = FakePodmanEngine{
			RunReportStdout: "1",
			RunReportStderr: "",
		}
	})

	Describe("Checking runtime user is not root", func() {
		Context("When runtime user is not root", func() {
			It("should succeed the check", func() {
				ok, err := run_as_nonroot.validate(engine, "dummy/image", logger)
				Expect(err).ToNot(HaveOccurred())
				Expect(ok).To(BeTrue())
			})
		})
		Context("When runtime user is root", func() {
			BeforeEach(func() {
				engine.RunReportStdout = "0"
			})
			It("should not succeed the check", func() {
				ok, err := run_as_nonroot.validate(engine, "dummy/image", logger)
				Expect(err).ToNot(HaveOccurred())
				Expect(ok).To(BeFalse())
			})
		})
	})
})
