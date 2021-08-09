package auth_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"bitbucket.org/meditekdevsteam/ml-check-mole-api/pkg/auth"
)

var _ = Describe("CognitoJWK", func() {
	It("should return at least 2 keys", func() {
		result, err := auth.CognitoJWK("")
		Expect(err).To(BeNil())
		Expect(len(result.Keys) >= 2).To(BeTrue())
	})
})
