package e2e

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestHelper provides common test utilities
type TestHelper struct {
	t *testing.T
}

// NewTestHelper creates a new TestHelper instance
func NewTestHelper(t *testing.T) *TestHelper {
	return &TestHelper{
		t: t,
	}
}

// WaitForNextBlock waits for the next block to be produced
func (h *TestHelper) WaitForNextBlock() {
	time.Sleep(2 * time.Second)
}

// RequireNoError checks if there is no error
func (h *TestHelper) RequireNoError(err error, msgAndArgs ...interface{}) {
	require.NoError(h.t, err, msgAndArgs...)
}

// RequireError checks if there is an error
func (h *TestHelper) RequireError(err error, msgAndArgs ...interface{}) {
	require.Error(h.t, err, msgAndArgs...)
}

// RequireContains checks if a string contains a substring
func (h *TestHelper) RequireContains(s, substr string, msgAndArgs ...interface{}) {
	require.Contains(h.t, s, substr, msgAndArgs...)
}

// RequireNotEqual checks if two values are not equal
func (h *TestHelper) RequireNotEqual(expected, actual interface{}, msgAndArgs ...interface{}) {
	require.NotEqual(h.t, expected, actual, msgAndArgs...)
}
