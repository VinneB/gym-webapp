package server

import (
	"log"
	"net/http"
	"slices"
	"strconv"
	"time"
	"unicode"

	"github.com/VinneB/gym-webapp/internal/sql"
	"github.com/VinneB/gym-webapp/internal/structapi"
)

func ExercisesGetHandler(w http.ResponseWriter, r *http.Request) {
}

func ExercisesPostHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Exercises Post Handler")
	renderer := newTemplate()
	var data structapi.Data
	for {
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
			SendError(w, r, http.StatusInternalServerError, err.Error())
			break
		}
		log.Print(r.Form)
		log.Printf("detailName ya: %s\n", r.Form["detailName"])
		log.Println("POST to exercises")
		// Length Sanitation
		log.Println(len(r.Form["notreal"]))
		if len(r.Form["muscleName"]) < 1 || len(r.Form["muscleDetail"]) != len(r.Form["muscleName"]) || len(r.Form["exerciseName"]) != 1 {
			log.Println("Failed length sanitation")
			SendError(w, r, http.StatusUnprocessableEntity, "Failed length sanitation")
			break
		}
		// Make fractional muscle objects
		var muscles []structapi.MuscleFraction
		for index, value := range r.Form["muscleName"] {
			// Value sanitation
			muscleDetailFloat, err := strconv.ParseFloat(r.Form["muscleDetail"][index], 32)
			capitalizedValue := []rune(value)
			capitalizedValue[0] = unicode.ToUpper(capitalizedValue[0])
			if err != nil {
				log.Println(err)
				SendError(w, r, http.StatusUnprocessableEntity, err.Error())
				break
			}
			if !slices.Contains(structapi.Muscles, string(capitalizedValue)) {
				log.Println("Invalid muscle name")
				SendError(w, r, http.StatusUnprocessableEntity, "Invalid muscle name")
				break
			} else if muscleDetailFloat > 1.0 || muscleDetailFloat < 0.0 {
				log.Println("Muscle fractional detail out of range")
				SendError(w, r, http.StatusUnprocessableEntity, "Invalid muscle name")
				break
			}
			muscles = append(muscles, structapi.MuscleFraction{Name: string(capitalizedValue), Fraction: float32(muscleDetailFloat)})
		}
		exercise := structapi.Exercise{Name: r.Form["exerciseName"][0], MuscleFractions: muscles}
		err = sql.AddExercise(exercise)
		if err != nil {
			log.Println(err)
			SendError(w, r, http.StatusUnprocessableEntity, err.Error())
			break
		}
		// Render
		data, err = getData(r.URL.Path)
		if err != nil {
			log.Println(err)
			SendError(w, r, http.StatusUnprocessableEntity, "Error")
			break
		}
		renderer.Render(w, "add_exercise_list", data)
		break
	}
	renderer.Render(w, "add_exercise_form", data)
}

func WorkoutsPostHandler(w http.ResponseWriter, r *http.Request) {
	renderer := newTemplate()
	var data structapi.Data
	for {
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
			SendError(w, r, http.StatusInternalServerError, err.Error())
			break
		}
		log.Print(r.Form)
		if len(r.Form["exerciseName"]) < 1 || len(r.Form["exerciseName"]) != len(r.Form["partialRepAmount"]) || len(r.Form["exerciseName"]) != len(r.Form["repAmount"]) || len(r.Form["exerciseName"]) != len(r.Form["setAmount"]) || len(r.Form["exerciseName"]) != len(r.Form["weightAmount"]) {
			log.Println("Failed length sanitation")
			SendError(w, r, http.StatusUnprocessableEntity, "Failed length sanitation")
			break
		}
		exercises, err := sql.GetExercises()
		if err != nil {
			log.Println(err)
			SendError(w, r, http.StatusInternalServerError, "InternalError")
			break
		}
		var exerciseNames []string
		for _, exercise := range exercises {
			exerciseNames = append(exerciseNames, exercise.Name)
		}
		var sets []structapi.Set
		for index := range r.Form["exerciseName"] {
			repAmountInt, err := strconv.Atoi(r.Form["repAmount"][index])
			if err != nil {
				log.Println("Failed to parse repAmount" + err.Error())
				SendError(w, r, http.StatusUnprocessableEntity, "InternalError")
				break
			}
			partialRepAmountInt, err := strconv.Atoi(r.Form["partialRepAmount"][index])
			if err != nil {
				log.Println("Failed to parse partialRepAmount" + err.Error())
				SendError(w, r, http.StatusUnprocessableEntity, "InternalError")
				break
			}

			setAmountInt, err := strconv.Atoi(r.Form["setAmount"][index])
			if err != nil {
				log.Println("Failed to parse setAmount" + err.Error())
				SendError(w, r, http.StatusUnprocessableEntity, "InternalError")
				break
			}

			weightAmountFloat, err := strconv.ParseFloat(r.Form["setAmount"][index], 32)
			if err != nil {
				log.Println("Failed to parse weightAmount" + err.Error())
				SendError(w, r, http.StatusUnprocessableEntity, "InternalError")
				break
			}

			if !slices.Contains(exerciseNames, r.Form["exerciseName"][index]) {
				log.Println("Invald exercise name")
				SendError(w, r, http.StatusInternalServerError, "InternalError")
				break
			}
			time, err := time.Parse("2006-01-02T15:04", r.Form["startTime"][index])
			if err != nil {
				log.Println("Time format invalid")
				SendError(w, r, http.StatusUnprocessableEntity, "Bad time")
				break
			}
			for i := 0; i < setAmountInt; i++ {
				sets = append(sets, structapi.Set{ExerciseName: r.Form["exerciseName"][index], Reps: int16(repAmountInt), PartialReps: int16(partialRepAmountInt), Weight: int16(weightAmountFloat), Time: time})
			}
		}

		starttime, err := time.Parse("2006-01-02T15:04", r.Form["workout-start-time"][0])
		if err != nil {
			log.Println("Workout start time format invalid")
			SendError(w, r, http.StatusUnprocessableEntity, "Bad time")
			break
		}
		endtime, err := time.Parse("2006-01-02T15:04", r.Form["workout-end-time"][0])
		if err != nil {
			log.Println("Workout end time format invalid")
			SendError(w, r, http.StatusUnprocessableEntity, "Bad time")
			break
		}
		workoutId, err := sql.AddWorkoutInstance(structapi.WorkoutInstance{UserEmail: "not_implemented_yet", StartTime: starttime, EndTime: endtime})
		log.Println(workoutId)
		if err != nil {
			log.Println("AddWorkoutInstance func failed" + err.Error())
			SendError(w, r, http.StatusInternalServerError, "InternalError")
			break
		}

		// Add each set to database
		for _, set := range sets {
			set.WorkoutId = int(workoutId)
			set.UserEmail = "not_implemented_yet"
			set.Type = "not_implemented_yet"
			err := sql.AddSet(set)
			if err != nil {
				log.Println("sql.AddSet func failed" + err.Error())
				SendError(w, r, http.StatusInternalServerError, "InternalError")
				break
			}
		}

		// Render
		data, err = getData(r.URL.Path)
		if err != nil {
			log.Println(err)
			SendError(w, r, http.StatusUnprocessableEntity, "Error")
			break
		}

		log.Println("break")
		break

	}
	renderer.Render(w, "add_workout_form", data)
}

func WorkoutsGetHandler(w http.ResponseWriter, r *http.Request) {
}
