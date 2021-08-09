package auth_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bitbucket.org/meditekdevsteam/ml-check-mole-api/pkg/auth"
)

var _ = Describe("LookupPermissions", func() {
	It("should allow users to access their own lesions", func() {
		Expect(auth.LookupPermissions([]string{}, "owned.lesions.read")).To(BeTrue())
	})

	It("should forbid users and doctors from writing questions", func() {
		Expect(auth.LookupPermissions([]string{}, "questions.write")).To(BeFalse())
		Expect(auth.LookupPermissions([]string{"Doctors"}, "questions.write")).To(BeFalse())
	})

	It("should allow doctors to respond to requests", func() {
		Expect(auth.LookupPermissions([]string{"Doctors"}, "requests.respond")).To(BeTrue())
	})

	It("should allow administrators to write questions", func() {
		Expect(auth.LookupPermissions([]string{"Administrators"}, "questions.write")).To(BeTrue())
	})

	It("should handle unknown groups", func() {
		Expect(auth.LookupPermissions([]string{"Unknown"}, "questions.write")).To(BeFalse())
	})
})
