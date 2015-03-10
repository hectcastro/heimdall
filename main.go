package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/hectcastro/heimdall/heimdall"
)

const (
	DEFAULT_DATABASE_URL   = "postgres://postgres@localhost/postgres?sslmode=disable"
	DEFAULT_LOCK_NAME      = "heimdall"
	DEFUALT_LOCK_NAMESPACE = "heimdall"
)

func main() {
	var debug bool
	var database, namespace, name string

	if os.Getenv("GOMAXPROCS") == "" {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	flag.Usage = func() { fmt.Print(Usage()) }
	flag.BoolVar(&debug, "debug", false, "Debug mode enabled")
	flag.StringVar(&database, "database", DEFAULT_DATABASE_URL, "A database URL")
	flag.StringVar(&namespace, "namespace", DEFUALT_LOCK_NAMESPACE, "A lock namespace")
	flag.StringVar(&name, "name", DEFAULT_LOCK_NAME, "A lock name")
	flag.Parse()

	args := flag.Args()

	if debug {
		log.SetLevel(log.DebugLevel)
	}

	if len(args) == 0 {
		log.Fatal("You must supply a program to run")
	}

	program := args[0]
	programArgs := args[1:]

	log.Debug(fmt.Sprintf("Database: %v", database))
	log.Debug(fmt.Sprintf("Namespace: %v", namespace))
	log.Debug(fmt.Sprintf("Name: %v", name))
	log.Debug(fmt.Sprintf("Program: %v", program))
	log.Debug(fmt.Sprintf("Program arguments: %v", programArgs))

	lock := heimdall.New(database, namespace, name)

	lockAcquired := lock.Acquire()
	defer lock.Release()

	if lockAcquired {
		log.Debug("Lock was acquired")

		os.Exit(Run(program, programArgs))
	} else {
		log.Debug("Lock was not acquired")

		os.Exit(1)
	}
}

func Run(program string, args []string) int {
	var waitStatus syscall.WaitStatus
	var cmdOut, cmdErr bytes.Buffer

	cmd := exec.Command(program, args...)
	cmd.Stdout = &cmdOut
	cmd.Stderr = &cmdErr

	if err := cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			waitStatus = exitError.Sys().(syscall.WaitStatus)
		} else {
			return 1
		}

		fmt.Fprint(os.Stderr, cmdErr.String())
	} else {
		waitStatus = cmd.ProcessState.Sys().(syscall.WaitStatus)
	}

	fmt.Fprint(os.Stdout, cmdOut.String())

	return waitStatus.ExitStatus()
}

func Usage() string {
	helpText := `
Usage: heimdall [options] PROGRAM

  Run a proram with an exclusive lock acquired from PostgreSQL

Options:

  --debug                      Debug mode enabled
  --database                   A database URL
  --namespace                  A lock namespace
  --name                       A lock name
`

	return strings.TrimSpace(helpText)
}
