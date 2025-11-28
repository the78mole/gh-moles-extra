package cmd

import (
	"testing"
)

func TestCleanupCmdExists(t *testing.T) {
	// Verify the cleanup command is registered properly
	if cleanupCmd == nil {
		t.Error("cleanupCmd should not be nil")
	}

	if cleanupCmd.Use != "cleanup [KEEP_COUNT]" {
		t.Errorf("expected Use to be 'cleanup [KEEP_COUNT]', got %s", cleanupCmd.Use)
	}
}

func TestCleanupCmdFlags(t *testing.T) {
	// Test that flags are defined
	yesFlag := cleanupCmd.Flags().Lookup("yes")
	if yesFlag == nil {
		t.Error("yes flag should be defined")
	}

	failedFlag := cleanupCmd.Flags().Lookup("failed")
	if failedFlag == nil {
		t.Error("failed flag should be defined")
	}

	// Test shorthand flags
	yesShort := cleanupCmd.Flags().ShorthandLookup("y")
	if yesShort == nil {
		t.Error("y shorthand flag should be defined")
	}

	failedShort := cleanupCmd.Flags().ShorthandLookup("f")
	if failedShort == nil {
		t.Error("f shorthand flag should be defined")
	}
}

func TestRunCmdExists(t *testing.T) {
	// Verify the run command is registered properly
	if runCmd == nil {
		t.Error("runCmd should not be nil")
	}

	if runCmd.Use != "run" {
		t.Errorf("expected Use to be 'run', got %s", runCmd.Use)
	}
}

func TestRootCmdExists(t *testing.T) {
	// Verify the root command is registered properly
	if rootCmd == nil {
		t.Error("rootCmd should not be nil")
	}

	if rootCmd.Use != "gh-moles" {
		t.Errorf("expected Use to be 'gh-moles', got %s", rootCmd.Use)
	}
}
