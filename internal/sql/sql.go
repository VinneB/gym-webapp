package sql

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/VinneB/gym-webapp/internal/structapi"
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

var addExerciseText string = `INSERT INTO exercises (name, data) VALUES (:name, :data);`

var addWorkoutText string = `INSERT INTO workouts (user_email, start_time, end_time) VALUES (:user_email, :start_time, :end_time);`

var getUserWorkoutsText string = `SELECT id FROM workouts WHERE user_email=?;`

var getAllExercisesText string = `SELECT * FROM exercises;`

var rootPathDB string = "data/"
var dataPathDB string = rootPathDB + "data.db"

var db *sqlx.DB

var addExerciseStmt *sqlx.NamedStmt
var addWorkoutStmt *sqlx.NamedStmt
var getAllUserWorkoutsStmt *sqlx.Stmt
var getAllExercisesStmt *sqlx.Stmt

func Connect() error {
	temp_db, err := sqlx.Connect("sqlite", dataPathDB)
	if err != nil {
		log.Println(err)
		return err
	}
	db = temp_db
	temp_addExerciseStmt, err := db.PrepareNamed(addExerciseText)
	if err != nil {
		log.Println(err)
		return err
	}
	temp_addWorkoutStmt, err := db.PrepareNamed(addWorkoutText)
	if err != nil {
		log.Println(err)
		return err
	}
	temp_getUserWorkoutsStmt, err := db.Preparex(getUserWorkoutsText)
	if err != nil {
		log.Println(err)
		return err
	}
	temp_getAllExercisesStmt, err := db.Preparex(getAllExercisesText)
	if err != nil {
		log.Println(err)
		return err
	}
	addExerciseStmt = temp_addExerciseStmt
	addWorkoutStmt = temp_addWorkoutStmt
	getAllUserWorkoutsStmt = temp_getUserWorkoutsStmt
	getAllExercisesStmt = temp_getAllExercisesStmt
	log.Println("Connected to sqlite3 database")
	return nil
}

func CloseDatabase() {
	db.Close()
	addExerciseStmt.Close()
	addWorkoutStmt.Close()
	getAllUserWorkoutsStmt.Close()
	getAllExercisesStmt.Close()
	log.Println("Closed sqlite3 database")
}

func dep_AddExercise(exercise structapi.Exercise) error {
	jsonData, err := json.Marshal(exercise)
	if err != nil {
		log.Println(err)
		return err
	}
	stmt, err := db.Prepare(addExerciseText)
	if err != nil {
		log.Println(err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(exercise.Name, jsonData)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func AddExercise(exercise structapi.Exercise) error {
	sqlExercise := structapi.ExerciseSqlForm{}
	sqlExercise.Name = exercise.Name
	stringData, err := json.Marshal(exercise.MuscleFractions)
	if err != nil {
		log.Println(err)
		return err
	}
	sqlExercise.MuscleFractions = string(stringData)
	_, err = addExerciseStmt.Exec(&sqlExercise)
	if err != nil {
		log.Println(err)
		log.Println("ya")
		return err
	}
	return nil
}

func dep_AddWorkoutInstance(workout structapi.WorkoutInstance, user_email string) error {
	jsonData, err := json.Marshal(workout)
	if err != nil {
		log.Println(err)
		return err
	}
	stmt, err := db.Prepare(addWorkoutText)
	if err != nil {
		log.Println(err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(workout.StartTime.Unix(), user_email, jsonData)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func AddWorkoutInstance(workout structapi.WorkoutInstance, user_email string) error {
	_, err := addWorkoutStmt.Exec(workout)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
func GetAllUserWorkouts(email string) ([]structapi.WorkoutInstance, error) {
	workouts := []structapi.WorkoutInstance{}
	err := getAllUserWorkoutsStmt.Select(workouts, email)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return workouts, nil
}

func dep_GetAllUserWorkouts(email string) ([]structapi.WorkoutInstance, error) {
	stmt, err := db.Prepare(getUserWorkoutsText)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(email)
	if err != nil {
		log.Println(err)
		return nil, err
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
			log.Println(err)
			return nil, err
		}
		json.Unmarshal(rawWorkout, &workout)
		workouts = append(workouts, workout)
	}
	return workouts, nil
}

func GetExercises() ([]structapi.Exercise, error) {
	exercisesRaw := []structapi.ExerciseSqlForm{}
	err := getAllExercisesStmt.Select(&exercisesRaw)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	exercises := []structapi.Exercise{}
	for _, value := range exercisesRaw {
		exercise, err := sqlFormToExercise(value)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		exercises = append(exercises, exercise)
	}

	return exercises, nil
}

func dep_GetExercises() ([]structapi.Exercise, error) {
	stmt, err := db.Prepare(getAllExercisesText)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var exercises []structapi.Exercise
	var rawExercise []byte
	var exercise structapi.Exercise
	for rows.Next() {
		exercise = structapi.Exercise{}
		rawExercise = nil
		err := rows.Scan(&rawExercise)
		fmt.Printf("row: %s\n", rawExercise)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		json.Unmarshal(rawExercise, &exercise)
		exercises = append(exercises, exercise)
	}
	return exercises, nil
}

func sqlFormToExercise(exerciseRaw structapi.ExerciseSqlForm) (structapi.Exercise, error) {
	muscleFractions := []structapi.MuscleFraction{}
	err := json.Unmarshal([]byte(exerciseRaw.MuscleFractions), &muscleFractions)
	if err != nil {
		return structapi.Exercise{}, err
	}
	return structapi.Exercise{Name: exerciseRaw.Name, MuscleFractions: muscleFractions}, nil
}
