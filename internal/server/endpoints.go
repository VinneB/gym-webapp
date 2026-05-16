package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/VinneB/gym-webapp/internal/calc"
	"github.com/VinneB/gym-webapp/internal/middleware"
	"github.com/VinneB/gym-webapp/internal/sql"
	"github.com/VinneB/gym-webapp/internal/structapi"
)

func ExercisesGetHandler(w http.ResponseWriter, r *http.Request) {
}

func ExercisesPostHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Exercises Post Handler")
	renderer := middleware.NewTemplate()
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
		data, err = getData(r)
		if err != nil {
			log.Println(err)
			SendError(w, r, http.StatusUnprocessableEntity, "Error")
			break
		}
		break
	}
	renderer.Render(w, "addexercise_page", data)
}

func WorkoutsPostHandler(w http.ResponseWriter, r *http.Request) {
	renderer := middleware.NewTemplate()
	var data structapi.Data
	for {
		err := r.ParseForm()
		log.Println(r.Form)
		if err != nil {
			log.Println(err)
			SendError(w, r, http.StatusInternalServerError, err.Error())
			break
		}
		log.Print(r.Form)
		if (r.FormValue("type") == "Custom" && (len(r.Form["exerciseName"]) < 1 || len(r.Form["exerciseName"]) != len(r.Form["partialRepAmount"]) || len(r.Form["exerciseName"]) != len(r.Form["repAmount"]) || len(r.Form["exerciseName"]) != len(r.Form["setAmount"]) || len(r.Form["exerciseName"]) != len(r.Form["weightAmount"]))) || (r.FormValue("type") == "WorkoutPlan" && (len(r.Form["exerciseName"]) < 1 || len(r.Form["exerciseName"]) != len(r.Form["partialRepAmount"]) || len(r.Form["exerciseName"]) != len(r.Form["repAmount"]) || len(r.Form["exerciseName"]) != len(r.Form["weightAmount"]))) {
			log.Println("Failed length sanitation")
			SendError(w, r, http.StatusUnprocessableEntity, "Failed length sanitation")
			break
		}
		exercises, err := sql.GetExercises(0)
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

			var setAmountInt int
			if r.FormValue("type") == "Custom" {
				setAmountInt, err = strconv.Atoi(r.Form["setAmount"][index])
				if err != nil {
					log.Println("Failed to parse setAmount" + err.Error())
					SendError(w, r, http.StatusUnprocessableEntity, "InternalError")
					break
				}
			} else {
				setAmountInt = 1
			}

			weightAmountFloat, err := strconv.ParseFloat(r.Form["weightAmount"][index], 32)
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
			timeStringPrefix := strings.Split(r.Form["workout-start-time"][0], "T")[0] + "T"
			log.Println(timeStringPrefix)
			time, err := time.Parse("2006-01-02T15:04", timeStringPrefix+r.Form["startTime"][index])
			if err != nil {
				log.Println("Time format invalid")
				SendError(w, r, http.StatusUnprocessableEntity, "Bad time")
				break
			}
			for i := 0; i < setAmountInt; i++ {
				sets = append(sets, structapi.Set{ExerciseName: r.Form["exerciseName"][index], Reps: int16(repAmountInt), PartialReps: int16(partialRepAmountInt), Weight: float32(weightAmountFloat), Time: time})
			}
		}

		planName := ""
		if r.FormValue("type") == "WorkoutPlan" {
			planName = r.FormValue("plan")
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

		workoutId, err := sql.AddWorkoutInstance(structapi.WorkoutInstance{UserEmail: "not_implemented_yet", StartTime: starttime, EndTime: endtime, PlanName: planName})
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
		data, err = getData(r)
		if err != nil {
			log.Println(err)
			SendError(w, r, http.StatusUnprocessableEntity, "Error")
			break
		}

		log.Println("break")
		break

	}
	data.WorkoutInstanceType = r.FormValue("type")
	renderer.Render(w, "add_workout_form", data)
}

func WorkoutsGetHandler(w http.ResponseWriter, r *http.Request) {
}

func GraphsPostHandler(w http.ResponseWriter, r *http.Request) {
}

func GraphsDeleteHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}
	log.Println(r.Form)
	// renderer := middleware.NewTemplate()
	// data := structapi.Data{}

	// renderer.Render(w, "graph_page", data)
}

func GraphsCalcTypeGetHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r)
	calcType := r.URL.Query().Get("graph-dataset-calctype-select")
	log.Println(calcType)
	options, _ := calc.GetOptions(calcType)
	log.Println(options)
	middleware.NewTemplate().Render(w, "graph_dataset_calctype_option", options)
}

func GraphsGetHandler(w http.ResponseWriter, r *http.Request) {
	data := structapi.Data{}
	r.ParseForm()
	log.Println(r.Form)
	calcType := r.Form["graph-dataset-calctype-select"][0]
	option := r.Form["graph-dataset-calctype-option-select"][0]
	log.Println(calcType + " " + option)
	points, _, _ := calc.Calculate(calcType, option, "not_implemented_yet")
	log.Println(points)
	data.ChartDataSets = append(data.ChartDataSets, structapi.ChartDataSet{Label: option, Points: points, CalcType: calcType})
	graphAppData, _ := json.Marshal(data.ChartDataSets)
	data.ChartDataSetsJson = template.JS(graphAppData)
	middleware.NewTemplate().Render(w, "graph_chart", data)
	middleware.NewTemplate().Render(w, "graph_dataset_info", data)
}

func WorkoutPlansPostHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	log.Println(r.Form)
	if err != nil {
		log.Println(err)
	}
	workoutPlan := structapi.PlanWorkout{}
	planName := r.Form["newWorkoutPlanName"][0]
	log.Println("test2")
	planSets := []structapi.PlanSet{}
	for i := range r.Form["exercise-name"] {
		exerciseName := r.Form["exercise-name"][i]
		lowerRepLimit, err := strconv.Atoi(r.Form["lower-rep-limit"][i])
		upperRepLimit, err := strconv.Atoi(r.Form["upper-rep-limit"][i])
		repsInReserve, err := strconv.Atoi(r.Form["reps-in-reserve"][i])
		setAmount, err := strconv.Atoi(r.Form["set-amount"][i])
		if err != nil {
			log.Println(err)
			return
		}
		for j := 0; j < setAmount; j++ {
			planSets = append(planSets, structapi.PlanSet{ExerciseName: exerciseName, RepLowerRange: int16(lowerRepLimit), RepUpperRange: int16(upperRepLimit), RIR: int16(repsInReserve)})
		}
	}
	workoutPlan = structapi.PlanWorkout{Name: planName, UserEmail: "not_implemented_yet", Sets: planSets}
	sql.AddPlanWorkout(workoutPlan)
	data, err := getData(r)
	if err != nil {
		log.Println("Failed to parse data")
	}
	middleware.NewTemplate().Render(w, "workoutplans_page", data)
}

func AddWorkoutSelectTypeGetHandler(w http.ResponseWriter, r *http.Request) {
	data, err := getData(r)
	if err != nil {
		log.Println(err)
		return
	}
	data.WorkoutInstanceType = r.FormValue("type")
	middleware.NewTemplate().Render(w, "add_workout_form", data)
	log.Println(data.WorkoutInstanceType)
}

func AddWorkoutSelectPlanGetHandler(w http.ResponseWriter, r *http.Request) {
	data, err := getData(r)
	if err != nil {
		log.Println(err)
		return
	}
	data.WorkoutInstanceType = r.FormValue("type")
	r.ParseForm()
	fmt.Println(r.Form)
	middleware.NewTemplate().Render(w, "add_workout_planned_form", data)
}
