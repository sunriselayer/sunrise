package e2e

import (
	"time"
)

// convertTokens executes the token conversion command
func (s *E2eTestSuite) convertTokens(amount string) ([]byte, error) {
	return s.execDockerCommand("sunrised", "tx", "tokenconverter", "convert", amount, "--from", KeyUser, "--keyring-backend", "test", "-y")
}

// TestTokenConverter tests the token converter functionality
func (s *E2eTestSuite) TestTokenConverter() {
	// Get initial balance
	initialBalance, err := s.getBalance("validator")
	s.Require().NoError(err)

	// Test cases for successful conversions
	tests := []struct {
		name     string
		amount   string
		expected string
	}{
		{
			name:     "Convert small amount",
			amount:   "1000000",
			expected: "code: 0",
		},
		{
			name:     "Convert medium amount",
			amount:   "10000000",
			expected: "code: 0",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			// Convert tokens
			output, err := s.convertTokens(tc.amount)
			s.Require().NoError(err)
			s.Require().Contains(string(output), tc.expected)

			// Wait for transaction to be included in a block
			time.Sleep(2 * time.Second)

			// Get new balance
			newBalance, err := s.getBalance("validator")
			s.Require().NoError(err)

			// Verify balance has changed
			s.Require().NotEqual(string(initialBalance), string(newBalance))
		})
	}
}

func (s *E2eTestSuite) TestTokenConverterErrorCases() {
	// Test cases for error scenarios
	tests := []struct {
		name     string
		amount   string
		expected string
	}{
		{
			name:     "Convert with insufficient balance",
			amount:   "1000000000000",
			expected: "insufficient funds",
		},
		{
			name:     "Convert with zero amount",
			amount:   "0",
			expected: "amount must be positive",
		},
		{
			name:     "Convert with negative amount",
			amount:   "-1000000",
			expected: "invalid amount",
		},
		{
			name:     "Convert with invalid amount format",
			amount:   "invalid",
			expected: "invalid amount",
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			output, err := s.convertTokens(tc.amount)
			s.Require().Error(err)
			s.Require().Contains(string(output), tc.expected)
		})
	}
}
