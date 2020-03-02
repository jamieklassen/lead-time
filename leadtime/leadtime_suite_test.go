package leadtime_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestLeadtime(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Leadtime Suite")
}
