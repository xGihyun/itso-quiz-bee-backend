# Development

This covers the development flow of the back end.

## Architecture

It follows a mixture of:

- Repository Pattern
- Vertical Slice Architecture

## Project Structure

```
internal
├── api
│   ├── api.go
│   └── json.go
├── database
│   ├── database.go
│   └── migrations
│       ├── 0001_sample_migration.sql
├── middleware
│   └── logging.go
├── user
│   ├── dto.go
│   ├── handler.go
│   ├── repository.go
│   ├── service.go
│   └── user.go
└── ws
    ├── client.go
    ├── pool.go
    ├── service.go
    └── websocket.go
main.go
```

...
