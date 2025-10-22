package heimdall

import (
	"context"
	"log"
	"net/url"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var testDatabaseURL string

func TestMain(m *testing.M) {
	ctx := context.Background()

	// Start PostgreSQL container
	postgresContainer, err := postgres.Run(ctx,
		"postgres:17",
		postgres.WithDatabase("test"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("heimdall"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2),
		),
	)
	if err != nil {
		log.Fatalf("Failed to start PostgreSQL container: %v", err)
	}

	// Get connection string
	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to get connection string: %v", err)
	}

	testDatabaseURL = connStr

	// Run tests
	code := m.Run()

	// Cleanup
	if err := testcontainers.TerminateContainer(postgresContainer); err != nil {
		log.Fatalf("Failed to terminate container: %v", err)
	}

	os.Exit(code)
}

func TestAcquire(t *testing.T) {
	namespace := uuid.New().String()
	name := uuid.New().String()

	lock, err := New(testDatabaseURL, namespace, name)
	if err != nil {
		t.Error(err)
	}
	defer lock.Release()

	lockAcquired, err := lock.Acquire()
	if err != nil {
		t.Error(err)
	}

	if !lockAcquired {
		t.Errorf("Unable to acquire lock")
	}
}

func TestEncode(t *testing.T) {
	for i := 0; i <= 1000; i++ {
		lockID := encode(uuid.New().String())

		if lockID < -2147483648 || lockID > 2147483647 {
			t.Errorf("Lock ID is out of range")
		}
	}
}

func TestLockContention(t *testing.T) {
	namespace := uuid.New().String()
	name := uuid.New().String()

	lock, err := New(testDatabaseURL, namespace, name)
	if err != nil {
		t.Error(err)
	}
	defer lock.Release()

	lockAcquired, err := lock.Acquire()
	if err != nil {
		t.Error(err)
	}

	if lockAcquired {
		secondLock, err := New(testDatabaseURL, namespace, name)
		if err != nil {
			t.Error(err)
		}
		defer secondLock.Release()

		secondLockAcquired, err := secondLock.Acquire()
		if err != nil {
			t.Error(err)
		}

		if secondLockAcquired {
			t.Errorf("Second lock acquired before first released")
		}
	}
}

func TestLibPqEnvironment(t *testing.T) {
	namespace := uuid.New().String()
	name := uuid.New().String()

	dbURL, _ := url.Parse(testDatabaseURL)

	// Split host and port
	host := dbURL.Hostname()
	port := dbURL.Port()

	os.Setenv("PGHOST", host)
	if port != "" {
		os.Setenv("PGPORT", port)
	}
	os.Setenv("PGUSER", dbURL.User.Username())
	password, _ := dbURL.User.Password()
	os.Setenv("PGPASSWORD", password)
	os.Setenv("PGDATABASE", dbURL.Path[1:len(dbURL.Path)])

	params, _ := url.ParseQuery(dbURL.RawQuery)
	os.Setenv("PGSSLMODE", params["sslmode"][0])

	lock, err := New("", namespace, name)
	if err != nil {
		t.Error(err)
	}
	defer lock.Release()

	lockAcquired, err := lock.Acquire()
	if err != nil {
		t.Error(err)
	}

	if !lockAcquired {
		t.Errorf("Unable to acquire lock")
	}
}

func TestNewWithUnreachableDatabase(t *testing.T) {
	namespace := uuid.New().String()
	name := uuid.New().String()

	// Use a valid connection string format but unreachable host
	// Port 54321 should not have PostgreSQL running
	unreachableURL := "postgres://postgres:heimdall@localhost:54321/test?sslmode=disable"

	lock, err := New(unreachableURL, namespace, name)

	if err == nil {
		t.Error("Expected error when connecting to unreachable database, got nil")
	}

	if lock != nil {
		t.Errorf("Expected nil lock when connection fails, got non-nil lock")
	}
}
