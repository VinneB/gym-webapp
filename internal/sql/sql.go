package sql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/VinneB/gym-webapp/internal/structapi"
	_ "modernc.org/sqlite"
)

var addExerciseText string = `INSERT INTO exercises (id, data) VALUES (?,?);`

var addWorkoutText string = `INSERT INTO workouts (start_time, user_email, data) VALUES (?,?,?);`

var getUserWorkoutsText string = `SELECT data FROM workouts WHERE user_email=?;`

var rootPathDB string = "../../data/"
var dataPathDB string = rootPathDB + "data.db"

func Connect() *sql.DB {
	db, err := sql.Open("sqlite", dataPathDB)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to sqlite3 database")
	return db
}

func Close(db *sql.DB) {
	db.Close()
}

func AddExercise(db *sql.DB, exercise structapi.Exercise) {
	jsonData, err := json.Marshal(exercise.MuscleFractions)
	if err != nil {
		log.Fatal(err)
		return
	}
	stmt, err := db.Prepare(addExerciseText)
	defer stmt.Close()
	if err != nil {
		log.Fatal(err)
		return
	}
	_, err = stmt.Exec(exercise.Name, jsonData)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func AddWorkoutInstance(db *sql.DB, workout structapi.WorkoutInstance, user_email string) {
	jsonData, err := json.Marshal(workout)
	if err != nil {
		log.Fatal(err)
		return
	}
	stmt, err := db.Prepare(addWorkoutText)
	defer stmt.Close()
	if err != nil {
		log.Fatal(err)
		return
	}
	_, err = stmt.Exec(workout.StartTime.Unix(), user_email, jsonData)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func GetAllUserWorkouts(db *sql.DB, email string) []structapi.WorkoutInstance {
	stmt, err := db.Prepare(getUserWorkoutsText)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	defer stmt.Close()
	rows, err := stmt.Query(email)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	var workouts []structapi.WorkoutInstance
	var rawWorkout []byte
	var workout structapi.WorkoutInstance
	for rows.Next() {
		workout = structapi.WorkoutInstance{}
		rawWorkout = nil
		err := rows.Scan(&rawWorkout)
		fmt.Printf("row: %s\n", rawWorkout)
		if err != nil {
			log.Fatal(err)
		}
		json.Unmarshal(rawWorkout, &workout)
		workouts = append(workouts, workout)
	}
	return workouts
}
