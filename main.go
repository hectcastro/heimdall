package main

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/hectcastro/heimdall/heimdall"
	log "github.com/sirupsen/logrus"
)

// Default values for the command line interface.
const (
	DefaultDatabaseURL   = ""
	DefaultLockName      = "heimdall"
	DefaultLockNamespace = "heimdall"
	DefaultLockTimeout   = 5
)

func main() {
	os.Exit(run())
}

func run() int {
	var debug bool
	var database, namespace, name string
	var timeout int

	flag.Usage = func() { fmt.Print(usage()) }
	flag.BoolVar(&debug, "debug", false, "Debug mode enabled")
	flag.StringVar(&database, "database", DefaultDatabaseURL, "A database URL")
	flag.StringVar(&namespace, "namespace", DefaultLockNamespace, "A lock namespace")
	flag.StringVar(&name, "name", DefaultLockName, "A lock name")
	flag.IntVar(&timeout, "timeout", DefaultLockTimeout, "A lock timeout")
	flag.Parse()

	args := flag.Args()

	if debug {
		log.SetLevel(log.DebugLevel)
	}

	if len(args) == 0 {
		fmt.Fprint(os.Stderr, errors.New("heimdall: you must supply a program to run"))
		return 1
	}

	program := args[0]
	programArgs := args[1:]

	log.Debug(fmt.Sprintf("Database: %v", database))
	log.Debug(fmt.Sprintf("Namespace: %v", namespace))
	log.Debug(fmt.Sprintf("Name: %v", name))
	log.Debug(fmt.Sprintf("Timeout: %v", timeout))
	log.Debug(fmt.Sprintf("Program: %v", program))
	log.Debug(fmt.Sprintf("Program arguments: %v", programArgs))

	ctx := context.Background()

	lock, err := heimdall.New(ctx, database, namespace, name)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		return 1
	}

	lockAcquired, err := lock.Acquire(ctx)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		return 1
	}
	defer lock.Release()

	if lockAcquired {
		log.Debug("Lock was acquired")

		return Run(program, programArgs, timeout)
	} else {
		log.Debug("Lock was not acquired")

		return 1
	}
}

// Run executes a program and returns its exit status. Its
// arguments are a program, an array of arguments to that
// program, and a timeout.
func Run(program string, args []string, timeout int) int {
	var exitStatus int
	var cmdOut, cmdErr bytes.Buffer

	cmd := exec.Command(program, args...)
	cmd.Stdout = &cmdOut
	cmd.Stderr = &cmdErr

	if err := cmd.Start(); err != nil {
		log.Error(err)
		return 1
	}

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case err := <-done:
		if err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				// The command's failure exit status
				exitStatus = exitError.Sys().(syscall.WaitStatus).ExitStatus()
			} else {
				exitStatus = 1
			}
			fmt.Fprint(os.Stderr, cmdErr.String())
		} else {
			// The command's successful exit status
			exitStatus = cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()
		}
	case <-time.After(timeoutDuration(timeout)):
		if cmd.Process != nil {
			if err := cmd.Process.Kill(); err != nil {
				log.Debug(fmt.Sprintf("Failed to kill process: %v", err))
			}
		}
		<-done

		log.Debug("Process killed due to timeout")

		// Print any stderr output captured before timeout
		fmt.Fprint(os.Stderr, cmdErr.String())

		exitStatus = 1
	}

	fmt.Fprint(os.Stdout, cmdOut.String())

	return exitStatus
}

// timeoutDuration converts its integer argument into a
// time.Duration. If 0 is passed as the timeout, the duration
// becomes the maximum integer value (simulating infinity).
func timeoutDuration(timeout int) time.Duration {
	if timeout == 0 {
		// Maximum integer timeout
		return time.Duration(int(^uint(timeout) >> 1))
	}

	return time.Duration(timeout) * time.Second
}

// usage returns the usage text for this program's command line
// interface.
func usage() string {
	helpText := `
Usage: heimdall [options] PROGRAM

  Run a program with an exclusive lock acquired from PostgreSQL

Options:

  --debug                      Debug mode enabled
  --database                   A database URL
  --namespace                  A lock namespace
  --name                       A lock name
  --timeout                    A lock timeout
`

	return strings.TrimSpace(helpText)
}
