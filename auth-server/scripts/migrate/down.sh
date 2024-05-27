#!/bin/bash

goose -dir internal/db/migrations sqlite3 .data/database.db down

