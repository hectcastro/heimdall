package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/hectcastro/heimdall/heimdall"
)

// Default values for the command line interface.
const (
	DEFAULT_DATABASE_URL   = ""
	DEFAULT_LOCK_NAME      = "heimdall"
	DEFUALT_LOCK_NAMESPACE = "heimdall"
	DEFAULT_LOCK_TIMEOUT   = 5
)

func main() {
	var debug bool
	var database, namespace, name string
	var timeout int

	if os.Getenv("GOMAXPROCS") == "" {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	flag.Usage = func() { fmt.Print(usage()) }
	flag.BoolVar(&debug, "debug", false, "Debug mode enabled")
	flag.StringVar(&database, "database", DEFAULT_DATABASE_URL, "A database URL")
	flag.StringVar(&namespace, "namespace", DEFUALT_LOCK_NAMESPACE, "A lock namespace")
	flag.StringVar(&name, "name", DEFAULT_LOCK_NAME, "A lock name")
	flag.IntVar(&timeout, "timeout", DEFAULT_LOCK_TIMEOUT, "A lock timeout")
	flag.Parse()

	args := flag.Args()

	if debug {
		log.SetLevel(log.DebugLevel)
	}

	if len(args) == 0 {
		exitError(errors.New("heimdall: you must supply a program to run"))
	}

	program := args[0]
	programArgs := args[1:]

	log.Debug(fmt.Sprintf("Database: %v", database))
	log.Debug(fmt.Sprintf("Namespace: %v", namespace))
	log.Debug(fmt.Sprintf("Name: %v", name))
	log.Debug(fmt.Sprintf("Timeout: %v", timeout))
	log.Debug(fmt.Sprintf("Program: %v", program))
	log.Debug(fmt.Sprintf("Program arguments: %v", programArgs))

	lock, err := heimdall.New(database, namespace, name)
	if err != nil {
		exitError(err)
	}

	lockAcquired, err := lock.Acquire()
	if err != nil {
		exitError(err)
	}
	defer lock.Release()

	if lockAcquired {
		log.Debug("Lock was acquired")

		os.Exit(Run(program, programArgs, timeout))
	} else {
		log.Debug("Lock was not acquired")

		os.Exit(1)
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
		exitError(err)
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
		cmd.Process.Kill()
		<-done

		log.Debug("Process killed due to timeout")

		exitStatus = 1
	}

	fmt.Fprint(os.Stdout, cmdOut.String())

	return exitStatus
}

// timeoutDuration converts its integer argument into a
// time.Duration. If 0 is passes as the timeout, the duration
// becomes the maximum integer value (simulating infinity).
func timeoutDuration(timeout int) time.Duration {
	if timeout == 0 {
		// Maximum integer timeout
		return time.Duration(int(^uint(timeout) >> 1))
	} else {
		return time.Duration(timeout) * time.Second
	}
}

// exitError is a convenience function for printing an error
// message to Stderr and returning 1 as the program's exit status.
func exitError(err error) {
	fmt.Fprint(os.Stderr, err)

	os.Exit(1)
}

// usage returns the usage text for this program's command line
// interface.
func usage() string {
	helpText := `
Usage: heimdall [options] PROGRAM

  Run a proram with an exclusive lock acquired from PostgreSQL

Options:

  --debug                      Debug mode enabled
  --database                   A database URL
  --namespace                  A lock namespace
  --name                       A lock name
  --timeout                    A lock timeout
`

	return strings.TrimSpace(helpText)
}
