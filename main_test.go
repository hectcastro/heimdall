package main

import (
	"testing"
)

func TestRunningProgram(t *testing.T) {
	exitStatus := Run("true", []string{}, DefaultLockTimeout)

	if exitStatus != 0 {
		t.Errorf("Running `true` failed")
	}
}

func TestFailingProgram(t *testing.T) {
	exitStatus := Run("false", []string{}, DefaultLockTimeout)

	if exitStatus != 1 {
		t.Errorf("Running `false` didn't fail")
	}
}

func TestTimeout(t *testing.T) {
	exitStatus := Run("sleep", []string{"3"}, 1)

	if exitStatus != 1 {
		t.Errorf("Timeout didn't kill process")
	}
}

func TestNonexistentCommand(t *testing.T) {
	exitStatus := Run("this-command-definitely-does-not-exist", []string{}, DefaultLockTimeout)

	if exitStatus != 1 {
		t.Errorf("Nonexistent command should return exit status 1, got %d", exitStatus)
	}
}
