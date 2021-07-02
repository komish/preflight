package shell

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("HasNoProhibitedPackages", func() {
	var (
		engine                     FakePodmanEngine
		has_no_prohibited_packages HasNoProhibitedPackagesCheck
	)

	BeforeEach(func() {
		engine = FakePodmanEngine{
			RunReportStdout: "",
			RunReportStderr: "",
		}
	})

	Describe("Checking if it has an prohibited packages", func() {
		Context("When there are no prohibited packages found", func() {
			It("should succeed the check", func() {
				ok, err := has_no_prohibited_packages.validate(engine, "dummy/image", logger)
				Expect(err).ToNot(HaveOccurred())
				Expect(ok).To(BeTrue())
			})
		})
		Context("When there was a prohibited packages found", func() {
			BeforeEach(func() {
				engine.RunReportStdout = "grub"
			})
			It("should not succeed the check", func() {
				ok, err := has_no_prohibited_packages.validate(engine, "dummy/image", logger)
				Expect(err).ToNot(HaveOccurred())
				Expect(ok).To(BeFalse())
			})
		})
	})
})
