#!/bin/bash

# This script is used to run db migrations
goose -dir internal/db/migrations sqlite3 .data/database.db up

