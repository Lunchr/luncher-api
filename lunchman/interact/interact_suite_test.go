package interact_test

import (
	"strings"

	"github.com/deiwin/luncher-api/lunchman/interact"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestInteract(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Interact Suite")
}

var (
	actor     interact.Actor
	userInput string
)

var _ = JustBeforeEach(func(done Done) {
	defer close(done)
	actor = interact.NewActor(strings.NewReader(userInput))
})
