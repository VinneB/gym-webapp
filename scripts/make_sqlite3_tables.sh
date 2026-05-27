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
    start_time DATETIME NOT NULL,
    end_time DATETIME,
    user_email TEXT NOT NULL,
    plan_name TEXT
);

CREATE TABLE IF NOT EXISTS live_workout_session (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    start_time DATETIME NOT NULL,
    end_time DATETIME,
    user_email TEXT NOT NULL,
    plan_name TEXT,
    sets TEXT
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
  weight REAL,
  workout_id INTEGER NOT NULL,
  time DATETIME,
  type TEXT,
  user_email TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS plan_workouts (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  user_email TEXT NOT NULL,
  data TEXT NOT NULL
);

"

# Execute the SQL commands using sqlite3
sqlite3 "$DB_FILE" "$SQL_COMMANDS"

echo "Database '$DB_FILE' created and tables created."
