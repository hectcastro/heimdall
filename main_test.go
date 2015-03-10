package main

import (
	"testing"
)

func TestRunningProgram(t *testing.T) {
	exitStatus := Run("true", []string{})

	if exitStatus != 0 {
		t.Errorf("Running `true` failed")
	}
}

func TestFailingProgram(t *testing.T) {
	exitStatus := Run("false", []string{})

	if exitStatus != 1 {
		t.Errorf("Running `false` didn't fail")
	}
}
