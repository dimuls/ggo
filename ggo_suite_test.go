package ggo

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGgo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Ggo Suite")
}
