package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"slices"
	"strings"
	"unicode"

	"github.com/VinneB/gym-webapp/internal/calc"
	"github.com/VinneB/gym-webapp/internal/middleware"
	"github.com/VinneB/gym-webapp/internal/sql"
	"github.com/VinneB/gym-webapp/internal/structapi"
	"github.com/VinneB/gym-webapp/internal/util"
)

var validPages []string = []string{"/addexercise", "/addworkout", "/graph", "/workoutsummary", "/workoutplans"}

func StartServer() {
	fmt.Println("Entering start server")
	mux := http.NewServeMux()
	mux.HandleFunc("/", htmlTemplateHandler)
	mux.HandleFunc("GET /htmx/exercises", ExercisesGetHandler)
	mux.HandleFunc("POST /htmx/exercises", ExercisesPostHandler)
	mux.HandleFunc("GET /htmx/workouts", WorkoutsGetHandler)
	mux.HandleFunc("POST /htmx/workouts", WorkoutsPostHandler)
	mux.HandleFunc("POST /htmx/graphs", GraphsPostHandler)
	mux.HandleFunc("DELETE /htmx/graphs", GraphsDeleteHandler)
	mux.HandleFunc("GET /htmx/graphs", GraphsGetHandler)
	mux.HandleFunc("GET /htmx/graph_calctype_option", GraphsCalcTypeGetHandler)
	mux.HandleFunc("POST /htmx/workoutplans", WorkoutPlansPostHandler)
	mux.HandleFunc("GET /htmx/addworkout/selecttype", AddWorkoutSelectTypeGetHandler)
	mux.HandleFunc("GET /htmx/addworkout/selectplan", AddWorkoutSelectPlanGetHandler)
	mux.HandleFunc("/static/", middleware.HandleFileServerWithLogging)
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "test/favicon.ico")
	})
	err := sql.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer sql.CloseDatabase()
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func htmlTemplateHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("html " + r.URL.Path)
	renderer := middleware.NewTemplate()
	data, err := getData(r)
	if err != nil {
		log.Println(err)
		SendError(w, r, http.StatusInternalServerError, "Sorry")
	}
	if slices.Contains(validPages, r.URL.Path) {
		// data.Page = r.URL.Path
		if r.Header.Get("HX-Request") == "true" {
			renderer.Render(w, data.ContentTemplate, data)
			return
		}
		renderer.Render(w, "layout", data)
		// renderer.Render(w, "index_body", data)
	}
	//else {
	//	data.Page = "error"
	//	fmt.Println("Error Page")
	//	renderer.Render(w, "index", data)
	//}
}

func htmxHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("htmx")
	renderer := middleware.NewTemplate()
	data := structapi.Data{}
	renderer.Render(w, "error_section", data)
	endpoint := strings.TrimPrefix(r.URL.Path, "/htmx/")
	log.Printf("htmx endpoint - %s\n", endpoint)
	switch endpoint {
	case "exercises":
	}
}

func SendError(w http.ResponseWriter, r *http.Request, statusCode int, err string) {
	w.WriteHeader(statusCode)
	data := structapi.Data{Errors: []string{err}, Page: "error"}
	middleware.NewTemplate().Render(w, "error_section", data)
	log.Println("Respondede")
}

