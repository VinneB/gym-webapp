package sql

import (
	"encoding/json"
	"errors"
	"fmt"
	"hash"
	"log"
	"slices"
	"strconv"

	"github.com/VinneB/gym-webapp/internal/structapi"
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

var addExerciseText string = `INSERT INTO exercises (name, data) VALUES (:name, :data);`

var addWorkoutText string = `INSERT INTO workouts (user_email, start_time, end_time, plan_name) VALUES (:user_email, :start_time, :end_time, :plan_name);`

var addSetText string = `INSERT INTO sets (exercise_name, reps, partial_reps, weight, workout_id, time, type, user_email) VALUES (:exercise_name, :reps, :partial_reps, :weight, :workout_id, :time, :type, :user_email)`

var getAllSetsOfExerciseText string = `SELECT * FROM sets WHERE user_email=? AND exercise_name=?;`

var getAllSetsOfExerciseInListText string = `SELECT * FROM sets WHERE user_email=? AND exercise_name IN (?);`

var getAllSetsText string = `SELECT * FROM sets WHERE user_email=?;`

var getUserWorkoutsText string = `SELECT * FROM workouts WHERE user_email=?;`

var getAllExercisesText string = `SELECT * FROM exercises;`

var getFirstExercisesText string = `SELECT * FROM exercises LIMIT ?;`

var getAllSetsWithWorkoutIdText string = `SELECT * FROM sets WHERE user_email=? AND workout_id=?;`

var getAllPlanWorkoutsText string = `SELECT * FROM plan_workouts WHERE user_email=?;`

var addPlanWorkoutText string = `INSERT INTO plan_workouts (name, user_email, data) VALUES (:name, :user_email, :data);`

var getPlanWorkoutText string = `SELECT * FROM plan_workouts WHERE user_email=? AND name=?;`

var addLiveWorkoutSessionText string = `INSERT INTO live_workout_session (user_email, start_time, end_time, plan_name, sets) VALUES (:user_email, :start_time, :end_time, :plan_name, :sets);`

var (
	rootPathDB string = "data/"
	dataPathDB string = rootPathDB + "data.db"
)

var db *sqlx.DB

var (
	addExerciseStmt             *sqlx.NamedStmt
	addWorkoutStmt              *sqlx.NamedStmt
	addSetStmt                  *sqlx.NamedStmt
	addPlanWorkoutStmt          *sqlx.NamedStmt
	getAllSetsOfExerciseStmt    *sqlx.Stmt
	getAllUserWorkoutsStmt      *sqlx.Stmt
	getAllExercisesStmt         *sqlx.Stmt
	getFirstExercisesStmt       *sqlx.Stmt
	getAllSetsStmt              *sqlx.Stmt
	getAllSetsWithWorkoutIdStmt *sqlx.Stmt
	getAllPlanWorkoutStmt       *sqlx.Stmt
	getPlanWorkoutStmt          *sqlx.Stmt
	addLiveWorkoutSessionStmt   *sqlx.NamedStmt
)

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
	temp_addSetStmt, err := db.PrepareNamed(addSetText)
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
	temp_getAllSetsOfExercise, err := db.Preparex(getAllSetsOfExerciseText)
	if err != nil {
		log.Println(err)
		return err
	}
	temp_getFirstExercises, err := db.Preparex(getFirstExercisesText)
	if err != nil {
		log.Println(err)
		return err
	}
	temp_getAllSets, err := db.Preparex(getAllSetsText)
	if err != nil {
		log.Println(err)
		return err
	}
	temp_getAllSetsWithWorkoutIdStmt, err := db.Preparex(getAllSetsWithWorkoutIdText)
	if err != nil {
		log.Println(err)
		return err
	}
	temp_getAllPlanWorkoutStmt, err := db.Preparex(getAllPlanWorkoutsText)
	if err != nil {
		log.Println(err)
		return err
	}
	temp_addPlanWorkoutStmt, err := db.PrepareNamed(addPlanWorkoutText)
	if err != nil {
		log.Println(err)
		return err
	}
	temp_getPlanWorkoutStmt, err := db.Preparex(getPlanWorkoutText)
	if err != nil {
		log.Println(err)
		return err
	}

	temp_addLiveWorkoutSessionStmt, err := db.PrepareNamed(addLiveWorkoutSessionText)
	if err != nil {
		log.Println(err)
		return err
	}

	addExerciseStmt = temp_addExerciseStmt
	addWorkoutStmt = temp_addWorkoutStmt
	getAllUserWorkoutsStmt = temp_getUserWorkoutsStmt
	getAllExercisesStmt = temp_getAllExercisesStmt
	addSetStmt = temp_addSetStmt
	getAllSetsOfExerciseStmt = temp_getAllSetsOfExercise
	getFirstExercisesStmt = temp_getFirstExercises
	getAllSetsStmt = temp_getAllSets
	getAllSetsWithWorkoutIdStmt = temp_getAllSetsWithWorkoutIdStmt
	getAllPlanWorkoutStmt = temp_getAllPlanWorkoutStmt
	addPlanWorkoutStmt = temp_addPlanWorkoutStmt
	getPlanWorkoutStmt = temp_getPlanWorkoutStmt
	addLiveWorkoutSessionStmt = temp_addLiveWorkoutSessionStmt
	log.Println("Connected to sqlite3 database")
	return nil
}

