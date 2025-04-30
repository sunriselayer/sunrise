package e2e

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// E2eTestSuite represents the e2e test suite
type E2eTestSuite struct {
	suite.Suite
	ctx         context.Context
	cancel      context.CancelFunc
	containerId string
	helper      *TestHelper
}

// SetupSuite runs once before all tests
func (s *E2eTestSuite) SetupSuite() {
	s.ctx, s.cancel = context.WithTimeout(context.Background(), 5*time.Minute)
	s.helper = NewTestHelper(s.T())

	// Build and run the container
	require.NoError(s.T(), s.buildAndRunContainer())
}

const (
	KeyUser      = "user"
	KeyValidator = "validator"
)

// TearDownSuite runs once after all tests
func (s *E2eTestSuite) TearDownSuite() {
	if s.containerId != "" {
		// Stop and remove the container
		cmd := exec.CommandContext(s.ctx, "docker", "stop", s.containerId)
		require.NoError(s.T(), cmd.Run())

		cmd = exec.CommandContext(s.ctx, "docker", "rm", s.containerId)
		require.NoError(s.T(), cmd.Run())
	}

	if s.cancel != nil {
		s.cancel()
	}
}

// buildAndRunContainer builds and runs the application container
func (s *E2eTestSuite) buildAndRunContainer() error {
	// Get the absolute path to the project root
	projectRoot, err := filepath.Abs("../../")
	if err != nil {
		return fmt.Errorf("failed to get project root: %w", err)
	}

	// Build the Docker image
	buildCmd := exec.CommandContext(s.ctx, "docker", "build", "-t", "sunrise-e2e", "-f", "./tests/Dockerfile", ".")
	buildCmd.Dir = projectRoot
	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("failed to build Docker image: %w", err)
	}

	// Run the container
	runCmd := exec.CommandContext(s.ctx, "docker", "run", "-d",
		"--name", "sunrise-e2e-test",
		"-p", "26656:26656",
		"-p", "26657:26657",
		"-p", "1317:1317",
		"-p", "9090:9090",
		"sunrise-e2e")
	output, err := runCmd.Output()
	if err != nil {
		return fmt.Errorf("failed to run container: %w", err)
	}

	s.containerId = string(output)

	// Initialize the chain
	if output, err := s.execDockerCommand("sunrised", "init", "test", "--chain-id", "sunrise-test"); err != nil {
		return fmt.Errorf("failed to initialize chain: %w\nOutput: %s", err, output)
	}

	// Create a user account
	if _, err := s.execDockerCommand("sunrised", "keys", "add", KeyUser, "--keyring-backend", "test"); err != nil {
		return fmt.Errorf("failed to create user account: %w", err)
	}

	// Create a validator account
	if _, err := s.execDockerCommand("sunrised", "keys", "add", KeyValidator, "--keyring-backend", "test"); err != nil {
		return fmt.Errorf("failed to create validator account: %w", err)
	}

	userAddr, err := s.getUserAddress()
	if err != nil {
		return err
	}

	// Add genesis account
	validatorAddr, err := s.getValidatorAddress()
	if err != nil {
		return err
	}

	output, err = s.execDockerCommand("sunrised", "genesis", "add-genesis-account", userAddr, "1000000000uvrise,1000000000urise", "--keyring-backend", "test")
	if err != nil {
		return fmt.Errorf("failed to add genesis account: %w\nOutput: %s", err, output)
	}

	output, err = s.execDockerCommand("sunrised", "genesis", "add-genesis-account", validatorAddr, "1000000000uvrise,1000000000urise", "--keyring-backend", "test")
	if err != nil {
		return fmt.Errorf("failed to add genesis account: %w\nOutput: %s", err, output)
	}

	// Create gentx
	if _, err := s.execDockerCommand("sunrised", "genesis", "gentx", KeyValidator, "1000000000uvrise", "--chain-id", "sunrise-test", "--keyring-backend", "test"); err != nil {
		return fmt.Errorf("failed to create gentx: %w", err)
	}

	// Collect gentxs
	if _, err := s.execDockerCommand("sunrised", "genesis", "collect-gentxs"); err != nil {
		return fmt.Errorf("failed to collect gentxs: %w", err)
	}

	// Start the chain
	startCmd := exec.CommandContext(s.ctx, "docker", "exec", s.containerId, "sunrised", "start")
	if err := startCmd.Start(); err != nil {
		return fmt.Errorf("failed to start chain: %w", err)
	}

	// Wait for the chain to be ready
	maxRetries := 30
	for range maxRetries {
		if _, err := s.execDockerCommand("sunrised", "status"); err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}

	return nil
}

// execDockerCommand executes a command inside the container
func (s *E2eTestSuite) execDockerCommand(args ...string) ([]byte, error) {
	cmd := exec.CommandContext(s.ctx, "docker", append([]string{"exec", s.containerId}, args...)...)
	return cmd.CombinedOutput()
}

// getUserAddress returns the user address
func (s *E2eTestSuite) getUserAddress() (string, error) {
	output, err := s.execDockerCommand("sunrised", "keys", "show", KeyUser, "-a", "--keyring-backend", "test")
	if err != nil {
		return "", fmt.Errorf("failed to get user address: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

// getValidatorAddress returns the validator address
func (s *E2eTestSuite) getValidatorAddress() (string, error) {
	output, err := s.execDockerCommand("sunrised", "keys", "show", KeyValidator, "-a", "--keyring-backend", "test")
	if err != nil {
		return "", fmt.Errorf("failed to get validator address: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

// getBalance returns the balance of the specified address
func (s *E2eTestSuite) getBalance(address string) ([]byte, error) {
	return s.execDockerCommand("sunrised", "query", "bank", "balances", address)
}
