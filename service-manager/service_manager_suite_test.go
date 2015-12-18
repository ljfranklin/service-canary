package service_manager_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestServiceManager(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Service Manager Suite")
}
