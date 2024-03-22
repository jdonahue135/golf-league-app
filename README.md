# Golf League App

This is a webapp that tracks results for golf leagues.

## Local Setup

- Create Postgres database
- Run command `./run.sh`

## Testing

- Run command `go test ./...`

## Testing with Coverage

- First time setup: run `chmod +x coverage.sh` to register coverage command (edit default browser in command if needed)
- Run `./coverage.sh`

## Dependencies

- Built in Go version 1.19
- Uses the [chi router](https://github.com/go-chi/chi)
- Uses [alex edwards SCS](https://github.com/alexedwards/scs) session management
- Uses [nosurf](https://github.com/justinas/nosurf)

## References
