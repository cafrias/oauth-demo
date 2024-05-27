# Auth Server

This is the Auth Server of a demo Oauth2 implementation. It is a simple server that provides endpoints to authenticate users and generate tokens.

## Getting Started

Requirements:
- Go 1.22
- Sqlite3
- Air 1.51
- Goose 3.2
- sqlc 1.26

Install Project dependencies:
```bash
go mod download
```

The app requires the following envinronment variables:
```bash

# The path to the sqlite database
export AUTH_SERVER_DB_PATH=.data/database.db
```

Run the migrations:
```bash
./scripts/migrate/up.sh
```

Start the dev server:
```bash
air
```

## Migrations

Inside the directory `scripts/migrate` you'll find the scripts to run the migrations.

Migration files live in `internal/db/migrations`.

## Generating DB Code

This project uses `sqlc` to generate the database code. The code is generated inside the directory `internal/db`.

Please read more about `sqlc` [here](https://docs.sqlc.dev/en/stable/index.html).
