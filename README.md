# heimdall [![Build Status](https://travis-ci.org/hectcastro/heimdall.svg?branch=develop)](https://travis-ci.org/hectcastro/heimdall)

This is an experimental program that wraps an executable program inside of an exclusive lock provided by PostgreSQL's `pg_try_advisory_lock`.

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
