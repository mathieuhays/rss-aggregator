# rss-aggregator

Project based on boot.dev course

## Development

### Useful commands

`sqlc generate` run this in the project root to update `/internal/database`

`goose postgres postgres://host:port/database up` in `sql/schema` to migrate the db. replace `up` with `down` to rollback changes

## useful references

[SQLC Documentation](https://docs.sqlc.dev/en/latest/tutorials/getting-started-postgresql.html)

[SQLMock Documentation](https://pkg.go.dev/github.com/DATA-DOG/go-sqlmock)

## Going further

- [ ] Support pagination on listing endpoints
- [ ] Support different sorting and filtering options
- [ ] Create a CLI tool that interfaces with the API
- [ ] Support more types of feeds
- [ ] Add integration tests
- [ ] Add bookmarking / liking
- [ ] Create simple front-end. HTMX?