func CloseDatabase() {
	db.Close()
	addExerciseStmt.Close()
	addWorkoutStmt.Close()
	getAllUserWorkoutsStmt.Close()
	getAllExercisesStmt.Close()
	addSetStmt.Close()
	getAllSetsOfExerciseStmt.Close()
	getFirstExercisesStmt.Close()
	getAllSetsStmt.Close()
	getAllSetsWithWorkoutIdStmt.Close()
	getAllPlanWorkoutStmt.Close()
	addPlanWorkoutStmt.Close()
	addLiveWorkoutSessionStmt.Close()
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

func AddSet(exercise structapi.Set) error {
	_, err := addSetStmt.Exec(exercise)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func GetAllSets(userName string) ([]structapi.Set, error) {
	sets := []structapi.Set{}
	err := getAllSetsStmt.Select(&sets, userName)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return sets, nil
}

func GetAllSetsOfExercise(exerciseName string, userName string) ([]structapi.Set, error) {
	sets := []structapi.Set{}
	err := getAllSetsOfExerciseStmt.Select(&sets, userName, exerciseName)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return sets, nil
}

func GetAllSetsWithWorkoutId(workoutId int, userName string) ([]structapi.Set, error) {
	sets := []structapi.Set{}
	err := getAllSetsWithWorkoutIdStmt.Select(&sets, userName, workoutId)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return sets, nil
}

func GetAllSetsThatTargetAMuscle(muscle string, userName string) ([]structapi.Set, []structapi.Exercise, error) {
	validExercises, err := GetAllExercisesThatTargetAMuscle(muscle, userName)
	if err != nil {
		return nil, nil, err
	}
	exerciseNameList := []string{}
	for _, exercise := range validExercises {
		exerciseNameList = append(exerciseNameList, exercise.Name)
	}
	log.Printf("exerciseNameList= %s\n", exerciseNameList)
	query, args, err := sqlx.In(getAllSetsOfExerciseInListText, userName, exerciseNameList)
	if err != nil {
		return nil, nil, err
	}
	query = db.Rebind(query)
	log.Println("query: " + query)
	log.Println(args)
	sets := []structapi.Set{}
	err = db.Select(&sets, query, args...)
	if err != nil {
		log.Println(err)
	}
	return sets, validExercises, nil
}

func GetAllExercisesThatTargetAMuscle(muscle string, userName string) ([]structapi.Exercise, error) {
	exercises, err := GetExercises(0)
	results := []structapi.Exercise{}
	if err != nil {
		log.Println(err)
		return nil, err
	}
	for _, exercise := range exercises {
		muscleNames := []string{}
		for _, val := range exercise.MuscleFractions {
			muscleNames = append(muscleNames, val.Name)
		}
		if slices.Contains(muscleNames, muscle) {
			results = append(results, exercise)
		}
	}
	return results, nil
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

func AddWorkoutInstance(workout structapi.WorkoutInstance) (int64, error) {
	result, err := addWorkoutStmt.Exec(workout)
	if err != nil {
		log.Println(err)
		return -1, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Println(err)
		return -1, err
	}
	log.Println("id = " + strconv.Itoa(int(id)))
	return id, nil
}

func GetAllUserWorkouts(email string) ([]structapi.WorkoutInstance, error) {
	workouts := []structapi.WorkoutInstance{}
	err := getAllUserWorkoutsStmt.Select(&workouts, email)
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

func GetExercises(count int) ([]structapi.Exercise, error) {
	var err error
	exercisesRaw := []structapi.ExerciseSqlForm{}
	if count > 0 {
		err = getFirstExercisesStmt.Select(&exercisesRaw, count)
	} else {
		err = getAllExercisesStmt.Select(&exercisesRaw)
	}
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

func sqlFormToPlanWorkout(planWorkoutRaw structapi.PlanWorkoutSqlForm) (structapi.PlanWorkout, error) {
	planSets := []structapi.PlanSet{}
	err := json.Unmarshal([]byte(planWorkoutRaw.Sets), &planSets)
	if err != nil {
		return structapi.PlanWorkout{}, err
	}
	return structapi.PlanWorkout{Name: planWorkoutRaw.Name, Id: planWorkoutRaw.Id, UserEmail: planWorkoutRaw.UserEmail, Sets: planSets}, nil
}

func planWorkoutToSqlForm(planWorkout structapi.PlanWorkout) (structapi.PlanWorkoutSqlForm, error) {
	sqlFormSets, err := json.Marshal(planWorkout.Sets)
	if err != nil {
		return structapi.PlanWorkoutSqlForm{}, err
	}
	return structapi.PlanWorkoutSqlForm{Name: planWorkout.Name, Id: planWorkout.Id, UserEmail: planWorkout.UserEmail, Sets: string(sqlFormSets)}, nil
}

func GetAllPlanWorkouts(userEmail string) ([]structapi.PlanWorkout, error) {
	planWorkouts := []structapi.PlanWorkoutSqlForm{}
	err := getAllPlanWorkoutStmt.Select(&planWorkouts, userEmail)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	returnVal := []structapi.PlanWorkout{}
	for _, planWorkout := range planWorkouts {
		workout, err := sqlFormToPlanWorkout(planWorkout)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		returnVal = append(returnVal, workout)
	}
	return returnVal, nil
}

func AddPlanWorkout(planWorkout structapi.PlanWorkout) error {
	sqlFormPlanWorkout, err := planWorkoutToSqlForm(planWorkout)
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = addPlanWorkoutStmt.Exec(&sqlFormPlanWorkout)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func GetPlanWorkout(workoutPlanName string, userEmail string) (structapi.PlanWorkout, error) {
	planWorkout := []structapi.PlanWorkoutSqlForm{}
	err := getPlanWorkoutStmt.Select(&planWorkout, userEmail, workoutPlanName)
	if len(planWorkout) < 1 {
		newErr := errors.New("No Plan workout with this name found")
		return structapi.PlanWorkout{}, newErr
	}
	if err != nil {
		log.Println("getPlanWorkoutStmt failed: " + err.Error())
		return structapi.PlanWorkout{}, err
	}
	workout, err := sqlFormToPlanWorkout(planWorkout[0])
	if err != nil {
		log.Println("sqlFormToPlanWorkout failed: " + err.Error())
		return structapi.PlanWorkout{}, err
	}
	return workout, nil
}

func AddLiveWorkoutSession(session structapi.LiveWorkoutSession) error {
	sessionSqlForm, err := liveWorkoutSessionToSqlForm(session)
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = addLiveWorkoutSessionStmt.Exec(&sessionSqlForm)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func GetLiveWorkoutSession(userEmail string) error {
	return nil
}

func liveWorkoutSessionToSqlForm(session structapi.LiveWorkoutSession) (structapi.LiveWorkoutSessionSqlForm, error) {

	return structapi.LiveWorkoutSessionSqlForm{}, nil
}

func sqlFormToLiveWorkoutSession(session structapi.LiveWorkoutSessionSqlForm) (structapi.LiveWorkoutSession, error) {

	return structapi.LiveWorkoutSession{}, nil
}
