package heimdall

import (
	"database/sql"
	"errors"
	"fmt"
	"hash/fnv"

	log "github.com/Sirupsen/logrus"
	_ "github.com/lib/pq"
)

type Lock struct {
	Database  *sql.DB
	Namespace int32
	Name      int32
}

func New(database, namespace, name string) (lock *Lock, err error) {
	db, err := sql.Open("postgres", database)
	if err != nil {
		err = errors.New("heimdall: unable to establish database connection")
	}

	return &Lock{
		Database:  db,
		Namespace: encode(namespace),
		Name:      encode(name),
	}, err

}

func (l *Lock) Acquire() (lockStatus bool, err error) {
	var lockAcquired bool

	rows, err := l.Database.Query("SELECT pg_try_advisory_lock($1,$2)", l.Namespace, l.Name)
	if err != nil {
		return false, errors.New("heimdall: unable to execute query for advisory lock")
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&lockAcquired); err != nil {
			return false, errors.New("heimdall: unable to scan query results")
		}
	}

	log.Debug(fmt.Sprintf("Lock acquired?: %v", lockAcquired))

	return lockAcquired, nil
}

func (l *Lock) Release() {
	l.Database.Close()
}

func encode(input string) int32 {
	log.Debug(fmt.Sprintf("Before encoding: %v", input))

	h := fnv.New32()
	h.Write([]byte(input))

	output := h.Sum32()

	log.Debug(fmt.Sprintf("After encoding: %v", output))
	log.Debug(fmt.Sprintf("After casting to int32: %v", int32(output)))

	return int32(output)
}
