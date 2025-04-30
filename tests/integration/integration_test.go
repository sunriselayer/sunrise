package integration

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// TestIntegrationSuite runs all e2e tests
func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
