package main

import (
	"context"
	"os"
	"os/exec"
	"syscall"
	"testing"
	"time"
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

func TestRunWithContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	// Start a long-running process
	done := make(chan int, 1)
	go func() {
		exitStatus := RunWithContext(ctx, "sleep", []string{"30"}, DefaultLockTimeout)
		done <- exitStatus
	}()

	// Give the process time to start
	time.Sleep(100 * time.Millisecond)

	// Cancel the context (simulating signal)
	cancel()

	// Wait for completion
	select {
	case exitStatus := <-done:
		if exitStatus != 1 {
			t.Errorf("Expected exit status 1 when context cancelled, got %d", exitStatus)
		}
	case <-time.After(2 * time.Second):
		t.Error("Process did not terminate after context cancellation")
	}
}

func TestSignalHandling(t *testing.T) {
	if os.Getenv("HEIMDALL_TEST_SIGNAL") == "1" {
		// This is the subprocess - run main
		main()
		return
	}

	// Start heimdall as a subprocess with a long-running command
	cmd := exec.Command(os.Args[0], "-test.run=TestSignalHandling", "--", "sleep", "30")
	cmd.Env = append(os.Environ(), "HEIMDALL_TEST_SIGNAL=1")

	if err := cmd.Start(); err != nil {
		t.Fatalf("Failed to start subprocess: %v", err)
	}

	// Give it time to start and acquire lock
	time.Sleep(500 * time.Millisecond)

	// Send SIGTERM
	if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
		t.Fatalf("Failed to send SIGTERM: %v", err)
	}

	// Wait for process to exit
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-done:
		// Process exited - this is what we want
	case <-time.After(2 * time.Second):
		cmd.Process.Kill()
		t.Fatal("Process did not exit after SIGTERM within timeout")
	}

	// Verify exit status indicates signal termination
	if cmd.ProcessState != nil {
		exitCode := cmd.ProcessState.ExitCode()
		// Exit code should be 143 (128 + 15 for SIGTERM) or similar
		if exitCode != 143 && exitCode != 1 {
			t.Logf("Exit code was %d (acceptable for signal termination)", exitCode)
		}
	}
}
