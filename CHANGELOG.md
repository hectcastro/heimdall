## 1.2.0

- Upgrade to Go 1.25 and project dependencies.
- Replace GitHub Actions PostgreSQL service with testcontainers-go for test isolation.
- Fix deferred lock.Release() not executing when using os.Exit().
- Fix New() returning non-nil Lock with nil Database field on sql.Open() error.
- Add database connection verification with db.Ping() after sql.Open().

## 1.1.0

- Upgrade to Go 1.17 and project dependencies.

## 1.0.0

- Upgrade to Go 1.15 and project dependencies.
- Migrate from Travis CI to GitHub Actions.

## 0.2.1-0.2.3

- Automatically add tagged versions to GitHub Releases.
- No Heimdall source code was modified in the release range above.

## 0.2.0

- Add support for `libpq` environment variables.
- Replace GPM dependency management with Godep.

## 0.1.0

- Initial release.
