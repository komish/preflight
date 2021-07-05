package shell

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("HasLicense", func() {
	var (
		engine      FakePodmanEngine
		has_license HasLicenseCheck
	)

	BeforeEach(func() {
		engine = FakePodmanEngine{
			RunReportStdout: `/licenses`,
			RunReportStderr: "",
		}
	})

	Describe("Checking if license can be found", func() {
		Context("When license is found", func() {
			It("should succeed the check", func() {
				ok, err := has_license.validate(engine, "dummy/image", logger)
				Expect(err).ToNot(HaveOccurred())
				Expect(ok).To(BeTrue())
			})
		})
		Context("When license is not found", func() {
			BeforeEach(func() {
				engine.RunReportStdout = "No such file or directory"
			})
			It("should not succeed the check", func() {
				ok, err := has_license.validate(engine, "dummy/image", logger)
				Expect(err).ToNot(HaveOccurred())
				Expect(ok).To(BeFalse())
			})
		})
	})
})
