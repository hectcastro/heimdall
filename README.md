# heimdall

[![Build Status](https://travis-ci.org/hectcastro/heimdall.svg?branch=develop)](https://travis-ci.org/hectcastro/heimdall)
[![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/hectcastro/heimdall/heimdall)
[![Apache V2 License](http://img.shields.io/badge/license-Apache%20V2-blue.svg)](https://github.com/hectcastro/heimdall/blob/develop/LICENSE)

This is an experimental program that wraps an executable program inside of an
exclusive lock provided by PostgreSQL's `pg_try_advisory_lock`.

`pg_try_advisory_lock` is a sibling of `pg_advisory_lock`, which is a locking
function with two signatures. One accepts a 64-bit integer, and the other accepts
2 32-bit integers. The main difference between `pg_advisory_lock` and
`pg_try_advisory_lock` is that `pg_try_advisory_lock` does not wait until a lock
becomes available. Instead, `pg_try_advisory_lock` returns immediately with a `t`
if it successfully acquired the lock, and `f` if it was unable to acquire it.

One last behavior that is important to note: if the connection to PostgreSQL is
severed, any lock acquired via `pg_try_advisory_lock` is automatically released.

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

A quick way to test an exclusive lock is to use `sleep` as the command.
Executing `heimdall` once with `sleep 10` and then immediately with another
command in another shell should prevent the second command from executing.

### Quick

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

There are two ways to run the built-in test suite. One is intended to be run
locally, while the other is setup to run within a Docker container.

### Local

The local tests require a PostgreSQL database connection fed to the test suite
as a connection string via the `DATABASE_URL` environment variable:

```bash
$ make deps
$ DATABASE_URL="postgres://hector@localhost/hector?sslmode=disable" make test
```

### Docker

The Docker set includes a PostgreSQL container that gets linked to the container
running the test suite. Ensure that [`docker-compose`](https://docs.docker.com/compose/)
is installed before attempting to use this setup.

```bash
$ make docker-test
```

## Release

Releases of `heimdall` are standalone binaries built for specific platforms. By
default, the `Makefile` in this project emits binaries for 64-bit Linux and
Darwin architectures.

### Local

If you have Go setup locally, you can use the following `make` target to build
binaries tucked under the `pkg` directory:

```bash
$ make release
```

### Docker

If you don't have Go setup locally, or want to isolate the build process inside
of a Linux container, the included `Dockerfile` can cross-compile `heimdall`
binaries for several platforms:

```bash
$ make docker-release
```

## Attribution

I was made aware of `pg_try_advisory_lock` through
[@ryandotsmith](https://github.com/ryandotsmith)'s work in
[lock-smith](https://github.com/ryandotsmith/lock-smith).
