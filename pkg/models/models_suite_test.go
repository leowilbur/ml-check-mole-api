package models_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const baseDBURI = "postgres://molepatrol:molepatrol@localhost:5432"

func TestModels(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "molepatrol/models")
}
