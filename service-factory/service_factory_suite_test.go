package service_factory_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestServiceFactory(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Service Factory Suite")
}
