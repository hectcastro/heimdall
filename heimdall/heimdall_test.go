package heimdall

import (
	"os"
	"testing"

	uuid "github.com/satori/go.uuid"
)

func TestAcquire(t *testing.T) {
	databaseUrl := os.Getenv("DATABASE_URL")
	namespace := uuid.NewV4().String()
	name := uuid.NewV4().String()

	lock := New(databaseUrl, namespace, name)

	lockAcquired := lock.Acquire()
	defer lock.Release()

	if !lockAcquired {
		t.Errorf("Unable to acquire lock")
	}
}

func TestEncode(t *testing.T) {
	for i := 0; i <= 1000; i++ {
		lockId := encode(uuid.NewV4().String())

		if lockId < -2147483648 || lockId > 2147483647 {
			t.Errorf("Lock ID is out of range")
		}
	}
}

func TestLockContention(t *testing.T) {
	databaseUrl := os.Getenv("DATABASE_URL")
	namespace := uuid.NewV4().String()
	name := uuid.NewV4().String()

	lock := New(databaseUrl, namespace, name)
	defer lock.Release()

	if lock.Acquire() {
		secondLock := New(databaseUrl, namespace, name)
		defer secondLock.Release()

		if secondLock.Acquire() {
			t.Errorf("Second lock acquired before first released")
		}
	}
}