func getData(r *http.Request) (structapi.Data, error) {
	data := structapi.Data{}
	page_rune := []rune(strings.TrimPrefix(r.URL.Path, "/"))
	data.Page = string(page_rune) + "_page"
	page_rune[0] = unicode.ToUpper(page_rune[0])
	data.Title = string(page_rune)
	data.ContentTemplate = data.Page
	log.Println("ContentTemplate: " + data.ContentTemplate)
	muscles, err := json.Marshal(structapi.Muscles)
	if err != nil {
		return structapi.Data{}, err
	}
	data.MusclesJson = template.JS(muscles)
	switch r.URL.Path {
	case "/addexercise":
		exercises, err := sql.GetExercises(0)
		if err != nil {
			return structapi.Data{}, err
		}
		data.Exercises = exercises
	case "/htmx/exercises":
		exercises, err := sql.GetExercises(0)
		if err != nil {
			return structapi.Data{}, err
		}
		data.Exercises = exercises
	case "/htmx/workouts":
		data.WorkoutInstanceType = r.FormValue("type")
		exercises, err := sql.GetExercises(0)
		if err != nil {
			return structapi.Data{}, err
		}
		data.Exercises = exercises
		exercisesJson, _ := json.Marshal(structapi.ExerciseToJsonFormat(exercises))
		data.ExercisesJson = template.JS(exercisesJson)
		if r.FormValue("type") == "WorkoutPlan" {
			log.Println("Plan: " + r.FormValue("plan"))
			planWorkout, err := sql.GetPlanWorkout(r.FormValue("plan"), "not_implemented_yet")
			if err != nil {
				return structapi.Data{}, err
			}
			data.WorkoutPlans = []structapi.PlanWorkout{planWorkout}
		}
	case "/addworkout":
		data.WorkoutInstanceType = "Custom"
		exercises, err := sql.GetExercises(0)
		if err != nil {
			return structapi.Data{}, err
		}
		data.Exercises = exercises
		exercisesJson, _ := json.Marshal(structapi.ExerciseToJsonFormat(exercises))
		data.ExercisesJson = template.JS(exercisesJson)
	case "/graph":
		// Default behavior is to grab the first exercise associated with current user, and display the volume over workouts
		// Use first exercise for user if no exercise was provided
		exercises, err := sql.GetExercises(0)
		points, _, err := calc.Calculate(calc.DefaultCalcType, exercises[0].Name, "not_implemented_yet")
		if err != nil {
			return structapi.Data{}, err
		}
		data.ChartDataSets = append(data.ChartDataSets, structapi.ChartDataSet{Label: exercises[0].Name, Points: points, CalcType: calc.DefaultCalcType})
		graphAppData, err := json.Marshal(data.ChartDataSets)
		data.ChartDataSetsJson = template.JS(graphAppData)
		data.CalcTypes = calc.CalcTypes
		data.Exercises = exercises
		data.GraphDatasetOptions, err = calc.GetOptions(calc.DefaultCalcType)
		if err != nil {
			return structapi.Data{}, err
		}
	case "/htmx/graphs":

	case "/workoutsummary":
		sets, err := sql.GetAllSets("not_implemented_yet")
		if err != nil {
			log.Println("error getting all sets")
			return structapi.Data{}, err
		}
		workouts, err := sql.GetAllUserWorkouts("not_implemented_yet")
		if err != nil {
			log.Println("error getting all workouts")
			return structapi.Data{}, err
		}
		sortedWorkouts := util.SortWorkoutsByDate(workouts)
		setByWorkoutIdMap := util.SeparateSetsIntoMapByWorkoutId(sets)
		exerciseVolumeByWorkoutIdMap := calc.ExerciseVolumeForSets(sets)
		muscleVolumeByWorkoutIdMap := calc.MuscleVolumeForSets(sets)
		for _, workout := range sortedWorkouts {
			workoutSummary := structapi.WorkoutSummary{
				Sets:           setByWorkoutIdMap[workout.Id],
				StartTime:      workout.StartTime.Format("2006-01-02 03:04 PM"),
				ExerciseVolume: exerciseVolumeByWorkoutIdMap[workout.Id],
				MuscleVolume:   muscleVolumeByWorkoutIdMap[workout.Id],
				Duration:       util.WorkoutDuration(workout),
				WorkoutPlan:    workout.PlanName,
			}
			data.WorkoutSummary = append(data.WorkoutSummary, workoutSummary)
		}
	case "/workoutplans":
		exercises, err := sql.GetExercises(0)
		data.Exercises = exercises
		exercisesJson, _ := json.Marshal(structapi.ExerciseToJsonFormat(exercises))
		data.ExercisesJson = template.JS(exercisesJson)
		if err != nil {
			return structapi.Data{}, err
		}
		workoutplans, err := sql.GetAllPlanWorkouts("not_implemented_yet")
		if err != nil {
			return structapi.Data{}, err
		}
		for i := range workoutplans {
			workoutplans[i].SetCollections = structapi.SetsToSetCollections(workoutplans[i].Sets)
		}
		data.WorkoutPlans = workoutplans
		log.Println(data.WorkoutPlans)
	case "/htmx/workoutplans":
		exercises, err := sql.GetExercises(0)
		data.Exercises = exercises
		exercisesJson, _ := json.Marshal(structapi.ExerciseToJsonFormat(exercises))
		data.ExercisesJson = template.JS(exercisesJson)
		if err != nil {
			return structapi.Data{}, err
		}
		workoutplans, err := sql.GetAllPlanWorkouts("not_implemented_yet")
		for i := range workoutplans {
			workoutplans[i].SetCollections = structapi.SetsToSetCollections(workoutplans[i].Sets)
		}
		if err != nil {
			return structapi.Data{}, err
		}
		data.WorkoutPlans = workoutplans
		log.Println(data)
	case "/htmx/addworkout/selecttype":
		exercises, err := sql.GetExercises(0)
		if err != nil {
			return structapi.Data{}, err
		}
		data.Exercises = exercises
		exercisesJson, _ := json.Marshal(structapi.ExerciseToJsonFormat(exercises))
		data.ExercisesJson = template.JS(exercisesJson)
		if r.FormValue("type") == "WorkoutPlan" {
			workoutplans, err := sql.GetAllPlanWorkouts("not_implemented_yet")
			if err != nil {
				log.Println(err)
				return structapi.Data{}, err
			}
			data.WorkoutPlans = workoutplans
		}

		log.Println("Hit")
	case "/htmx/addworkout/selectplan":
		log.Println("PlanWorkout: " + r.FormValue("plan"))
		planWorkout, err := sql.GetPlanWorkout(r.FormValue("plan"), "not_implemented_yet")
		if err != nil {
			return structapi.Data{}, err
		}
		data.WorkoutPlans = []structapi.PlanWorkout{planWorkout}

	default:
	}

	return data, nil
}
