package watcher_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type serviceSuite struct {
	suite.Suite
}

func TestService(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(serviceSuite))
}

// TODO: Write real tests.
func (s *serviceSuite) TestFake() {
	s.True(true) //nolint:testifylint // Temporary
}
