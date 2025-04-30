package e2e

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// TestE2eSuite runs all e2e tests
func TestE2eSuite(t *testing.T) {
	suite.Run(t, new(E2eTestSuite))
}
