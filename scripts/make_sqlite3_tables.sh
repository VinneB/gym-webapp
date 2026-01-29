#!/bin/bash

DB_FILE="../data/data.db"

# SQL commands to create tables and insert data
# The semicolon at the end of each SQL statement is important
SQL_COMMANDS="
CREATE TABLE IF NOT EXISTS users (
    user_id INTEGER PRIMARY KEY AUTOINCREMENT,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    pass TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS workouts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    start_time INTEGER NOT NULL,
    end_time INTEGER,
    user_email TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS exercises (
    name TEXT PRIMARY KEY,
    data TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS sets (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  exercise_name TEXT NOT NULL,
  reps INTEGER NOT NULL,
  partial_reps INTEGER NOT NULL,
  weight INTEGER,
  workout_id INTEGER NOT NULL,
  time TEXT,
  type TEXT,
  user_email TEXT NOT NULL
);

"

# Execute the SQL commands using sqlite3
sqlite3 "$DB_FILE" "$SQL_COMMANDS"

echo "Database '$DB_FILE' created and tables created."
