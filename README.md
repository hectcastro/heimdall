# heimdall [![Build Status](https://travis-ci.org/hectcastro/heimdall.svg?branch=develop)](https://travis-ci.org/hectcastro/heimdall)

This is an experimental program that wraps an executable program inside of an exclusive lock provided by PostgreSQL's `pg_try_advisory_lock`.

`pg_try_advisory_lock` is a sibling of `pg_advisory_lock`, which is a locking function with two signatures. One accepts a 64-bit integer, and the other accepts 2 32-bit integers. The main difference between `pg_advisory_lock` and `pg_try_advisory_lock` is that `pg_try_advisory_lock` does not wait until a lock becomes available. Instead, `pg_try_advisory_lock` returns immediately with a `t` if it successfully acquired the lock, and `f` if it was unable to acquire it.

One last behiavor that is important to note: if the connection to PostgreSQL is severed, any lock acquired via `pg_try_advisory_lock` is automatically released.

## Usage

The `heimdall` program has four required arguments:

- A database URL (`postgres://postgres@localhost/postgres?sslmode=disable`)
- A lock namespace
- A lock name
- A command to wrap

An example wrapping the `ls` command looks like:

```bash
$ heimdall \
    --database "postgres://hector@localhost/hector?sslmode=disable" \
    --namespace joker \
    --name test \
    ls
```

## Testing

A quick way to test an exclusive lock is to use `sleep` as the command. Executing `heimdall` once with `sleep 10` and then immediately with another command in another shell should prevent the second command from executing.

Induce a deep sleep:

```bash
$ heimdall \
    --database "postgres://hector@localhost/hector?sslmode=disable" \
    --namespace joker \
    --name test \
    sleep 10
```

Try to get the date:

```bash
$ heimdall \
    --database "postgres://hector@localhost/hector?sslmode=disable" \
    --namespace joker \
    --name test \
    date
```

## Attribution

I was made aware of `pg_try_advisory_lock` through [@ryandotsmith](https://github.com/ryandotsmith)'s work in [lock-smith](https://github.com/ryandotsmith/lock-smith).
