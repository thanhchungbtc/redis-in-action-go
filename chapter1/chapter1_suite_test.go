package chapter1_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestChapter1(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Chapter1 Suite")
}
