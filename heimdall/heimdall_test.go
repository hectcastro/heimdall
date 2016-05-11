package heimdall

import (
	"net/url"
	"os"
	"testing"

	uuid "github.com/satori/go.uuid"
)

func TestAcquire(t *testing.T) {
	databaseURL := os.Getenv("DATABASE_URL")
	namespace := uuid.NewV4().String()
	name := uuid.NewV4().String()

	lock, err := New(databaseURL, namespace, name)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lock.Release()

	lockAcquired, err := lock.Acquire()
	if err != nil {
		t.Errorf(err.Error())
	}

	if !lockAcquired {
		t.Errorf("Unable to acquire lock")
	}
}

func TestEncode(t *testing.T) {
	for i := 0; i <= 1000; i++ {
		lockID := encode(uuid.NewV4().String())

		if lockID < -2147483648 || lockID > 2147483647 {
			t.Errorf("Lock ID is out of range")
		}
	}
}

func TestLockContention(t *testing.T) {
	databaseURL := os.Getenv("DATABASE_URL")
	namespace := uuid.NewV4().String()
	name := uuid.NewV4().String()

	lock, err := New(databaseURL, namespace, name)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lock.Release()

	lockAcquired, err := lock.Acquire()
	if err != nil {
		t.Errorf(err.Error())
	}

	if lockAcquired {
		secondLock, err := New(databaseURL, namespace, name)
		if err != nil {
			t.Errorf(err.Error())
		}
		defer secondLock.Release()

		secondLockAcquired, err := secondLock.Acquire()
		if err != nil {
			t.Errorf(err.Error())
		}

		if secondLockAcquired {
			t.Errorf("Second lock acquired before first released")
		}
	}
}

func TestLibPqEnvironment(t *testing.T) {
	namespace := uuid.NewV4().String()
	name := uuid.NewV4().String()

	dbURL, _ := url.Parse(os.Getenv("DATABASE_URL"))

	os.Setenv("PGHOST", dbURL.Host)
	os.Setenv("PGUSER", dbURL.User.Username())
	os.Setenv("PGDATABASE", dbURL.Path[1:len(dbURL.Path)])

	params, _ := url.ParseQuery(dbURL.RawQuery)
	os.Setenv("PGSSLMODE", params["sslmode"][0])

	lock, err := New("", namespace, name)
	if err != nil {
		t.Errorf(err.Error())
	}
	defer lock.Release()

	lockAcquired, err := lock.Acquire()
	if err != nil {
		t.Errorf(err.Error())
	}

	if !lockAcquired {
		t.Errorf("Unable to acquire lock")
	}
}
